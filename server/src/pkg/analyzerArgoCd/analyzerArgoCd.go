package analyzerArgoCd

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"sort"
)

func AnalyzeArgoCd(jsonData string) (string, error) {
	// Extract the manifest field
	manifestStr, err := jsonparser.GetString([]byte(jsonData), "manifest")
	if err != nil {
		return "", fmt.Errorf("manifest field is missing or not a string: %v", err)
	}

	// Extract the spec field from the manifest
	spec, _, _, err := jsonparser.Get([]byte(manifestStr), "spec")
	if err != nil {
		return "", fmt.Errorf("spec field is missing or not a map: %v", err)
	}

	// Extract the strategy field from the spec
	strategy, _, _, err := jsonparser.Get(spec, "strategy")
	if err != nil {
		return "", fmt.Errorf("strategy field is missing or not a map: %v", err)
	}

	// Extract the canary field from the strategy
	canary, _, _, err := jsonparser.Get(strategy, "canary")
	if err != nil {
		return "", fmt.Errorf("canary field is missing or not a map: %v", err)
	}

	// Extract the status field from the manifest
	status, _, _, err := jsonparser.Get([]byte(manifestStr), "status")
	if err != nil {
		return "", fmt.Errorf("status field is missing or not a map: %v", err)
	}

	// Extract the steps field from the canary
	steps, _, _, err := jsonparser.Get(canary, "steps")
	if err != nil {
		return "", fmt.Errorf("steps field is missing or not a list: %v", err)
	}

	// Extract the currentStepIndex from the status
	currentStepIndex, err := jsonparser.GetInt(status, "currentStepIndex")
	if err != nil {
		return "", fmt.Errorf("currentStepIndex field is missing or not an integer: %v", err)
	}

	// Handle the case where currentStepIndex is -1
	if currentStepIndex == -1 {
		return "currentStepIndex is -1", nil
	}

	// Check if currentStepIndex is within the bounds of the steps array
	stepsCount := 0
	_, err = jsonparser.ArrayEach(steps, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		stepsCount++
	})
	if err != nil {
		return "", fmt.Errorf("error counting steps: %v", err)
	}

	if currentStepIndex < 0 || currentStepIndex >= int64(stepsCount) {
		//return "", fmt.Errorf("currentStepIndex %d is out of bounds", currentStepIndex)
	} else {

		// Find the latest step and the setWeight value from the previous step
		var latestStep, setWeight string
		index := 0

		// Sort steps by date-time order
		sortedSteps := make([][]byte, 0)
		_, err = jsonparser.ArrayEach(steps, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			sortedSteps = append(sortedSteps, value)
		})
		if err != nil {
			return "", fmt.Errorf("error iterating steps: %v", err)
		}

		// Sort the steps array by date-time order
		sort.Slice(sortedSteps, func(i, j int) bool {
			dateTimeI, _ := jsonparser.GetString(sortedSteps[i], "dateTime")
			dateTimeJ, _ := jsonparser.GetString(sortedSteps[j], "dateTime")
			return dateTimeI < dateTimeJ
		})

		// Iterate through the sorted steps
		for _, step := range sortedSteps {
			if index == int(currentStepIndex) {
				latestStep = string(step)
				break
			} else {
				_, dataType, _, err := jsonparser.Get(step, "setWeight")
				if err != nil {
					//return "", fmt.Errorf("error getting setWeight: %v", err)
				} else {
					switch dataType {
					case jsonparser.String:
						setWeight, err = jsonparser.GetString(step, "setWeight")
						if err != nil {
							return "", fmt.Errorf("error getting setWeight: %v", err)
						}
					case jsonparser.Number:
						setWeightFloat, err := jsonparser.GetFloat(step, "setWeight")
						if err != nil {
							return "", fmt.Errorf("error getting setWeight: %v", err)
						}
						setWeight = fmt.Sprintf("%f", setWeightFloat)
					default:
						return "", fmt.Errorf("setWeight is not a valid type: %v", dataType)
					}
				}
			}
			index++
		}

		if latestStep == "" {
			return "", fmt.Errorf("no step found for currentStepIndex: %d", currentStepIndex)
		}

		// Create the result map
		resultWithIndex := map[string]interface{}{
			"step": []string{latestStep},
			//"index":  currentStepIndex,
			"weight": setWeight,
		}

		// Convert the result to JSON
		resultJSON, err := json.Marshal(resultWithIndex)
		if err != nil {
			return "", fmt.Errorf("error marshalling result to JSON: %v", err)
		}

		return string(resultJSON), nil
	}
	return string("{}"), nil

}

