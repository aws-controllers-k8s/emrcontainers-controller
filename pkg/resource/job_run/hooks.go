package job_run

import (
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/emrcontainers/types"
	"github.com/ghodss/yaml"
)

func configurationOverridesToString(cfg *svcsdktypes.ConfigurationOverrides) (*string, error) {
	configBytes, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	configStr := string(configBytes)
	return &configStr, nil
}

func stringToConfigurationOverrides(cfg *string) (*svcsdktypes.ConfigurationOverrides, error) {
	if cfg == nil {
		cfg = aws.String("")
	}

	var config svcsdktypes.ConfigurationOverrides
	err := yaml.Unmarshal([]byte(*cfg), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func customPreCompare(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {
	aConfig, err := stringToConfigurationOverrides(a.ko.Spec.ConfigurationOverrides)
	if err != nil {
		panic(err)
	}
	bConfig, err := stringToConfigurationOverrides(b.ko.Spec.ConfigurationOverrides)
	if err != nil {
		panic(err)
	}

	// background:
	// API Always return a non empty configuration
	// Users can give empty configuration, and API still might return something
	// Users can give empty fields and API still might return something

	// If we have a nil configuration and API returns a non nil configuration - something is wrong
	if ackcompare.HasNilDifference(aConfig.MonitoringConfiguration, bConfig.MonitoringConfiguration) {
		delta.Add("Spec.ConfigurationOverrides", aConfig.MonitoringConfiguration, bConfig.MonitoringConfiguration)
	} else if aConfig.MonitoringConfiguration != nil && bConfig.MonitoringConfiguration != nil {
		if ackcompare.HasNilDifference(aConfig.MonitoringConfiguration.PersistentAppUI, bConfig.MonitoringConfiguration.PersistentAppUI) {
			if aConfig.MonitoringConfiguration.PersistentAppUI == "" && bConfig.MonitoringConfiguration.PersistentAppUI == svcsdktypes.PersistentAppUIEnabled {
				// We do not consider this as a difference because the API defaults PersistentAppUI to "ENABLED"
			} else {
				delta.Add("Spec.ConfigurationOverrides.PersistentAppUI", aConfig.MonitoringConfiguration.PersistentAppUI, bConfig.MonitoringConfiguration.PersistentAppUI)
			}
		} else if aConfig.MonitoringConfiguration.PersistentAppUI != bConfig.MonitoringConfiguration.PersistentAppUI {
			delta.Add("Spec.ConfigurationOverrides.PersistentAppUI", aConfig.MonitoringConfiguration.PersistentAppUI, bConfig.MonitoringConfiguration.PersistentAppUI)
		}
		if ackcompare.HasNilDifference(
			aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration,
			bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration,
		) {
			delta.Add("Spec.ConfigurationOverrides.CloudWatchMonitoringConfiguration", aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration, bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration)
		} else if aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration != nil &&
			bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration != nil {

			if ackcompare.HasNilDifference(
				aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName,
				bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName,
			) {
				delta.Add("Spec.ConfigurationOverrides.CloudWatchMonitoringConfiguration.LogGroupName", aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName, bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName)
			} else if *aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName !=
				*bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName {
				delta.Add("Spec.ConfigurationOverrides.CloudWatchMonitoringConfiguration.LogGroupName", aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName, bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogGroupName)
			}
			if ackcompare.HasNilDifference(
				aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix,
				bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix,
			) {
				delta.Add("Spec.ConfigurationOverrides.CloudWatchMonitoringConfiguration.LogStreamNamePrefix", aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix, bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix)
			} else if *aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix !=
				*bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix {
				delta.Add("Spec.ConfigurationOverrides.CloudWatchMonitoringConfiguration.LogStreamNamePrefix", aConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix, bConfig.MonitoringConfiguration.CloudWatchMonitoringConfiguration.LogStreamNamePrefix)
			}
		}

		//

		if ackcompare.HasNilDifference(
			aConfig.MonitoringConfiguration.S3MonitoringConfiguration,
			bConfig.MonitoringConfiguration.S3MonitoringConfiguration,
		) {
			delta.Add("Spec.ConfigurationOverrides.S3MonitoringConfiguration", aConfig.MonitoringConfiguration.S3MonitoringConfiguration, bConfig.MonitoringConfiguration.S3MonitoringConfiguration)
		} else if aConfig.MonitoringConfiguration.S3MonitoringConfiguration != nil &&
			bConfig.MonitoringConfiguration.S3MonitoringConfiguration != nil {
			if ackcompare.HasNilDifference(
				aConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri,
				bConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri,
			) {
				delta.Add("Spec.ConfigurationOverrides.S3MonitoringConfiguration.LogUri", aConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri, bConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri)
			} else if *aConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri !=
				*bConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri {
				delta.Add("Spec.ConfigurationOverrides.S3MonitoringConfiguration.LogUri", aConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri, bConfig.MonitoringConfiguration.S3MonitoringConfiguration.LogUri)
			}
		}

	}

	// If two arrays have different sizes then automatically they are different
	if len(aConfig.ApplicationConfiguration) != len(bConfig.ApplicationConfiguration) {
		delta.Add("Spec.ApplicationConfiguration", aConfig.ApplicationConfiguration, bConfig.ApplicationConfiguration)
	} else if len(aConfig.ApplicationConfiguration) > 0 {
		// at this stage we know that they have the same size they contain at least one element
		// We assume that the EMRContainer API doesn't mess with the order of the provided Application
		// Configuration (To verify).
		if !reflect.DeepEqual(aConfig.ApplicationConfiguration, bConfig.ApplicationConfiguration) {
			delta.Add("Spec.ApplicationConfiguration", aConfig.ApplicationConfiguration, bConfig.ApplicationConfiguration)
		}
	}
}
