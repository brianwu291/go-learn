package utils

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var (
	// ErrInvalidJSON is returned when the input is not valid JSON
	ErrInvalidJSON = errors.New("invalid JSON input")
)

// TrimOptions configures the JSON trimming behavior
type TrimOptions struct {
	// Paths to trim with gjson path syntax
	// See: https://github.com/tidwall/gjson#path-syntax
	Paths []string
}

// TrimJSON takes a JSON string and trims whitespace from string values
// at the specified paths. Returns the trimmed JSON and a boolean
// indicating if any modifications were made.
func TrimJSON(jsonStr string, opts TrimOptions) (string, bool, error) {
	// Check if the input is valid JSON
	if !gjson.Valid(jsonStr) {
		return jsonStr, false, ErrInvalidJSON
	}

	// If no paths are specified, return the original JSON unmodified
	if len(opts.Paths) == 0 {
		return jsonStr, false, nil
	}

	return trimPaths(jsonStr, opts.Paths)
}

// TrimJSONBytes is like TrimJSON but works with byte slices
func TrimJSONBytes(data []byte, opts TrimOptions) ([]byte, bool, error) {
	str := string(data)
	result, modified, err := TrimJSON(str, opts)
	if err != nil {
		return nil, false, err
	}
	if !modified {
		return data, false, nil
	}
	return []byte(result), true, nil
}

// UnmarshalAndTrim unmarshals JSON data into a struct while trimming string fields
func UnmarshalAndTrim(data []byte, v interface{}, opts TrimOptions) error {
	trimmed, _, err := TrimJSONBytes(data, opts)
	if err != nil {
		return err
	}

	// Use standard json unmarshaler to fill the struct
	return json.Unmarshal(trimmed, v)
}

// ReadAndTrim reads from an io.Reader, trims the JSON using gjson/sjson,
// and unmarshals it into the provided value.
func ReadAndTrim(r io.Reader, v interface{}, opts TrimOptions) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return UnmarshalAndTrim(data, v, opts)
}

// trimPaths trims string values at specific paths using gjson/sjson
func trimPaths(jsonStr string, paths []string) (string, bool, error) {
	result := jsonStr
	modified := false

	for _, path := range paths {
		// Get the value at the path
		value := gjson.Get(result, path)

		if value.Exists() && value.IsArray() {
			for _, val := range value.Array() {
				res, hasChanged := trim(result, val, path)
				if hasChanged {
					modified = true
					result = res
				}

			}
			continue
		}

		res, hasChanged := trim(result, value, path)
		if hasChanged {
			modified = true
			result = res
		}
	}

	return result, modified, nil
}

func trim(jsonStr string, unTrimmedVal gjson.Result, path string) (string, bool) {
	if !unTrimmedVal.Exists() || unTrimmedVal.Type != gjson.String {
		return jsonStr, false
	}
	strValue := unTrimmedVal.String()
	trimmed := strings.TrimSpace(strValue)
	if trimmed != strValue {
		var err error
		result, err := sjson.Set(jsonStr, path, trimmed)
		if err != nil {
			return jsonStr, false
		}
		return result, true
	}
	return jsonStr, false
}
