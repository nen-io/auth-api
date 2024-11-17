package models

type User struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (u *User) Valid() bool {
	if u.FirstName == "" || u.LastName == "" || u.Email == "" || u.Password == "" {
		return false
	}
	return true
}
