# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
#	 http://aws.amazon.com/apache2.0/
#
# or in the "license" file accompanying this file. This file is distributed
# on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
# express or implied. See the License for the specific language governing
# permissions and limitations under the License.

"""Integration tests for the EMR on EKS VirtualCluster resource
"""

import boto3
import logging
import time
from typing import Dict

import pytest

from acktest.k8s import resource as k8s
from acktest.k8s import condition
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_eks_resource
from e2e.common.types import CLUSTER_RESOURCE_PLURAL
from e2e.replacement_values import REPLACEMENT_VALUES
from dataclasses import dataclass, field
# from e2e import bootstrap_directory
# from e2e.bootstrap_resources import BootstrapResources
from e2e.bootstrap_resources import get_bootstrap_resources

RESOURCE_PLURAL = "virtualclusters"

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the zone has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 10

def wait_for_cluster_active(eks_client, eks_cluster_name):
    waiter = eks_client.get_waiter('cluster_active')
    waiter.wait(name=eks_cluster_name)

@service_marker
@pytest.mark.canary
class TestVirtualCluster:
    def test_create_delete_ekscluster(self):
        eks_client = boto3.client("eks")
        eks_cluster_name = random_suffix_name("eks-cluster", 32)
        cluster_role = get_bootstrap_resources().ClusterRole.arn
        private_subnet_1 = get_bootstrap_resources().ClusterVPC.private_subnets.subnet_ids[0]
        private_subnet_2 = get_bootstrap_resources().ClusterVPC.private_subnets.subnet_ids[1]
        k8s_version = "1.20"

        # Create Kubernetes cluster.
        response = eks_client.create_cluster(
            name=eks_cluster_name,
            version=k8s_version,
            roleArn=cluster_role,
            resourcesVpcConfig={
                "subnetIds": [private_subnet_1, private_subnet_2]
            }
        )

        logging.info(f"Creating EKS Cluster {eks_cluster_name}")

        wait_for_cluster_active(eks_client, eks_cluster_name)

        assert eks_cluster_name is not None

        try:
            aws_res = eks_client.describe_cluster(name=eks_cluster_name)
            assert aws_res is not None
            logging.info(f"EKS Cluster {eks_cluster_name} is active")
        except eks_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find cluster '{eks_cluster_name}' in EKS")

        return eks_cluster_name

    def test_create_delete_virtualcluster(self):
        eks_cluster_name = self.test_create_delete_ekscluster()
        namespace_name = "default"
        emrcontainers_client = boto3.client("emr-containers")
        print("eks_cluster_name =", eks_cluster_name)

        virtual_cluster_name = random_suffix_name("emr-virtual-cluster", 32)

        replacements = REPLACEMENT_VALUES.copy()
        replacements["VIRTUALCLUSTER_NAME"] = virtual_cluster_name
        replacements["EKS_CLUSTER_NAME"] = eks_cluster_name
        replacements["NAMESPACE"] = namespace_name

        resource_data = load_eks_resource(
            "virtualcluster",
            additional_replacements=replacements,
        )
        logging.debug(resource_data)

        # Create the k8s resource
        ref = k8s.CustomResourceReference(
            CRD_GROUP, CRD_VERSION, RESOURCE_PLURAL,
            virtual_cluster_name, namespace=namespace_name,
        )
        k8s.create_custom_resource(ref, resource_data)
        cr = k8s.wait_resource_consumed_by_controller(ref)

        assert cr is not None
        assert k8s.get_resource_exists(ref)

        try:
            _, deleted = k8s.delete_custom_resource(ref, 3, 10)
            assert deleted
        except:
            pass

        virtual_cluster_id = cr["status"]["id"]
        print("virtual_cluster_id = ", virtual_cluster_id)

        assert virtual_cluster_id

        try:
            aws_res = emrcontainers_client.describe_virtual_cluster(id=virtual_cluster_id)
            assert aws_res is not None
        except emrcontainers_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find virtual cluster with ID '{virtual_cluster_id}' in EMR on EKS")

#     def cleanup(self):
#         """Deletes the EKS cluster.
#         """
        super().cleanup()

#         eks = self.ec2_resource.Vpc(self.vpc_id)
        eks_client.delete()
