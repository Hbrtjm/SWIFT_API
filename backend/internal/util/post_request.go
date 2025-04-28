package util

func NewPostRequest() *PostRequest {
	return &PostRequest{
		CountryISO2: "",
		SwiftCode:   "",
		CodeType:    "",
		BankName:    "",
		Address:     "",
		TownName:    "",
		CountryName: "",
		TimeZone:    "",
	}
}

type PostRequest struct {
	CountryISO2   string `json:"countryISO2"`
	SwiftCode     string `json:"swiftCode"`
	CodeType    string `json:"codeType"`
	BankName      string `json:"bankName"`
	Address       string `json:"address"`
	TownName	  string `json:"townName"`
	CountryName   string `json:"countryName"`
	TimeZone 	 string `json:"timezone"`
}
