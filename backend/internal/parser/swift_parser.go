// backend/internal/parser/swift_parser.go

package parser

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/internal/db/models"
)

// SwiftCodeParser parses SWIFT code data
type SwiftCodeParser struct{}

// NewSwiftCodeParser creates a new SwiftCodeParser
func NewSwiftCodeParser() *SwiftCodeParser {
	return &SwiftCodeParser{}
}

// ParseFile parses a CSV file containing SWIFT codes and returns array of models.Bank
func (p *SwiftCodeParser) ParseFile(filename string) ([]models.Bank, error) {
	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	reader.Comma = ';' // Set the delimiter to semicolon
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// Check for minimum record count
	if len(records) < 2 {
		return []models.Bank{}, nil
	}

	// Map headers for flexible access
	headers := records[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToUpper(header)] = i
	}

	// Process records into Bank models
	result := make([]models.Bank, 0, len(records)-1)
	for _, record := range records[1:] {
		if len(record) < len(headers) {
			continue
		}

		swiftCode := getFieldValue(record, headerMap, "SWIFT CODE")
		branchCode := ""
		isHeadOffice := false
		if len(swiftCode) == 11 {
			branchCode = swiftCode[8:]
			isHeadOffice = branchCode == "XXX"
		} else if len(swiftCode) == 8 {
			isHeadOffice = true
		}

		bank := models.Bank{
			CountryCode:  strings.ToUpper(getFieldValue(record, headerMap, "COUNTRY ISO2 CODE")),
			SwiftCode:    swiftCode,
			CodeType:     getFieldValue(record, headerMap, "CODE TYPE"),
			BankName:     getFieldValue(record, headerMap, "NAME"),
			Address:      getFieldValue(record, headerMap, "ADDRESS"),
			TownName:     getFieldValue(record, headerMap, "TOWN NAME"),
			CountryName:  strings.ToUpper(getFieldValue(record, headerMap, "COUNTRY NAME")),
			TimeZone:     getFieldValue(record, headerMap, "TIME ZONE"),
			IsHeadOffice: isHeadOffice,
			BranchCode:   branchCode,
		}

		result = append(result, bank)
	}

	return result, nil
}

// getFieldValue safely gets a field value from a record using the header map
func getFieldValue(record []string, headerMap map[string]int, fieldName string) string {
	if index, exists := headerMap[fieldName]; exists && index < len(record) {
		return record[index]
	}
	return ""
}
