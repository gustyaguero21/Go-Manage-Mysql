package models

type User struct {
	ID       string `gorm:"primaryKey;type:varchar(36);not null;unique" json:"id"`
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	Surname  string `gorm:"type:varchar(255);not null" json:"surname"`
	Username string `gorm:"type:varchar(255);not null;unique" json:"username"`
	Phone    string `gorm:"type:varchar(255);not null;unique" json:"phone"`
	Email    string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`
}

type UserResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}
