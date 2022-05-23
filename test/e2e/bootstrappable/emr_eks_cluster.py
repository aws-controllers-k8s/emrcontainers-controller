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

import base64
import boto3
import kubernetes
import tempfile
import yaml

from dataclasses import dataclass, field
from typing import Union
from botocore import session
from awscli.customizations.eks.get_token import STSClientFactory, TokenGenerator

from acktest.aws.identity import get_account_id
from acktest.bootstrapping import Bootstrappable
from acktest.bootstrapping.iam import ServiceLinkedRole
from acktest.bootstrapping.eks import Cluster as EKSCluster

EMR_K8S_ROLE_NAME = "emr-containers"
EMR_K8S_USER_NAME = "emr-containers"

AWS_AUTH_NAMESPACE = "kube-system"
AWS_AUTH_CONFIG_MAP_NAME = "aws-auth"

@dataclass
class EMREnabledEKSCluster(Bootstrappable):
    # Inputs
    name_prefix: str
    emr_namespace: str

    # Subresources
    cluster: EKSCluster = field(init=False, default=None)
    emr_slr: ServiceLinkedRole = field(init=False, default=None)

    # Outputs

    def __post_init__(self):
        self.cluster = EKSCluster(f'{self.name_prefix}-cluster')
        self.emr_slr = ServiceLinkedRole("emr-containers.amazonaws.com", "AWSServiceRoleForAmazonEMRContainers")

    @property
    def eks_client(self):
        return boto3.client("eks", region_name=self.region)

    @property
    def eks_resource(self):
        return boto3.resource("eks", region_name=self.region)

    def _write_cafile(self, data: str) -> tempfile.NamedTemporaryFile:
        cafile = tempfile.NamedTemporaryFile(delete=False)

        cadata_b64 = data
        cadata = base64.b64decode(cadata_b64)

        cafile.write(cadata)
        cafile.flush()

        return cafile

    def _get_eks_token(self, cluster_name: str) -> str:
        sts_client = STSClientFactory(session.get_session()).get_sts_client()
        return TokenGenerator(sts_client).get_token(cluster_name)

    def _k8s_api_client(self, endpoint: str, token: str, cafile: str) -> kubernetes.client.CoreV1Api:
        kconfig = kubernetes.config.kube_config.Configuration(
            host=endpoint,
            api_key={'authorization': 'Bearer ' + token}
        )
        kconfig.ssl_ca_cert = cafile
        return kubernetes.client.ApiClient(configuration=kconfig)

    def bootstrap(self):
        """Creates an EKS cluster and installs the EMR components into it.
        """
        super().bootstrap()

        # # check if cluster is Active
        # try:
        #     cluster = self.eks_client.describe_cluster(name=self.cluster.name)
        #     assert cluster is not None
        #     cluster_endpoint = cluster["cluster"]["endpoint"]
        # except self.eks_client.exceptions.ResourceNotFoundException:
        #     pytest.fail(f"Could not find cluster with '{self.cluster.name}' in EKS")

        cluster = self.eks_client.describe_cluster(name=self.cluster.name)
        cluster_endpoint = cluster["cluster"]["endpoint"]

        cert_data = cluster["cluster"]["certificateAuthority"]["data"]
        ca_file = self._write_cafile(cert_data)

        api_client = self._k8s_api_client(
            cluster_endpoint,
            self._get_eks_token(self.cluster.name),
            ca_file.name,
        )

        core_v1 = kubernetes.client.CoreV1Api(api_client)

        # Create the EMR namespace
        namespaces = core_v1.list_namespace()
        if not any(ns.metadata.name == self.emr_namespace for ns in namespaces.items):
            emr_ns = kubernetes.client.V1Namespace(metadata=kubernetes.client.V1ObjectMeta(name=self.emr_namespace))
            core_v1.create_namespace(emr_ns)

        rbac_v1 = kubernetes.client.RbacAuthorizationV1Api(api_client)

        # Create the EMR RBAC
        roles = rbac_v1.list_namespaced_role(self.emr_namespace)
        if not any(role.metadata.name == EMR_K8S_ROLE_NAME for role in roles.items):
            rbac_v1.create_namespaced_role(self.emr_namespace, kubernetes.client.V1Role(
                metadata=kubernetes.client.V1ObjectMeta(name=EMR_K8S_ROLE_NAME, namespace=self.emr_namespace),
                rules=[
                    kubernetes.client.V1PolicyRule(
                        api_groups=[""],
                        resources=["namespaces"],
                        verbs=["get"],
                    ),
                    kubernetes.client.V1PolicyRule(
                        api_groups=[""],
                        resources=["serviceaccounts", "services", "configmaps", "events", "pods", "pods/log"],
                        verbs=["get", "list", "watch", "describe", "create", "edit", "delete", "deletecollection", "annotate", "patch", "label"],
                    ),
                    kubernetes.client.V1PolicyRule(
                        api_groups=[""],
                        resources=["secrets"],
                        verbs=["create", "patch", "delete", "watch"],
                    ),
                    kubernetes.client.V1PolicyRule(
                        api_groups=["apps"],
                        resources=["statefulsets", "deployments"],
                        verbs=["get", "list", "watch", "describe", "create", "edit", "delete", "annotate", "patch", "label"],
                    ),
                    kubernetes.client.V1PolicyRule(
                        api_groups=["batch"],
                        resources=["jobs"],
                        verbs=["get", "list", "watch", "describe", "create", "edit", "delete", "annotate", "patch", "label"],
                    ),
                    kubernetes.client.V1PolicyRule(
                        api_groups=["extensions"],
                        resources=["ingresses"],
                        verbs=["get", "list", "watch", "describe", "create", "edit", "delete", "annotate", "patch", "label"],
                    ),
                    kubernetes.client.V1PolicyRule(
                        api_groups=["rbac.authorization.k8s.io"],
                        resources=["roles", "rolebindings"],
                        verbs=["get", "list", "watch", "describe", "create", "edit", "delete", "deletecollection", "annotate", "patch", "label"],
                    ),
                ]
            ))

        # Create the role binding
        bindings = rbac_v1.list_namespaced_role_binding(self.emr_namespace)
        if not any(binding.metadata.name == EMR_K8S_ROLE_NAME for binding in bindings.items):
            rbac_v1.create_namespaced_role_binding(self.emr_namespace, kubernetes.client.V1RoleBinding(
                metadata=kubernetes.client.V1ObjectMeta(name=EMR_K8S_ROLE_NAME, namespace=self.emr_namespace),
                subjects=[
                    kubernetes.client.V1Subject(
                        api_group="rbac.authorization.k8s.io",
                        kind="User",
                        name=EMR_K8S_USER_NAME,
                    )
                ],
                role_ref=kubernetes.client.V1RoleRef(kind="Role", name=EMR_K8S_ROLE_NAME, api_group="rbac.authorization.k8s.io")
            ))

        # Patch the auth configmap
        current_auth = core_v1.read_namespaced_config_map(AWS_AUTH_CONFIG_MAP_NAME, AWS_AUTH_NAMESPACE)
        map_roles = current_auth.data["mapRoles"]
        map_roles = yaml.safe_load(map_roles)

        if not any(role["username"] == EMR_K8S_USER_NAME for role in map_roles):
            map_roles.append({
                "rolearn": f"arn:aws:iam::{get_account_id()}:role/{self.emr_slr.role_name}",
                "username": EMR_K8S_USER_NAME
            })
            auth_patch = [{
                "op": "replace",
                "path": "/data/mapRoles",
                "value": yaml.dump(map_roles)
            }]
            core_v1.patch_namespaced_config_map(AWS_AUTH_CONFIG_MAP_NAME, AWS_AUTH_NAMESPACE, auth_patch)

    def cleanup(self):
        """Deletes the EKS cluster and all associated resources.
        """
        super().cleanup()
