# Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License"). You may
# not use this file except in compliance with the License. A copy of the
# License is located at
#
# http://aws.amazon.com/apache2.0/
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
JR_RESOURCE_PLURAL = "jobruns"

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the zone has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 180

# Maximum time to wait for EKS cluster to be active (5 minutes)
MAX_EKS_WAIT_SECONDS = 300


@pytest.fixture
def eks_client():
    return boto3.client("eks")


def wait_for_eks_cluster_active(eks_client, cluster_name: str, max_wait_seconds: int = MAX_EKS_WAIT_SECONDS) -> bool:
    """Wait for EKS cluster to be in ACTIVE status

    Args:
        eks_client: boto3 EKS client
        cluster_name: Name of the EKS cluster
        max_wait_seconds: Maximum time to wait in seconds

    Returns:
        bool: True if cluster is active, False if timeout reached
    """
    start_time = time.time()
    while (time.time() - start_time) < max_wait_seconds:
        try:
            response = eks_client.describe_cluster(name=cluster_name)
            status = response['cluster']['status']
            if status == 'ACTIVE':
                return True
            elif status in ['FAILED', 'DELETING', 'DELETED']:
                logging.error(
                    f"EKS cluster {cluster_name} in terminal state: {status}")
                return False
            logging.info(f"Waiting for EKS cluster to be active.")
            if not wait_for_eks_cluster_active(eks_client, cluster_name):
                pytest.fail(f"EKS cluster did not become active within {MAX_EKS_WAIT_SECONDS} seconds")
        except eks_client.exceptions.ResourceNotFoundException:
            logging.error(f"EKS cluster not found")
            return False
        except Exception as e:
            logging.warning(f"Error checking EKS cluster status: {str(e)}")
        time.sleep(30)
    return False


@pytest.fixture
def iam_client():
    return boto3.client("iam")


