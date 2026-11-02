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
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/emrcontainers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTagsClient struct {
	mock.Mock
}

func (m *mockTagsClient) TagResource(ctx context.Context, input *svcsdk.TagResourceInput, opts ...func(*svcsdk.Options)) (*svcsdk.TagResourceOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*svcsdk.TagResourceOutput), args.Error(1)
}

func (m *mockTagsClient) ListTagsForResource(ctx context.Context, input *svcsdk.ListTagsForResourceInput, opts ...func(*svcsdk.Options)) (*svcsdk.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*svcsdk.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockTagsClient) UntagResource(ctx context.Context, input *svcsdk.UntagResourceInput, opts ...func(*svcsdk.Options)) (*svcsdk.UntagResourceOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*svcsdk.UntagResourceOutput), args.Error(1)
}

type mockMetricsRecorder struct {
	mock.Mock
}

func (m *mockMetricsRecorder) RecordAPICall(opType string, opID string, err error) {
	m.Called(opType, opID, err)
}

func TestSyncResourceTags(t *testing.T) {
	ctx := context.Background()
	resourceARN := "arn:aws:emr-containers:us-west-2:123456789012:virtualclusters/test"

	tests := []struct {
		name        string
		desired     map[string]*string
		latest      map[string]*string
		expectTag   bool
		expectUntag bool
	}{
		{
			name: "add new tags",
			desired: map[string]*string{
				"key1": aws.String("value1"),
				"key2": aws.String("value2"),
			},
			latest:    map[string]*string{},
			expectTag: true,
		},
		{
			name:    "remove tags",
			desired: map[string]*string{},
			latest: map[string]*string{
				"key1": aws.String("value1"),
			},
			expectUntag: true,
		},
		{
			name: "update existing tags",
			desired: map[string]*string{
				"key1": aws.String("newvalue"),
			},
			latest: map[string]*string{
				"key1": aws.String("oldvalue"),
			},
			expectTag: true,
		},
		{
			name: "no changes",
			desired: map[string]*string{
				"key1": aws.String("value1"),
			},
			latest: map[string]*string{
				"key1": aws.String("value1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockTagsClient{}
			mr := &mockMetricsRecorder{}

			if tt.expectTag {
				client.On("TagResource", ctx, mock.AnythingOfType("*emrcontainers.TagResourceInput")).Return(&svcsdk.TagResourceOutput{}, nil)
				mr.On("RecordAPICall", "UPDATE", "TagResource", nil)
			}

			if tt.expectUntag {
				client.On("UntagResource", ctx, mock.AnythingOfType("*emrcontainers.UntagResourceInput")).Return(&svcsdk.UntagResourceOutput{}, nil)
				mr.On("RecordAPICall", "UPDATE", "UntagResource", nil)
			}

			err := SyncResourceTags(ctx, client, mr, resourceARN, tt.desired, tt.latest)
			assert.NoError(t, err)

			client.AssertExpectations(t)
			mr.AssertExpectations(t)
		})
	}
}

func TestSyncResourceTagsError(t *testing.T) {
	ctx := context.Background()
	resourceARN := "arn:aws:emr-containers:us-west-2:123456789012:virtualclusters/test"

	client := &mockTagsClient{}
	mr := &mockMetricsRecorder{}

	desired := map[string]*string{"key1": aws.String("value1")}
	latest := map[string]*string{}

	apiError := errors.New("tag error")
	client.On("TagResource", ctx, mock.AnythingOfType("*emrcontainers.TagResourceInput")).Return(&svcsdk.TagResourceOutput{}, apiError)
	mr.On("RecordAPICall", "UPDATE", "TagResource", apiError)
	expectedErr := fmt.Errorf("TagResource Error: %w", apiError)

	err := SyncResourceTags(ctx, client, mr, resourceARN, desired, latest)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
func TestSyncResourceTagsInputs(t *testing.T) {
	ctx := context.Background()
	resourceARN := "arn:aws:emr-containers:us-west-2:123456789012:virtualclusters/test"

	tests := []struct {
		name            string
		desired         map[string]*string
		latest          map[string]*string
		expectedTag     map[string]string
		expectedUntagKs []string
	}{
		{
			name: "no changes",
			desired: map[string]*string{
				"key1": aws.String("value1"),
			},
			latest: map[string]*string{
				"key1": aws.String("value1"),
			},
		},
		{
			name: "add tags",
			desired: map[string]*string{
				"key1": aws.String("value1"),
				"key2": aws.String("value2"),
			},
			latest: map[string]*string{},
			expectedTag: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:    "remove tags",
			desired: map[string]*string{},
			latest: map[string]*string{
				"key1": aws.String("value1"),
				"key2": aws.String("value2"),
			},
			// Keys must be in lexicographic order: SyncResourceTags calls slices.Sorted(maps.Keys(removed)).
			expectedUntagKs: []string{"key1", "key2"},
		},
		{
			name: "value change only adds the updated tag",
			desired: map[string]*string{
				"key1": aws.String("newvalue"),
				"key2": aws.String("value2"),
			},
			latest: map[string]*string{
				"key1": aws.String("oldvalue"),
				"key2": aws.String("value2"),
			},
			expectedTag: map[string]string{
				"key1": "newvalue",
			},
		},
		{
			name:            "add and remove simultaneously",
			desired:         map[string]*string{"key2": aws.String("v2")},
			latest:          map[string]*string{"key1": aws.String("v1")},
			expectedTag:     map[string]string{"key2": "v2"},
			expectedUntagKs: []string{"key1"},
		},
		{
			name: "nil desired value is not sent",
			desired: map[string]*string{
				"key1": nil,
			},
			latest: map[string]*string{},
		},
		{
			name: "nil desired value does not untag existing tag",
			desired: map[string]*string{
				"key1": nil,
			},
			latest: map[string]*string{
				"key1": aws.String("value1"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockTagsClient{}
			mr := &mockMetricsRecorder{}

			var gotTag *svcsdk.TagResourceInput
			var gotUntag *svcsdk.UntagResourceInput

			if len(tt.expectedTag) > 0 {
				client.On("TagResource", ctx, mock.AnythingOfType("*emrcontainers.TagResourceInput")).Run(func(args mock.Arguments) {
					gotTag = args.Get(1).(*svcsdk.TagResourceInput)
				}).Return(&svcsdk.TagResourceOutput{}, nil)
				mr.On("RecordAPICall", "UPDATE", "TagResource", nil)
			}

			if len(tt.expectedUntagKs) > 0 {
				client.On("UntagResource", ctx, mock.AnythingOfType("*emrcontainers.UntagResourceInput")).Run(func(args mock.Arguments) {
					gotUntag = args.Get(1).(*svcsdk.UntagResourceInput)
				}).Return(&svcsdk.UntagResourceOutput{}, nil)
				mr.On("RecordAPICall", "UPDATE", "UntagResource", nil)
			}

			err := SyncResourceTags(ctx, client, mr, resourceARN, tt.desired, tt.latest)
			assert.NoError(t, err)

			client.AssertExpectations(t)
			mr.AssertExpectations(t)
			if len(tt.expectedTag) > 0 {
				require.NotNil(t, gotTag)
				assert.Equal(t, resourceARN, aws.ToString(gotTag.ResourceArn))
				assert.Equal(t, tt.expectedTag, gotTag.Tags)
			}
			if len(tt.expectedUntagKs) > 0 {
				require.NotNil(t, gotUntag)
				assert.Equal(t, resourceARN, aws.ToString(gotUntag.ResourceArn))
				assert.Equal(t, tt.expectedUntagKs, gotUntag.TagKeys)
			}
			if len(tt.expectedTag) == 0 && len(tt.expectedUntagKs) == 0 {
				client.AssertNotCalled(t, "TagResource", mock.Anything, mock.Anything)
				client.AssertNotCalled(t, "UntagResource", mock.Anything, mock.Anything)
			}
		})
	}
}
