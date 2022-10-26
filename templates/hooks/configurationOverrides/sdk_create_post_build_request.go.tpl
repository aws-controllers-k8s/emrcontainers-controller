    // Unmarshall ConfigurationOverrides structure
    if desired.ko.Spec.ConfigurationOverrides != nil {
        input.ConfigurationOverrides, err = stringToConfigurationOverrides(desired.ko.Spec.ConfigurationOverrides)
        if err != nil {
          return nil, err
        }
    }
