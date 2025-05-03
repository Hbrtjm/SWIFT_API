package util

func NewPostRequest() *PostRequest {
	return &PostRequest{
		CountryISO2:   "",
		SwiftCode:     "",
		CodeType:      "",
		BankName:      "",
		Address:       "",
		TownName:      "",
		IsHeadquarter: false,
		CountryName:   "",
		TimeZone:      "",
	}
}

type PostRequest struct {
	CountryISO2   string `json:"countryISO2"`
	SwiftCode     string `json:"swiftCode"`
	CodeType      string `json:"codeType"`
	BankName      string `json:"bankName"`
	Address       string `json:"address"`
	TownName      string `json:"townName"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	TimeZone      string `json:"timeZone"`
}
