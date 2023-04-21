package cfg

import (
	"encoding/json"
	"errors"
	"os"
)

func ReadConfigFromJson(path string, obj interface{}) error {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if len(configFile) > 0 {
		err = json.Unmarshal(configFile, obj)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New(" empty config ! ")
}

func ReadConfigFromToml() error {

	return nil
}

func ReadConfigFromYaml() error {
	return nil
}
