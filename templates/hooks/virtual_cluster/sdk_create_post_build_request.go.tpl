
    if input.ContainerProvider != nil {
        // Clear any existing Info
        if input.ContainerProvider.Info != nil {
            input.ContainerProvider.Info = nil
        }

        // Set the Info field if it exists in the spec
        eksInfo := &svcsdktypes.EksInfo{}
        if desired.ko.Spec.ContainerProvider.Info != nil && 
           desired.ko.Spec.ContainerProvider.Info.EKSInfo != nil && 
           desired.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace != nil {
            eksInfo.Namespace = desired.ko.Spec.ContainerProvider.Info.EKSInfo.Namespace
        } else {
            eksInfo.Namespace = aws.String("default")
        }
        input.ContainerProvider.Info = &svcsdktypes.ContainerInfoMemberEksInfo{
            Value: *eksInfo,
        }
    }