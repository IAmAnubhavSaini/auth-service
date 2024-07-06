package types

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var ServiceUsers = map[string]string{}
