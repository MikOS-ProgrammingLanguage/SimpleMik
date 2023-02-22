package config

import (
	"MicNewDawn/errors"
	"MicNewDawn/mik_compiler/optimizer"
	"encoding/json"
	"os"
)

// Check wether file exists
// Read contents from file
// Unmarshal them as JSON
// return the resulting struct
func ParseOptimizationConfig(file *string) optimizer.OptimizationFlags {
	contents, err := os.ReadFile(*file)
	if err != nil {
		errors.CouldNotOpenOrLocateFile(*file)
	}

	returnStruct := optimizer.OptimizationFlags{}
	jsonErr := json.Unmarshal(contents, &returnStruct)
	if jsonErr != nil {
		errors.ThrowError(jsonErr.Error(), true)
	}

	return returnStruct
}
