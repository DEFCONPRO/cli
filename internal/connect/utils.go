package connect

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/confluentinc/cli/v3/pkg/errors"
)

func getConfig(cmd *cobra.Command) (*map[string]string, error) {
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return nil, err
	}

	options, err := parseConfigFile(configFile)
	if err != nil {
		return nil, fmt.Errorf(errors.UnableToReadConfigurationFileErrorMsg, configFile, err)
	}

	connectorType := options["confluent.connector.type"]
	if connectorType == "" {
		connectorType = "MANAGED"
	}

	_, nameExists := options["name"]
	_, classExists := options["connector.class"]

	if connectorType != "CUSTOM" && (!nameExists || !classExists) {
		return nil, fmt.Errorf(`required configs "name" and "connector.class" missing from connector config file "%s"`, configFile)
	}

	return &options, nil
}

func parseConfigFile(filename string) (map[string]string, error) {
	jsonFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf(errors.UnableToReadConfigurationFileErrorMsg, filename, err)
	}
	if len(jsonFile) == 0 {
		return nil, fmt.Errorf(`connector config file "%s" is empty`, filename)
	}

	kvPairs := make(map[string]string)
	var options map[string]any

	if err := json.Unmarshal(jsonFile, &options); err != nil {
		return nil, fmt.Errorf(errors.UnableToReadConfigurationFileErrorMsg, filename, err)
	}

	for key, val := range options {
		if val2, ok := val.(string); ok {
			kvPairs[key] = val2
		} else {
			// We support object-as-a-value only for "config" key.
			if key != "config" {
				return nil, fmt.Errorf(`only string values are permitted for the configuration "%s"`, key)
			}

			configMap, ok := val.(map[string]any)
			if !ok {
				return nil, fmt.Errorf(`value for the configuration "config" is malformed`)
			}

			for configKey, configVal := range configMap {
				value, isString := configVal.(string)
				if !isString {
					return nil, fmt.Errorf(`only string values are permitted for the configuration "%s"`, configKey)
				}
				kvPairs[configKey] = value
			}
		}
	}

	return kvPairs, err
}
