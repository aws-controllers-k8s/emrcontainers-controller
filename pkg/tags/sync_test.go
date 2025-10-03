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
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/emrcontainers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockSDKAPI struct {
	mock.Mock
}

func (m *mockSDKAPI) ListTagsForResourceWithContext(ctx context.Context, input *svcsdk.ListTagsForResourceInput, opts ...interface{}) (*svcsdk.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*svcsdk.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockSDKAPI) TagResourceWithContext(ctx context.Context, input *svcsdk.TagResourceInput, opts ...interface{}) (*svcsdk.TagResourceOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*svcsdk.TagResourceOutput), args.Error(1)
}

func (m *mockSDKAPI) UntagResourceWithContext(ctx context.Context, input *svcsdk.UntagResourceInput, opts ...interface{}) (*svcsdk.UntagResourceOutput, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*svcsdk.UntagResourceOutput), args.Error(1)
}

func TestGetTags(t *testing.T) {
	ctx := context.Background()
	resourceARN := "arn:aws:emr-containers:us-west-2:123456789012:jobtemplate/test-job-template"

	mockAPI := &mockSDKAPI{}
	mockAPI.On("ListTagsForResourceWithContext", ctx, &svcsdk.ListTagsForResourceInput{
		ResourceArn: aws.String(resourceARN),
	}).Return(&svcsdk.ListTagsForResourceOutput{
		Tags: map[string]*string{
			"key1": aws.String("value1"),
			"key2": aws.String("value2"),
		},
	}, nil)

	tags, err := GetTags(ctx, resourceARN, mockAPI)

	assert.NoError(t, err)
	assert.Len(t, tags, 2)
	assert.Equal(t, "value1", aws.StringValue(tags["key1"]))
	assert.Equal(t, "value2", aws.StringValue(tags["key2"]))

	mockAPI.AssertExpectations(t)
}

func TestSyncTags(t *testing.T) {
	ctx := context.Background()
	resourceARN := "arn:aws:emr-containers:us-west-2:123456789012:jobtemplate/test-job-template"

	// Test adding and updating tags
	mockAPI := &mockSDKAPI{}
	mockAPI.On("TagResourceWithContext", ctx, mock.MatchedBy(func(input *svcsdk.TagResourceInput) bool {
		return aws.StringValue(input.ResourceArn) == resourceARN &&
			aws.StringValue(input.Tags["key1"]) == "new-value1" &&
			aws.StringValue(input.Tags["key3"]) == "value3"
	})).Return(&svcsdk.TagResourceOutput{}, nil)

	mockAPI.On("UntagResourceWithContext", ctx, mock.MatchedBy(func(input *svcsdk.UntagResourceInput) bool {
		return aws.StringValue(input.ResourceArn) == resourceARN &&
			aws.StringValue(input.TagKeys[0]) == "key2"
	})).Return(&svcsdk.UntagResourceOutput{}, nil)

	desired := map[string]*string{
		"key1": aws.String("new-value1"),
		"key3": aws.String("value3"),
	}

	latest := map[string]*string{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
	}

	err := SyncTags(ctx, resourceARN, desired, latest, mockAPI)

	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)

	// Test with no changes
	mockAPI = &mockSDKAPI{}

	desired = map[string]*string{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
	}

	latest = map[string]*string{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
	}

	err = SyncTags(ctx, resourceARN, desired, latest, mockAPI)

	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
}

func TestMapsEqual(t *testing.T) {
	// Test equal maps
	map1 := map[string]*string{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
	}
	map2 := map[string]*string{
		"key1": aws.String("value1"),
		"key2": aws.String("value2"),
	}
	assert.True(t, mapsEqual(map1, map2))

	// Test different values
	map3 := map[string]*string{
		"key1": aws.String("different"),
		"key2": aws.String("value2"),
	}
	assert.False(t, mapsEqual(map1, map3))

	// Test different keys
	map4 := map[string]*string{
		"key1": aws.String("value1"),
		"key3": aws.String("value3"),
	}
	assert.False(t, mapsEqual(map1, map4))

	// Test different lengths
	map5 := map[string]*string{
		"key1": aws.String("value1"),
	}
	assert.False(t, mapsEqual(map1, map5))
}
