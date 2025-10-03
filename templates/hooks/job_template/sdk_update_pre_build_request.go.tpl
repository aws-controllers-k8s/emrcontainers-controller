// Sync tags if they've changed
if delta.DifferentAt("Spec.Tags") {
    err := rm.syncTags(
        ctx,
        latest,
        desired,
    )
    if err != nil {
        return nil, err
    }
}
// If the only difference is in the tags, we don't need to make an update call
if !delta.DifferentExcept("Spec.Tags") {
    return desired, nil
}