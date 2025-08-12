
// If the job has already reached a state where it has finished we can 
// consider the delete operation completed.
if err != nil {
    var awsErr smithy.APIError
    if errors.As(err, &awsErr) && awsErr.ErrorCode() == "ValidationException" && strings.HasSuffix(awsErr.ErrorMessage(), "is not in a cancellable state") {
        rm.log.Info("JobRun is not in a cancellable state. Allowing deletion of resource continue.", "JobRun", r.ko.Status.ID)
        return nil, nil
    }
}