package model

type RevokedToken struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Token string `json:"token"`
}
