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
"""Bootstraps the resources required to run the EMRcontainers integration tests.
"""
import boto3
import logging
import json
import time

from acktest.bootstrapping import Resources, BootstrapFailureException
from acktest.bootstrapping.iam import Role, UserPolicies
from acktest.bootstrapping.s3 import Bucket
from e2e import bootstrap_directory
from e2e.bootstrap_resources import BootstrapResources
from e2e.bootstrappable.emr_eks_cluster import EMREnabledEKSCluster

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the zone has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 10

def service_bootstrap() -> Resources:
    logging.getLogger().setLevel(logging.INFO)

    job_execution_policy = json.dumps({
        "Version": "2012-10-17",
        "Statement": [
            {
                "Effect": "Allow",
                "Action": [
                    "s3:PutObject",
                    "s3:GetObject",
                    "s3:ListBucket"
                ],
                "Resource": "*"
            },
            {
                "Effect": "Allow",
                "Action": [
                    "logs:PutLogEvents",
                    "logs:CreateLogStream",
                    "logs:DescribeLogGroups",
                    "logs:DescribeLogStreams"
                ],
                "Resource": [
                    "arn:aws:logs:*:*:*"
                ]
            }
        ]
    })

    resources = BootstrapResources(
        JobExecutionRole=Role("ack-emrcontainers-job-execution-role", "ec2.amazonaws.com",
            user_policies=UserPolicies("ack-emrcontainers-job-execution-policy", [job_execution_policy])
        ),
        EMREKSS3BucketName=Bucket("ack-emr-eks-logs"),
        HostCluster_VC=EMREnabledEKSCluster("ack-emr-eks", "emr-ns"),
        HostCluster_JR=EMREnabledEKSCluster("ack-emr-eks", "emr-ns")
    )

    try:
        resources.bootstrap()
    except BootstrapFailureException as ex:
        exit(254)
    return resources

if __name__ == "__main__":
    config = service_bootstrap()
    # Write config to current directory by default
    config.serialize(bootstrap_directory)
