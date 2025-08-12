package job_run

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/emrcontainers"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/emrcontainers/types"
	smithy "github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	svcapitypes "github.com/aws-controllers-k8s/emrcontainers-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/runtime/pkg/metrics"
)

// resourceManager directly uses *svcsdk.Client struct instead of interface. In order to mock HTTP requests
// we instead mock the underlying HttpClient.
type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(request *http.Request) (*http.Response, error) {
	args := m.Called(request)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Validate that sdkDelete exits without calling CancelJobRun when JobRun is in a
// state that isn't cancellable.
func TestSdkDelete_NonCancellableStates_ReturnsEarlyWithoutError(t *testing.T) {
	noncancellableStates := []string{
		string(svcsdktypes.JobRunStateCompleted),
		string(svcsdktypes.JobRunStateFailed),
		string(svcsdktypes.JobRunStateCancelled),
		string(svcsdktypes.JobRunStateCancelPending),
	}

	for _, state := range noncancellableStates {
		t.Run("state_"+state, func(t *testing.T) {
			r := &resource{
				ko: &svcapitypes.JobRun{
					Status: svcapitypes.JobRunStatus{
						State: aws.String(state),
						ID:    aws.String("fake-id"),
					},
					Spec: svcapitypes.JobRunSpec{
						VirtualClusterID: aws.String("test-cluster-id"),
					},
				},
			}

			mockHttpClient := &MockHttpClient{}
			emrClient := svcsdk.New(svcsdk.Options{
				HTTPClient: mockHttpClient,
				Region:     "no-region",
			})

			mockHttpClient.On("Do", mock.Anything).Return(nil, errors.New("Bad Request"))
			rm := &resourceManager{
				sdkapi:  emrClient,
				metrics: metrics.NewMetrics("test-emr"),
			}

			// Execute - should return early without calling API
			latest, err := rm.sdkDelete(context.Background(), r)

			mockHttpClient.AssertNotCalled(t, "Do", mock.Anything)
			assert.NoError(t, err)
			assert.Nil(t, latest)
		})
	}
}

// Validate that when JobRun is in a cancellable state, sdkDelete calls CancelJobRun.
func TestSdkDelete_CancellableStates_CallsCancelJobRun(t *testing.T) {
	cancellableStates := []string{
		string(svcsdktypes.JobRunStatePending),
		string(svcsdktypes.JobRunStateRunning),
		string(svcsdktypes.JobRunStateSubmitted),
	}

	for _, state := range cancellableStates {
		t.Run("state_"+state, func(t *testing.T) {
			r := &resource{
				ko: &svcapitypes.JobRun{
					Status: svcapitypes.JobRunStatus{
						State: aws.String(state),
						ID:    aws.String("fake-id"),
					},
					Spec: svcapitypes.JobRunSpec{
						VirtualClusterID: aws.String("test-cluster-id"),
					},
				},
			}

			mockHttpClient := &MockHttpClient{}
			output := svcsdk.CancelJobRunOutput{
				Id:               r.ko.Status.ID,
				VirtualClusterId: r.ko.Spec.VirtualClusterID,
			}
			jsonBytes, _ := json.Marshal(output)
			mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(string(jsonBytes))),
			}, nil)

			emrClient := svcsdk.New(svcsdk.Options{
				HTTPClient: mockHttpClient,
				Region:     "no-region",
			})
			rm := &resourceManager{
				sdkapi:  emrClient,
				metrics: metrics.NewMetrics("test-emr"),
			}

			latest, err := rm.sdkDelete(context.Background(), r)

			mockHttpClient.AssertCalled(t, "Do", mock.Anything)
			assert.NoError(t, err)
			assert.Nil(t, latest)
		})
	}
}

// Validate that in the event CancelJobRun returns a ValidationException with "Job run X is not in a cancellable state"
// sdkDelete does not block resource finalization.
func TestSdkDelete_ValidationErrorWithNoncancellableState_AllowsFinalization(t *testing.T) {
	t.Run("TestSdkDelete_ValidationErrorWithNoncancellableState_AllowsFinalization", func(t *testing.T) {
		// Create resource with cancellable state
		state := string(svcsdktypes.JobRunStateSubmitted)
		r := &resource{
			ko: &svcapitypes.JobRun{
				Status: svcapitypes.JobRunStatus{
					State: aws.String(state),
					ID:    aws.String("fake-id"),
				},
				Spec: svcapitypes.JobRunSpec{
					VirtualClusterID: aws.String("test-cluster-id"),
				},
			},
		}

		mockHttpClient := &MockHttpClient{}
		output := &smithy.GenericAPIError{
			Message: "Job run 13221 is not in a cancellable state",
			Code:    "ValidationException",
			Fault:   400,
		}
		jsonBytes, _ := json.Marshal(output)
		mockHttpClient.On("Do", mock.Anything).Return(&http.Response{
			StatusCode: 400,
			Status:     "400 Bad Request",
			Body:       io.NopCloser(strings.NewReader(string(jsonBytes))),
		}, nil)

		emrClient := svcsdk.New(svcsdk.Options{
			HTTPClient: mockHttpClient,
			Region:     "no-region",
		})
		rm := &resourceManager{
			sdkapi:  emrClient,
			metrics: metrics.NewMetrics("test-emr"),
		}

		latest, err := rm.sdkDelete(context.Background(), r)

		mockHttpClient.AssertCalled(t, "Do", mock.Anything)
		assert.NoError(t, err)
		assert.Nil(t, latest)
	})
}
