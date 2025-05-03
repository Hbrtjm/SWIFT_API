package service

import (
	"strconv"
)

// Helper function to safely extract boolean values from map
func getBool(data map[string]interface{}, key string) bool {
	if value, exists := data[key]; exists {
		if boolValue, ok := value.(bool); ok {
			return boolValue
		}

		if strValue, ok := value.(string); ok {
			if boolValue, err := strconv.ParseBool(strValue); err == nil {
				return boolValue
			}
		}

		if numValue, ok := value.(float64); ok {
			return numValue != 0
		}
	}
	return false
}
