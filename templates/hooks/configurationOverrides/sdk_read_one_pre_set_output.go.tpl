		// DescribeJobRun should output ConfigurationOverrides and show all available configuration
		if resp.JobRun.ConfigurationOverrides != nil {
			ko.Spec.ConfigurationOverrides, err = cfgToString(resp.JobRun.ConfigurationOverrides)
			if err != nil {
				return nil, err
			}
		}
