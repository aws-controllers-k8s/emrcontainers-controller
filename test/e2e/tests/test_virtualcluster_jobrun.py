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
from acktest.resources import random_suffix_name
from e2e import service_marker, CRD_GROUP, CRD_VERSION, load_resource
from e2e.replacement_values import REPLACEMENT_VALUES
from e2e.bootstrap_resources import get_bootstrap_resources

VC_RESOURCE_PLURAL = "virtualclusters"
JR_RESOURCE_PLURAL = "jobruns"

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the zone has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 10

@pytest.fixture
def virtualcluster_jobrun():
    virtual_cluster_name = random_suffix_name("emr-virtual-cluster", 32)
    job_run_name = random_suffix_name("emr-job-run", 32)

    replacements = REPLACEMENT_VALUES.copy()
    replacements["VIRTUALCLUSTER_NAME"] = virtual_cluster_name
    replacements["EKS_CLUSTER_NAME"] = get_bootstrap_resources().HostCluster.cluster.name

    resource_data = load_resource(
        "simple_cluster",
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

    virtual_cluster_id = vc_cr["status"]["id"]
    print(virtual_cluster_id)
    emr_release_label = "emr-6.3.0-latest"
    job_exection_role = get_bootstrap_resources().JobExecutionRole.arn

    try:
        aws_res = emrcontainers_client.update_role_trust_policy(cluster-name=EKS_CLUSTER_NAME,namespace=emr-ns)
        assert aws_res is not None
    except emrcontainers_client.exceptions.ResourceNotFoundException:
        pass

    replacements = REPLACEMENT_VALUES.copy()
    replacements["JOBRUN_NAME"] = job_run_name
    replacements["VIRTUALCLUSTER_ID"] = virtual_cluster_id
    replacements["EMR_RELEASE_LABEL"] = emr_release_label
    replacements["JOB_EXECUTION_ROLE"] = job_exection_role

    resource_data = load_resource(
        "job_run",
        additional_replacements=replacements,
    )
    logging.debug(resource_data)

    # Create the k8s resource for emr job run
    jr_ref = k8s.CustomResourceReference(
        CRD_GROUP, CRD_VERSION, JR_RESOURCE_PLURAL,
        job_run_name, namespace="default",
    )
    k8s.create_custom_resource(jr_ref, resource_data)
    jr_cr = k8s.wait_resource_consumed_by_controller(jr_ref)

    assert jr_cr is not None
    assert k8s.get_resource_exists(jr_ref)

    yield (vc_ref, vc_cr, jr_ref, jr_cr)

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(jr_ref, 3, 10)
        assert deleted
    except:
        pass

    # Try to delete, if doesn't already exist
    try:
        _, deleted = k8s.delete_custom_resource(vc_ref, 3, 10)
        assert deleted
    except:
        pass

@service_marker
@pytest.mark.canary
class TestVirtualCluster:
    def test_create_delete_virtualcluster_jobrun(self, virtualcluster_jobrun, emrcontainers_client):
        (vc_ref, vc_cr, jr_ref, jr_cr) = virtualcluster_jobrun
        assert vc_cr, jr_cr

        virtual_cluster_id = vc_cr["status"]["id"]
        assert virtual_cluster_id

        try:
            aws_res = emrcontainers_client.describe_virtual_cluster(id=virtual_cluster_id)
            assert aws_res is not None
        except emrcontainers_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find virtual cluster with ID '{virtual_cluster_id}' in EMR on EKS")

        jobrun_id = jr_cr["status"]["id"]
        assert jobrun_id

        try:
            aws_res = emrcontainers_client.describe_job_run(id=jobrun_id,virtualClusterId=virtual_cluster_id)
            assert aws_res is not None
        except emrcontainers_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find job run with ID '{jobrun_id}' in EMR on EKS")
