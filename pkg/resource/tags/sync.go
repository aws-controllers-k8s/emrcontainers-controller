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
	"fmt"

	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/emrcontainers"
)

type metricsRecorder interface {
	RecordAPICall(opType string, opID string, err error)
}

type tagsClient interface {
	TagResource(context.Context, *svcsdk.TagResourceInput, ...func(*svcsdk.Options)) (*svcsdk.TagResourceOutput, error)
	ListTagsForResource(context.Context, *svcsdk.ListTagsForResourceInput, ...func(*svcsdk.Options)) (*svcsdk.ListTagsForResourceOutput, error)
	UntagResource(context.Context, *svcsdk.UntagResourceInput, ...func(*svcsdk.Options)) (*svcsdk.UntagResourceOutput, error)
}

// SyncResourceTags uses TagResource and UntagResource API Calls to add, remove
// and update resource tags.
func SyncResourceTags(
	ctx context.Context,
	client tagsClient,
	mr metricsRecorder,
	resourceARN string,
	desiredTags map[string]*string,
	latestTags map[string]*string,
) error {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("common.SyncResourceTags")
	defer func() {
		exit(err)
	}()

	addedOrUpdated, removed := computeTagsDelta(desiredTags, latestTags)

	if len(removed) > 0 {
		_, err = client.UntagResource(
			ctx,
			&svcsdk.UntagResourceInput{
				ResourceArn: aws.String(resourceARN),
				TagKeys:     removed,
			},
		)
		mr.RecordAPICall("UPDATE", "UntagResource", err)
		if err != nil {
			return fmt.Errorf("UntagResource Error: %w", err)
		}
	}

	if len(addedOrUpdated) > 0 {
		_, err = client.TagResource(
			ctx,
			&svcsdk.TagResourceInput{
				ResourceArn: aws.String(resourceARN),
				Tags:        addedOrUpdated,
			},
		)
		mr.RecordAPICall("UPDATE", "TagResource", err)
		if err != nil {
			return fmt.Errorf("TagResource Error: %w", err)
		}
	}
	return nil
}

// computeTagsDelta compares two Tag maps and return two different list
// containing the addedOrupdated and removed tags. The removed tags array
// only contains the tags Keys.
func computeTagsDelta(
	a map[string]*string,
	b map[string]*string,
) (addedOrUpdated map[string]string, removed []string) {

	// Find the keys in the Spec have either been added or updated.
	addedOrUpdated = make(map[string]string)
	for aKey, aValue := range a {
		if bValue, exists := b[aKey]; !exists || (aValue != nil && bValue != nil && *aValue != *bValue) {
			if aValue != nil {
				addedOrUpdated[aKey] = *aValue
			}
		}
	}

	for bKey := range b {
		if _, exists := a[bKey]; !exists {
			removed = append(removed, bKey)
		}
	}

	return addedOrUpdated, removed
}
