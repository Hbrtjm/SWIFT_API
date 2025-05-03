// First, let's fix the ParseFile method in parser/parser.go to return both banks and countries

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

// ParseFile parses a CSV file containing SWIFT codes and returns array of models.Bank and an array of countries, that is models.Country
func (p *SwiftFileParser) ParseFile(filename string) ([]models.Bank, []models.Country, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // Set the delimiter to semicolon
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	// Check for minimum record count
	if len(records) < 2 {
		return []models.Bank{}, []models.Country{}, nil
	}

	headers := records[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.ToUpper(header)] = i
	}

	// Create maps to keep track of unique countries
	countryMap := make(map[string]models.Country)

	bankResults := make([]models.Bank, 0, len(records)-1)
	for _, record := range records[1:] {
		if len(record) < len(headers) {
			continue
		}

		swiftCode := getFieldValue(record, headerMap, "SWIFT CODE")
		branchCode := ""
		isHeadquarter := false
		if len(swiftCode) == 11 {
			branchCode = swiftCode[:8]
			isHeadquarter = swiftCode[8:] == "XXX"
		} else if len(swiftCode) == 8 {
			isHeadquarter = true
		}

		countryISO2 := strings.ToUpper(getFieldValue(record, headerMap, "COUNTRY ISO2 CODE"))
		countryName := strings.ToUpper(getFieldValue(record, headerMap, "COUNTRY NAME"))
		timeZone := getFieldValue(record, headerMap, "TIME ZONE")

		// Add country to the map if it doesn't exist
		if _, exists := countryMap[countryISO2]; !exists && countryISO2 != "" {
			countryMap[countryISO2] = models.Country{
				CountryISO2: countryISO2,
				CountryName: countryName,
				TimeZone:    timeZone,
			}
		}

		bank := models.Bank{
			CountryISO2: countryISO2,
			SwiftCode:   swiftCode,
			CodeType:    getFieldValue(record, headerMap, "CODE TYPE"),
			BankName:    getFieldValue(record, headerMap, "NAME"),
			Address:     getFieldValue(record, headerMap, "ADDRESS"),
			TownName:    getFieldValue(record, headerMap, "TOWN NAME"),
			IsHeadquarter: isHeadquarter,
			BranchCode:    branchCode,
		}

		bankResults = append(bankResults, bank)
	}

	// Convert countries map to slice
	countryResults := make([]models.Country, 0, len(countryMap))
	for _, country := range countryMap {
		countryResults = append(countryResults, country)
	}

	return bankResults, countryResults, nil
}

func getFieldValue(record []string, headerMap map[string]int, fieldName string) string {
	if index, exists := headerMap[fieldName]; exists && index < len(record) {
		return strings.TrimSpace(record[index])
	}
	return ""
}
