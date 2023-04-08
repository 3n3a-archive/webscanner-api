package file

import (
	"os"
	"gopkg.in/yaml.v3"
	validate "github.com/3n3a/webscanner-api/modules/validation"
)

// Reads YAML-File into specified struct
func ReadYAMLIntoStruct[T any](filepath string) (T, error) {
	var element T

	var data []byte
	data, err := os.ReadFile(filepath)
	if validate.IsErrorState(err) {
		return element, err
	}

	if err := yaml.Unmarshal(data, &element); validate.IsErrorState(err) {
		return element, err
	}

	return element, nil
} 