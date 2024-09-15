package parse

import (
	"os"

	jsonIter "github.com/json-iterator/go"
)

// JsonFileToStruct loads json file and parse into given struct
func JsonFileToStruct(jsonPath string, target any) (any, error) {
	content, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	err = ToStruct(content, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

// ToStruct parses json content into given struct
func ToStruct(content []byte, target any) error {
	err := Decode(content, target)
	if err != nil {
		return err
	}

	return nil
}

// ToJson parses given struct into json
func ToJson(target any, indent ...bool) ([]byte, error) {
	content, err := Encode(target, indent...)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// ToJsonString parses given struct into json string
func ToJsonString(target any) (string, error) {
	content, err := ToJson(target)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Decode parses given json into given struct
func Decode(data []byte, v any) error {
	var json = jsonIter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, v)
}

// Encode parses given struct into json
func Encode(v any, indent ...bool) ([]byte, error) {
	var json = jsonIter.ConfigCompatibleWithStandardLibrary
	if len(indent) > 0 && indent[0] {
		return json.MarshalIndent(v, "", "  ")
	}

	return json.Marshal(v)
}
