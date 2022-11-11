		// DescribeJobRun should output ConfigurationOverrides and show all available configuration
		if resp.JobRun.ConfigurationOverrides != nil {
			ko.Spec.ConfigurationOverrides, err = configurationOverridesToString(resp.JobRun.ConfigurationOverrides)
			if err != nil {
				return nil, err
			}
		}
