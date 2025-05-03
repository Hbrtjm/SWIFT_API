package models

type Country struct {
	CountryISO2 string `bson:"countryISO2" json:"countryISO2"`
	CountryName string `bson:"countryName" json:"countryName"`
	TimeZone    string `bson:"timeZone" json:"timeZone"`
}
