package models

type Bank struct {
	CountryCode  string `bson:"countryCode" json:"countryCode"`
	SwiftCode    string `bson:"swiftCode" json:"swiftCode"`
	CodeType     string `bson:"codeType" json:"codeType"` // This could be a constant, it can be "BIC11" for all of the rows in the CSV file
	BankName     string `bson:"bankName" json:"bankName"`
	Address      string `bson:"address" json:"address"` // This should be optional since some rows in the CSV file don't have an address
	TownName     string `bson:"townName" json:"townName"`
	CountryName  string `bson:"countryName" json:"countryName"`
	TimeZone     string `bson:"timeZone" json:"timeZone"`
	IsHeadOffice bool   `bson:"isHeadOffice" json:"isHeadOffice"`
	BranchCode   string `bson:"branchCode" json:"branchCode"` // This would either specify a branch or the head office (XXX)
}
