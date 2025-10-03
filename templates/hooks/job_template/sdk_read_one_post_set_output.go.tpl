// Retrieve and set the tags for the JobTemplate resource
if ko.Status.ACKResourceMetadata != nil && ko.Status.ACKResourceMetadata.ARN != nil {
    ko.Spec.Tags = rm.getTags(ctx, string(*ko.Status.ACKResourceMetadata.ARN))
}