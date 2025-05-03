package validators

import (
	"testing"
)

func TestValidInput(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]interface{}
		expectedError bool
	}{
		{
			name: "Valid input",
			input: map[string]interface{}{
				"field1": "validValue",
				"field2": "anotherValidValue",
			},
			expectedError: false,
		},
		{
			name: "Input with non-string field",
			input: map[string]interface{}{
				"field1": "validValue",
				"field2": 12345,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Sanitize(tt.input)
			if (err != nil) != tt.expectedError {
				t.Errorf("Sanitize() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func SingleCaracter(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]interface{}
		expectedError bool
	}{
		{
			name: "Input with illegal character $",
			input: map[string]interface{}{
				"field1": "validValue",
				"field2": "illegal$value",
			},
			expectedError: true,
		},
		{
			name: "Input with illegal character {",
			input: map[string]interface{}{
				"field1": "validValue",
				"field2": "illegal{value",
			},
			expectedError: true,
		},
		{
			name: "Input with illegal character }",
			input: map[string]interface{}{
				"field1": "validValue",
				"field2": "illegal}value",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Sanitize(tt.input)
			if (err != nil) != tt.expectedError {
				t.Errorf("Sanitize() error = %v, expectedError %v", err, tt.expectedError)
			}
		})
	}
}

func TestForInjection(t *testing.T) {
	input := map[string]interface{}{
		"field1": "${RCE}",
		"field2": "${id}",
	}
	err := Sanitize(input)
	if err == nil {
		t.Error("Sanitize() expected error but got none")
	}
}
