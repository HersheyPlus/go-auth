package models

type User struct {
    Base
    Username  string  `gorm:"type:varchar(100);not null;unique;index" json:"username"`
    FirstName *string `gorm:"type:varchar(100)" json:"first_name,omitempty"`
    LastName  *string `gorm:"type:varchar(100)" json:"last_name,omitempty"`
    Phone     string  `gorm:"type:varchar(20);not null" json:"phone"`
    Email     string  `gorm:"type:varchar(100);unique;not null;index" json:"email"`
    Password  string  `gorm:"type:varchar(255);not null" json:"-"`
}
