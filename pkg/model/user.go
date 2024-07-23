package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
