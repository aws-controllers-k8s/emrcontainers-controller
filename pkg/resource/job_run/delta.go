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

package job_run

import (
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}
	customPreCompare(delta, a, b)

	if ackcompare.HasNilDifference(a.ko.Spec.ExecutionRoleARN, b.ko.Spec.ExecutionRoleARN) {
		delta.Add("Spec.ExecutionRoleARN", a.ko.Spec.ExecutionRoleARN, b.ko.Spec.ExecutionRoleARN)
	} else if a.ko.Spec.ExecutionRoleARN != nil && b.ko.Spec.ExecutionRoleARN != nil {
		if *a.ko.Spec.ExecutionRoleARN != *b.ko.Spec.ExecutionRoleARN {
			delta.Add("Spec.ExecutionRoleARN", a.ko.Spec.ExecutionRoleARN, b.ko.Spec.ExecutionRoleARN)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.JobDriver, b.ko.Spec.JobDriver) {
		delta.Add("Spec.JobDriver", a.ko.Spec.JobDriver, b.ko.Spec.JobDriver)
	} else if a.ko.Spec.JobDriver != nil && b.ko.Spec.JobDriver != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.JobDriver.SparkSubmitJobDriver, b.ko.Spec.JobDriver.SparkSubmitJobDriver) {
			delta.Add("Spec.JobDriver.SparkSubmitJobDriver", a.ko.Spec.JobDriver.SparkSubmitJobDriver, b.ko.Spec.JobDriver.SparkSubmitJobDriver)
		} else if a.ko.Spec.JobDriver.SparkSubmitJobDriver != nil && b.ko.Spec.JobDriver.SparkSubmitJobDriver != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint, b.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint) {
				delta.Add("Spec.JobDriver.SparkSubmitJobDriver.EntryPoint", a.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint, b.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint)
			} else if a.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint != nil && b.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint != nil {
				if *a.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint != *b.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint {
					delta.Add("Spec.JobDriver.SparkSubmitJobDriver.EntryPoint", a.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint, b.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPoint)
				}
			}
			if !ackcompare.SliceStringPEqual(a.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPointArguments, b.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPointArguments) {
				delta.Add("Spec.JobDriver.SparkSubmitJobDriver.EntryPointArguments", a.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPointArguments, b.ko.Spec.JobDriver.SparkSubmitJobDriver.EntryPointArguments)
			}
			if ackcompare.HasNilDifference(a.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters, b.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters) {
				delta.Add("Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters", a.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters, b.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters)
			} else if a.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters != nil && b.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters != nil {
				if *a.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters != *b.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters {
					delta.Add("Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters", a.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters, b.ko.Spec.JobDriver.SparkSubmitJobDriver.SparkSubmitParameters)
				}
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Name, b.ko.Spec.Name) {
		delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
	} else if a.ko.Spec.Name != nil && b.ko.Spec.Name != nil {
		if *a.ko.Spec.Name != *b.ko.Spec.Name {
			delta.Add("Spec.Name", a.ko.Spec.Name, b.ko.Spec.Name)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ReleaseLabel, b.ko.Spec.ReleaseLabel) {
		delta.Add("Spec.ReleaseLabel", a.ko.Spec.ReleaseLabel, b.ko.Spec.ReleaseLabel)
	} else if a.ko.Spec.ReleaseLabel != nil && b.ko.Spec.ReleaseLabel != nil {
		if *a.ko.Spec.ReleaseLabel != *b.ko.Spec.ReleaseLabel {
			delta.Add("Spec.ReleaseLabel", a.ko.Spec.ReleaseLabel, b.ko.Spec.ReleaseLabel)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.Tags, b.ko.Spec.Tags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	} else if a.ko.Spec.Tags != nil && b.ko.Spec.Tags != nil {
		if !ackcompare.MapStringStringPEqual(a.ko.Spec.Tags, b.ko.Spec.Tags) {
			delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.VirtualClusterID, b.ko.Spec.VirtualClusterID) {
		delta.Add("Spec.VirtualClusterID", a.ko.Spec.VirtualClusterID, b.ko.Spec.VirtualClusterID)
	} else if a.ko.Spec.VirtualClusterID != nil && b.ko.Spec.VirtualClusterID != nil {
		if *a.ko.Spec.VirtualClusterID != *b.ko.Spec.VirtualClusterID {
			delta.Add("Spec.VirtualClusterID", a.ko.Spec.VirtualClusterID, b.ko.Spec.VirtualClusterID)
		}
	}
	if !reflect.DeepEqual(a.ko.Spec.VirtualClusterRef, b.ko.Spec.VirtualClusterRef) {
		delta.Add("Spec.VirtualClusterRef", a.ko.Spec.VirtualClusterRef, b.ko.Spec.VirtualClusterRef)
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ConfigurationOverrides, b.ko.Spec.ConfigurationOverrides) {
		delta.Add("Spec.ConfigurationOverrides", a.ko.Spec.ConfigurationOverrides, b.ko.Spec.ConfigurationOverrides)
	} else if a.ko.Spec.ConfigurationOverrides != nil && b.ko.Spec.ConfigurationOverrides != nil {
		if *a.ko.Spec.ConfigurationOverrides != *b.ko.Spec.ConfigurationOverrides {
			delta.Add("Spec.ConfigurationOverrides", a.ko.Spec.ConfigurationOverrides, b.ko.Spec.ConfigurationOverrides)
		}
	}

	return delta
}
