package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`   // never expose
	Role     string `json:"role"` // "admin" | "user"
}
