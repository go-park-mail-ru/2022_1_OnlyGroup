package models

type Profile struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Birthday  string   `json:"birthday"`
	City      string   `json:"city"`
	Interests []string `json:"interests"`
	AboutUser string   `json:"aboutUser"`
	UserId    int      `json:"userId"`
	Gender    string   `json:"gender"`
}
