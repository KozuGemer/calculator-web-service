// user.go
package models

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"` // На практике пароль хранится в зашифрованном виде
}
