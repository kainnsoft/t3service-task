package entity

type User struct {
	Id    int    `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}

func NewUser(email string) *User {
	return &User{Email: email}
}