@pytest.fixture
def jobrun(eks_client):
    virtual_cluster_name = random_suffix_name("emr-virtual-cluster", 32)
    job_run_name = random_suffix_name("emr-job-run", 32)

    # Wait for EKS cluster to be active before proceeding
    eks_cluster_name = get_bootstrap_resources().HostCluster_JR.cluster.name
    if not wait_for_eks_cluster_active(eks_client, eks_cluster_name):
        pytest.fail(f"EKS cluster {eks_cluster_name} did not become active within {MAX_EKS_WAIT_SECONDS}seconds")

    replacements = REPLACEMENT_VALUES.copy()
    replacements["VIRTUALCLUSTER_NAME"] = virtual_cluster_name
    replacements["EKS_CLUSTER_NAME"] = eks_cluster_name

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

    # Add retry mechanism for VirtualCluster creation
    max_retries = 5
    retry_delay = 30  # seconds
    for attempt in range(max_retries):
        try:
            k8s.create_custom_resource(vc_ref, resource_data)
            vc_cr = k8s.wait_resource_consumed_by_controller(vc_ref)

            # Check if the resource exists and has an ID
            if vc_cr is not None and k8s.get_resource_exists(vc_ref):
                if "status" in vc_cr and "id" in vc_cr["status"]:
                    break

            # If we get here, the creation succeeded but ID is not set
            # Delete and retry
            k8s.delete_custom_resource(vc_ref, 3, 10)

        except Exception as e:
            logging.warning(f"Attempt {attempt + 1} failed: {str(e)}")
            if attempt < max_retries - 1:
                time.sleep(retry_delay)
                continue
            raise

    assert vc_cr is not None
    assert k8s.get_resource_exists(vc_ref)
    assert "status" in vc_cr and "id" in vc_cr["status"], "VirtualCluster status.id not set after retries"

    virtual_cluster_id = vc_cr["status"]["id"]
    emr_release_label = "emr-6.3.0-latest"
    eks_clustername = get_bootstrap_resources().HostCluster_JR.cluster.name
    job_execution_role = get_bootstrap_resources().JobExecutionRole.arn
    emr_eks_s3_logs_name = get_bootstrap_resources().EMREKSS3BucketName.name

    replacements = REPLACEMENT_VALUES.copy()
    replacements["JOBRUN_NAME"] = job_run_name
    replacements["VIRTUALCLUSTER_NAME"] = virtual_cluster_name
    replacements["EMR_RELEASE_LABEL"] = emr_release_label
    replacements["JOB_EXECUTION_ROLE"] = job_execution_role
    replacements["EMREKSS3BucketName"] = emr_eks_s3_logs_name

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

    # Introducing sleep for emr job to finish
    time.sleep(CHECK_STATUS_WAIT_SECONDS)

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
class Test_JobRun:

    def base36_str_to_int(self, request):
        """Method to convert given string into decimal representation"""
        result = 0
        for char in request:
            result = result * 256 + ord(char)

        return result

    def base36_encode(self, request):
        """Method to return base36 encoded form of the input string"""
        decimal_number = self.base36_str_to_int(str(request))
        alphabet, base36 = ['0123456789abcdefghijklmnopqrstuvwxyz', '']

        while decimal_number:
            decimal_number, i = divmod(decimal_number, 36)
            base36 = alphabet[i] + base36

        return base36 or alphabet[0]

    def check_if_statement_exists(self, expected_statement, actual_assume_role_document):
        if actual_assume_role_document is None:
            return False

        existing_statements = actual_assume_role_document.get("Statement", [])
        for existing_statement in existing_statements:
            matches = self.check_if_dict_matches(
                expected_statement, existing_statement)
            if matches:
                return True
        return False

    def check_if_dict_matches(self, expected_dict, actual_dict):
        if len(expected_dict) != len(actual_dict):
            return False
        for key in expected_dict:
            key_str = str(key)
            val = expected_dict[key_str]
            if isinstance(val, dict):
                if not check_if_dict_matches(val, actual_dict.get(key_str, {})):
                    return False
            else:
                if key_str not in actual_dict or actual_dict[key_str] != str(val):
                    return False
        return True

    def get_assume_role_policy(self, iam_client, job_execution_role_name):
        """Method to retrieve trust policy of given role name"""
        role = self.iam_client.get_role(RoleName=job_execution_role_name)
        return role.get("Role").get("AssumeRolePolicyDocument")

    def update_assume_role(self, oidc_provider_arn, iam_client):
        job_execution_role_arn = get_bootstrap_resources().JobExecutionRole.arn
        job_execution_role_name = job_execution_role_arn.split('role/')[1]
        oidc_provider = oidc_provider_arn.split('oidc-provider/')[1]
        emr_namespace = "emr-ns"
        base36_encoded_role_name = self.base36_encode(job_execution_role_name)
        account_id = get_account_id()
        LOG = logging.getLogger(__name__)
        TRUST_POLICY_STATEMENT_ALREADY_EXISTS = "Trust policy statement already " \
            "exists for role %s. No changes " \
            "were made!"
        TRUST_POLICY_UPDATE_SUCCESSFUL = "Successfully updated trust policy of role %s"
        TRUST_POLICY_STATEMENT_FORMAT = '{ \
        "Effect": "Allow", \
        "Principal": { \
            "Federated": "%(OIDC_PROVIDER_ARN)s" \
        }, \
        "Action": "sts:AssumeRoleWithWebIdentity", \
        "Condition": { \
            "StringLike": { \
                "%(OIDC_PROVIDER)s:sub": "system:serviceaccount:%(NAMESPACE)s' \
                                    ':emr-containers-sa-*-*-%(AWS_ACCOUNT_ID)s-' \
                                    '%(BASE36_ENCODED_ROLE_NAME)s" \
            } \
        } \
        }'

        job_execution_trust_policy = json.loads(TRUST_POLICY_STATEMENT_FORMAT % {
            "AWS_ACCOUNT_ID": account_id,
            "OIDC_PROVIDER_ARN": oidc_provider_arn,
            "OIDC_PROVIDER": oidc_provider,
            "NAMESPACE": emr_namespace,
            "BASE36_ENCODED_ROLE_NAME": base36_encoded_role_name
        })

        assume_role_policy = iam_client.get_role(
            RoleName=job_execution_role_name)
        assume_role_document = assume_role_policy.get(
            "Role").get("AssumeRolePolicyDocument")

        matches = self.check_if_statement_exists(job_execution_trust_policy,
                                                 assume_role_document)

        if not matches:
            existing_statements = assume_role_document.get("Statement")
            if existing_statements is None:
                assume_role_document["Statement"] = [
                    job_execution_trust_policy]
            else:
                existing_statements.append(job_execution_trust_policy)

            iam_client.update_assume_role_policy(
                RoleName=job_execution_role_name, PolicyDocument=json.dumps(assume_role_document))
            return TRUST_POLICY_UPDATE_SUCCESSFUL % job_execution_role_name
        else:
            return TRUST_POLICY_STATEMENT_ALREADY_EXISTS % job_execution_role_name

    def test_create_delete_jobrun(self, jobrun, emrcontainers_client, iam_client):
        oidc_provider_arn = get_bootstrap_resources().HostCluster_JR.export_oidc_arn

        # Update Job Execution Role
        role_update = self.update_assume_role(oidc_provider_arn, iam_client)
        assert role_update

        (vc_ref, vc_cr, jr_ref, jr_cr) = jobrun
        assert vc_cr, jr_cr

        virtual_cluster_id = vc_cr["status"]["id"]
        assert virtual_cluster_id

        jobrun_id = jr_cr["status"]["id"]
        assert jobrun_id

        try:
            aws_res = emrcontainers_client.describe_job_run(
                id=jobrun_id, virtualClusterId=virtual_cluster_id)
            assert aws_res is not None
        except emrcontainers_client.exceptions.ResourceNotFoundException:
            pytest.fail(f"Could not find job run with ID in EMR on EKS")

        # delete oidc provider
        try:
            aws_res = iam_client.delete_open_id_connect_provider(
                OpenIDConnectProviderArn=oidc_provider_arn)
            assert aws_res is not None
        except iam_client.exceptions.InvalidInputException:
            pytest.fail(f"Could not delete oidc identity provider")

        # check if JobRun is deleted
        try:
            jr_deleted = emrcontainers_client.describe_job_run(
                id=jobrun_id, virtualClusterId=virtual_cluster_id)
            logging.debug('%s is deleted during cleanup', jobrun_id)
            assert jr_deleted
        except:
            logging.debug('some resources such as %s did not cleanup as expected', jobrun_id)

        # check if VirtualCluster is deleted
        try:
            vc_deleted = emrcontainers_client.describe_virtual_cluster(
                id=virtual_cluster_id)
            logging.debug('%s is deleted during cleanup', virtual_cluster_id)
            assert vc_deleted
        except:
            logging.debug('some resources such as %s did not cleanup as expected', virtual_cluster_id)
