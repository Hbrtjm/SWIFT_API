package parser

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/Hbrtjm/SWIFT_API/backend/internal/db/models"
)

// SwiftFileParser parses SWIFT code data
type SwiftFileParser struct{}

func NewSwiftFileParser() *SwiftFileParser {
	return &SwiftFileParser{}
}

// ParseFile parses a CSV file containing SWIFT codes and returns array of models.Bank
func (p *SwiftFileParser) ParseFile(filename string) ([]models.Bank, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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

	headers := records[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToUpper(header)] = i
	}

	result := make([]models.Bank, 0, len(records)-1)
	for _, record := range records[1:] {
		if len(record) < len(headers) {
			continue
		}

		swiftCode := getFieldValue(record, headerMap, "SWIFT CODE")
		branchCode := ""
		IsHeadquarter := false
		if len(swiftCode) == 11 {
			branchCode = swiftCode[:8]
			IsHeadquarter = swiftCode[8:] == "XXX"
		} else if len(swiftCode) == 8 {
			IsHeadquarter = true
		}

		bank := models.Bank{
			CountryCode:   strings.ToUpper(getFieldValue(record, headerMap, "COUNTRY ISO2 CODE")),
			SwiftCode:     swiftCode,
			CodeType:      getFieldValue(record, headerMap, "CODE TYPE"),
			BankName:      getFieldValue(record, headerMap, "NAME"),
			Address:       getFieldValue(record, headerMap, "ADDRESS"),
			TownName:      getFieldValue(record, headerMap, "TOWN NAME"),
			CountryName:   strings.ToUpper(getFieldValue(record, headerMap, "COUNTRY NAME")),
			TimeZone:      getFieldValue(record, headerMap, "TIME ZONE"),
			IsHeadquarter: IsHeadquarter,
			BranchCode:    branchCode,
		}

		result = append(result, bank)
	}

	return result, nil
}

func getFieldValue(record []string, headerMap map[string]int, fieldName string) string {
	if index, exists := headerMap[fieldName]; exists && index < len(record) {
		return strings.TrimSpace(record[index])
	}
	return ""
}
