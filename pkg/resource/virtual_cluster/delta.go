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

package virtual_cluster

import (
	"bytes"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
	_ = &acktags.Tags{}
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

	if ackcompare.HasNilDifference(a.ko.Spec.ContainerProvider, b.ko.Spec.ContainerProvider) {
		delta.Add("Spec.ContainerProvider", a.ko.Spec.ContainerProvider, b.ko.Spec.ContainerProvider)
	} else if a.ko.Spec.ContainerProvider != nil && b.ko.Spec.ContainerProvider != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ContainerProvider.ID, b.ko.Spec.ContainerProvider.ID) {
			delta.Add("Spec.ContainerProvider.ID", a.ko.Spec.ContainerProvider.ID, b.ko.Spec.ContainerProvider.ID)
		} else if a.ko.Spec.ContainerProvider.ID != nil && b.ko.Spec.ContainerProvider.ID != nil {
			if *a.ko.Spec.ContainerProvider.ID != *b.ko.Spec.ContainerProvider.ID {
				delta.Add("Spec.ContainerProvider.ID", a.ko.Spec.ContainerProvider.ID, b.ko.Spec.ContainerProvider.ID)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ContainerProvider.Info, b.ko.Spec.ContainerProvider.Info) {
			delta.Add("Spec.ContainerProvider.Info", a.ko.Spec.ContainerProvider.Info, b.ko.Spec.ContainerProvider.Info)
		} else if a.ko.Spec.ContainerProvider.Info != nil && b.ko.Spec.ContainerProvider.Info != nil {
			if ackcompare.HasNilDifference(a.ko.Spec.ContainerProvider.Info.EKSInfo, b.ko.Spec.ContainerProvider.Info.EKSInfo) {
				delta.Add("Spec.ContainerProvider.Info.EKSInfo", a.ko.Spec.ContainerProvider.Info.EKSInfo, b.ko.Spec.ContainerProvider.Info.EKSInfo)
			} else if a.ko.Spec.ContainerProvider.Info.EKSInfo != nil && b.ko.Spec.ContainerProvider.Info.EKSInfo != nil {
				if ackcompare.HasNilDifference(a.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace, b.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace) {
					delta.Add("Spec.ContainerProvider.Info.EKSInfo.Namespace", a.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace, b.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace)
				} else if a.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace != nil && b.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace != nil {
					if *a.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace != *b.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace {
						delta.Add("Spec.ContainerProvider.Info.EKSInfo.Namespace", a.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace, b.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace)
					}
				}
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ContainerProvider.Type, b.ko.Spec.ContainerProvider.Type) {
			delta.Add("Spec.ContainerProvider.Type", a.ko.Spec.ContainerProvider.Type, b.ko.Spec.ContainerProvider.Type)
		} else if a.ko.Spec.ContainerProvider.Type != nil && b.ko.Spec.ContainerProvider.Type != nil {
			if *a.ko.Spec.ContainerProvider.Type != *b.ko.Spec.ContainerProvider.Type {
				delta.Add("Spec.ContainerProvider.Type", a.ko.Spec.ContainerProvider.Type, b.ko.Spec.ContainerProvider.Type)
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
	desiredACKTags, _ := convertToOrderedACKTags(a.ko.Spec.Tags)
	latestACKTags, _ := convertToOrderedACKTags(b.ko.Spec.Tags)
	if !ackcompare.MapStringStringEqual(desiredACKTags, latestACKTags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	}

	return delta
}
