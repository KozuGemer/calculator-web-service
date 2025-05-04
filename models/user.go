// user.go
package models

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"` // На практике пароль хранится в зашифрованном виде
	Token    string `json:"token"`
}

func (u *User) UpdateUserToken(d int, tokenString string) (any, error) {
	panic("unimplemented")
}
