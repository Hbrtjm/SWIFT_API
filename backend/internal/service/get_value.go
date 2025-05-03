package service

// Helper function to safely extract string values from map
func getValue(data map[string]interface{}, key string) string {
	if value, exists := data[key]; exists {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return ""
}
