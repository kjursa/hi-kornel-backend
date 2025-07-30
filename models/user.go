package models

type User struct {
	Uid  string
	Name string
}

type Token struct {
	Access  string
	Refresh string
}

type UserToken struct {
	User  User
	Token Token
}

type Skill struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

type Profile struct {
	Description string  `json:"description"`
	Skills      []Skill `json:"skills"`
}

type Message struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Email   string `json:"email"`
}

type Chat struct {
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	Timestamp int64  `json:"timestamp"`
	MessageId string `json:"messageId"`
	ChatId    string `json:"chatId"`
}

type Experience struct {
	Name      string    `json:"name"`
	Title     string    `json:"title"`
	Companies []Company `json:"companies"`
}

type Company struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Period      string    `json:"period"`
	Projects    []Project `json:"projects"`
}

type Project struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Points      []string `json:"points"`
	Technology  []string `json:"technology"`
}
