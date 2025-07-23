package models

type Client struct {
	UserId       int     `json:"user_id"`
	Name         string  `json:"name"`
	Passport     string  `json:"passport"`
	USD          float64 `json:"usd"`
	EUR          float64 `json:"eur"`
	Currency     string  `json:"currency"`
	ReceivedName string  `json:"received_name"`
	CreatedDate  string  `json:"created_date"`
}
type Login struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type SearchClient struct {
	Passport string `json:"passport"`
}
