// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

type CertificateProviderType string

const (
	CertificateProviderType_PEM CertificateProviderType = "PEM"
)

type ContainerProviderType string

const (
	ContainerProviderType_EKS ContainerProviderType = "EKS"
)

type EndpointState string

const (
	EndpointState_ACTIVE                 EndpointState = "ACTIVE"
	EndpointState_CREATING               EndpointState = "CREATING"
	EndpointState_TERMINATED             EndpointState = "TERMINATED"
	EndpointState_TERMINATED_WITH_ERRORS EndpointState = "TERMINATED_WITH_ERRORS"
	EndpointState_TERMINATING            EndpointState = "TERMINATING"
)

type FailureReason string

const (
	FailureReason_CLUSTER_UNAVAILABLE FailureReason = "CLUSTER_UNAVAILABLE"
	FailureReason_INTERNAL_ERROR      FailureReason = "INTERNAL_ERROR"
	FailureReason_USER_ERROR          FailureReason = "USER_ERROR"
	FailureReason_VALIDATION_ERROR    FailureReason = "VALIDATION_ERROR"
)

type JobRunState string

const (
	JobRunState_CANCELLED      JobRunState = "CANCELLED"
	JobRunState_CANCEL_PENDING JobRunState = "CANCEL_PENDING"
	JobRunState_COMPLETED      JobRunState = "COMPLETED"
	JobRunState_FAILED         JobRunState = "FAILED"
	JobRunState_PENDING        JobRunState = "PENDING"
	JobRunState_RUNNING        JobRunState = "RUNNING"
	JobRunState_SUBMITTED      JobRunState = "SUBMITTED"
)

type PersistentAppUI string

const (
	PersistentAppUI_DISABLED PersistentAppUI = "DISABLED"
	PersistentAppUI_ENABLED  PersistentAppUI = "ENABLED"
)

type TemplateParameterDataType string

const (
	TemplateParameterDataType_NUMBER TemplateParameterDataType = "NUMBER"
	TemplateParameterDataType_STRING TemplateParameterDataType = "STRING"
)

type VirtualClusterState string

const (
	VirtualClusterState_ARRESTED    VirtualClusterState = "ARRESTED"
	VirtualClusterState_RUNNING     VirtualClusterState = "RUNNING"
	VirtualClusterState_TERMINATED  VirtualClusterState = "TERMINATED"
	VirtualClusterState_TERMINATING VirtualClusterState = "TERMINATING"
)
