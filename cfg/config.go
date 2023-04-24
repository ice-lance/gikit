package cfg

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func CheckConfig(path string, cfg []byte) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		f.Write(cfg)
		f.Close()

	}
	return nil
}

func Read(path string, obj any) error {
	pathSilce := strings.Split(path, ".")
	n := len(pathSilce)
	fileSuffix := pathSilce[n-1]
	switch fileSuffix {
	case "json":
		return ReadConfigFromJson(path, obj)
	case "yaml":
	case "yml":
		return ReadConfigFromYaml(path, obj)
	case "toml":
		return ReadConfigFromToml(path, obj)
	default:
		return errors.New(" Unsupported format! ")
	}
	return nil
}

func ReadConfigFromJson(path string, obj any) error {
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

func ReadConfigFromToml(path string, obj any) error {

	return nil
}

func ReadConfigFromYaml(path string, obj any) error {
	return nil
}
