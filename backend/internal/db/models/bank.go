package models

// Renamed from models, since we have a models package already
type Bank struct {
	CountryISO2 string `bson:"countryISO2" json:"countryISO2"`
	SwiftCode   string `bson:"swiftCode" json:"swiftCode"`
	CodeType    string `bson:"codeType" json:"codeType"` // This could be a constant, it can be "BIC11" for all of the rows in the CSV file
	BankName    string `bson:"bankName" json:"bankName"`
	Address     string `bson:"address" json:"address"` // This should be optional since some rows in the CSV file don't have an address
	TownName    string `bson:"townName" json:"townName"`
	// CountryName   string `bson:"countryName" json:"countryName"` // Deprecated, in the final patch I will remove this field
	// TimeZone      string `bson:"timeZone" json:"timeZone"` // Deprecated, in the final patch I will remove this field
	IsHeadquarter bool   `bson:"isHeadquarter" json:"isHeadquarter"`
	BranchCode    string `bson:"branchCode" json:"branchCode"` // This would either specify a branch or the head office (XXX)
}
