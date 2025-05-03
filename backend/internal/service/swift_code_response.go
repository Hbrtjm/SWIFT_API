package service

type SwiftCodeResponse struct {
	Address       string                   `json:"address"`
	BankName      string                   `json:"bankName"`
	CountryISO2   string                   `json:"countryISO2"`
	CountryName   string                   `json:"countryName"`
	IsHeadquarter bool                     `json:"isHeadquarter"`
	SwiftCode     string                   `json:"swiftCode"`
	Branches      []map[string]interface{} `json:"branches,omitempty"`
}
