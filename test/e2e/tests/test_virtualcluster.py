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
import json
import logging
import time
from typing import Dict
import pytest

from acktest.k8s import resource as k8s
from acktest.resources import random_suffix_name
from acktest.aws.identity import get_account_id
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources

VC_RESOURCE_PLURAL = "virtualclusters"

@pytest.fixture
def iam_client():
    return boto3.client("iam")

@pytest.fixture
def virtualcluster():
    virtual_cluster_name = random_suffix_name("emr-virtual-cluster", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["VIRTUALCLUSTER_NAME"] = virtual_cluster_name
    replacements["EKS_CLUSTER_NAME"] = get_bootstrap_resources().HostCluster_VC.cluster.name

    resource_data = load_resource(
        "emr_virtual_cluster",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource for emr virtual cluster
    vc_ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, VC_RESOURCE_PLURAL,
        virtual_cluster_name, namespace="default",
    )
    k8s.create_custom_resource(vc_ref, resource_data)
    vc_cr = k8s.wait_resource_consumed_by_controller(vc_ref)

    assert vc_cr is not None
    assert k8s.get_resource_exists(vc_ref)

    yield (vc_ref, vc_cr)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(vc_ref, 3, 10)
        assert deleted
    except:
        pass


@service_marker
@pytest.mark.canary
class Test_VirtualCluster:
    def test_create_delete_virtualcluster(self, virtualcluster, emrcontainers_client, iam_client):
        oidc_provider_arn = get_bootstrap_resources().HostCluster_VC.export_oidc_arn

        (vc_ref, vc_cr) = virtualcluster
        assert vc_cr

        print("vc_cr=", vc_cr)

        virtual_cluster_id = vc_cr["status"]["id"]
        assert virtual_cluster_id

        try:
            aws_res = emrcontainers_client.describe_virtual_cluster(id=virtual_cluster_id)
            assert aws_res is not None
        except emrcontainers_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find virtual cluster with ID '{virtual_cluster_id}' in EMR on EKS")

        # delete oidc provider
        try:
            aws_res = iam_client.delete_open_id_connect_provider(OpenIDConnectProviderArn=oidc_provider_arn)
            assert aws_res is not None
        except iam_client.exceptions.InvalidInputException:
            pytest.fail(f"Could not delete oidc identity provider")

        # check if VirtualCluster is deleted
        try:
            vc_deleted = emrcontainers_client.describe_virtual_cluster(id=virtual_cluster_id)
            logging.debug('%s is deleted during cleanup', virtual_cluster_id)
            assert vc_deleted
        except:
            logging.debug('some resources such as %s did not cleanup as expected', virtual_cluster_id)
