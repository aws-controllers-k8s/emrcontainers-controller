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
	"maps"
	"slices"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
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

	from := toACKTags(latestTags)
	to := toACKTags(desiredTags)
	added, _, removed := ackcompare.GetTagsDifference(from, to)

	// Nil-valued desired keys are absent from `to` (toACKTags skips them), so
	// GetTagsDifference may place them in `removed` if they exist in `from`.
	// Delete only nil-valued desired keys from `removed` to preserve the
	// previous computeTagsDelta semantics of never untagging a nil-valued key.
	for k, v := range desiredTags {
		if v == nil {
			delete(removed, k)
		}
	}

	if len(added) == 0 && len(removed) == 0 {
		return nil
	}

	if len(removed) > 0 {
		_, err = client.UntagResource(
			ctx,
			&svcsdk.UntagResourceInput{
				ResourceArn: aws.String(resourceARN),
				TagKeys:     slices.Sorted(maps.Keys(removed)),
			},
		)
		mr.RecordAPICall("UPDATE", "UntagResource", err)
		if err != nil {
			return fmt.Errorf("UntagResource Error: %w", err)
		}
	}

	if len(added) > 0 {
		_, err = client.TagResource(
			ctx,
			&svcsdk.TagResourceInput{
				ResourceArn: aws.String(resourceARN),
				Tags:        added,
			},
		)
		mr.RecordAPICall("UPDATE", "TagResource", err)
		if err != nil {
			return fmt.Errorf("TagResource Error: %w", err)
		}
	}
	return nil
}

// toACKTags converts an EMR Containers tag map into the ACK runtime tags type.
// Entries with nil values are skipped, since EMR Containers rejects empty tag values.
func toACKTags(tags map[string]*string) acktags.Tags {
	m := make(acktags.Tags)
	for k, v := range tags {
		if v == nil {
			continue
		}
		m[k] = *v
	}
	return m
}
