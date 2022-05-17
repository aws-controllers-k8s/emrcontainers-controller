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
import time

from acktest.bootstrapping import Resources, BootstrapFailureException
from acktest.bootstrapping.iam import Role
from acktest.bootstrapping.vpc import VPC
from acktest.resources import random_suffix_name
from acktest.aws.identity import get_region
from e2e import bootstrap_directory
from e2e.bootstrap_resources import BootstrapResources, get_bootstrap_resources

# Time to wait after modifying the CR for the status to change
MODIFY_WAIT_AFTER_SECONDS = 10

# Time to wait after the zone has changed status, for the CR to update
CHECK_STATUS_WAIT_SECONDS = 10

def service_bootstrap() -> Resources:
    logging.getLogger().setLevel(logging.INFO)
    region = get_region()

    resources = BootstrapResources(
        ClusterRole=Role("cluster-role", "eks.amazonaws.com", managed_policies=["arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"]),
        NodegroupRole=Role("nodegroup-role", "ec2.amazonaws.com", managed_policies=[
            "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
            "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
            "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
        ]),
        ClusterVPC=VPC(name_prefix="cluster-vpc", num_public_subnet=2, num_private_subnet=2),
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
