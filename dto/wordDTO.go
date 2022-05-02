package dto

type WordDTO struct {
	English      string `bson:"english"`
	Turkish      string `bson:"turkish"`
	Abbreviation string `bson:"abbreviation"`
	Description  string `bson:"description"`
}
