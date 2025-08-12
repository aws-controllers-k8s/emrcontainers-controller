package job_run

import (
	"testing"

	svcapitypes "github.com/aws-controllers-k8s/emrcontainers-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/emrcontainers/types"
	"github.com/stretchr/testify/assert"
)

func TestJobInCancellableState_NonCancellableStates(t *testing.T) {
	nonCancellableStates := []string{
		string(svcsdktypes.JobRunStateCompleted),
		string(svcsdktypes.JobRunStateFailed),
		string(svcsdktypes.JobRunStateCancelled),
		string(svcsdktypes.JobRunStateCancelPending),
	}

	for _, state := range nonCancellableStates {
		t.Run("state_"+state, func(t *testing.T) {
			r := &resource{
				ko: &svcapitypes.JobRun{
					Status: svcapitypes.JobRunStatus{
						State: aws.String(state),
					},
				},
			}

			// Verify that job is NOT in cancellable state
			assert.False(t, jobInCancellableState(r), "Job in state %s should not be cancellable", state)
		})
	}
}

func TestJobInCancellableState_CancellableStates(t *testing.T) {
	cancellableStates := []string{
		string(svcsdktypes.JobRunStateSubmitted),
		string(svcsdktypes.JobRunStatePending),
		string(svcsdktypes.JobRunStateRunning),
	}

	for _, state := range cancellableStates {
		t.Run("state_"+state, func(t *testing.T) {
			r := &resource{
				ko: &svcapitypes.JobRun{
					Status: svcapitypes.JobRunStatus{
						State: aws.String(state),
					},
				},
			}

			// Verify that job IS in cancellable state
			assert.True(t, jobInCancellableState(r), "Job in state %s should be cancellable", state)
		})
	}
}
