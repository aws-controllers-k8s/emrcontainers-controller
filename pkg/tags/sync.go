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

package tags

import (
	"context"

	svcsdk "github.com/aws/aws-sdk-go-v2/service/emrcontainers"
)

// GetTags returns the tags for a given resource ARN
func GetTags(
	ctx context.Context,
	resourceARN string,
	sdkapi *svcsdk.Client,
) (map[string]*string, error) {
	input := &svcsdk.ListTagsForResourceInput{
		ResourceArn: &resourceARN,
	}
	resp, err := sdkapi.ListTagsForResource(ctx, input)
	if err != nil {
		return nil, err
	}

	// Convert to map[string]*string for compatibility with ACK
	tags := make(map[string]*string)
	for k, v := range resp.Tags {
		value := v // Create a copy to avoid issues with the loop variable
		tags[k] = &value
	}

	return tags, nil
}

// SyncTags synchronizes the tags between the spec and the AWS resource
func SyncTags(
	ctx context.Context,
	resourceARN string,
	desired map[string]*string,
	latest map[string]*string,
	sdkapi *svcsdk.Client,
) error {
	// Check if tags are the same
	if mapsEqual(desired, latest) {
		return nil
	}

	// Determine which tags to add/update
	toAdd := make(map[string]string)
	for k, v := range desired {
		if v != nil {
			latestVal, exists := latest[k]
			if !exists || *v != *latestVal {
				toAdd[k] = *v
			}
		}
	}

	// Determine which tags to remove
	toRemove := []string{}
	for k := range latest {
		if _, exists := desired[k]; !exists {
			toRemove = append(toRemove, k)
		}
	}

	if len(toAdd) > 0 {
		_, err := sdkapi.TagResource(
			ctx,
			&svcsdk.TagResourceInput{
				ResourceArn: &resourceARN,
				Tags:        toAdd,
			},
		)
		if err != nil {
			return err
		}
	}

	if len(toRemove) > 0 {
		_, err := sdkapi.UntagResource(
			ctx,
			&svcsdk.UntagResourceInput{
				ResourceArn: &resourceARN,
				TagKeys:     toRemove,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// mapsEqual compares two string pointer maps for equality
func mapsEqual(a, b map[string]*string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v1 := range a {
		v2, ok := b[k]
		if !ok || v1 == nil || v2 == nil || *v1 != *v2 {
			return false
		}
	}

	return true
}
