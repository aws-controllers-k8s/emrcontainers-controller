// Add tags to the JobTemplate resource after creation
if ko.Spec.Tags != nil {
    // Mark the resource as not synced to trigger a requeue and apply tags
    rm.setStatusSynced(ctx, &resource{ko}, nil, nil, false)
}