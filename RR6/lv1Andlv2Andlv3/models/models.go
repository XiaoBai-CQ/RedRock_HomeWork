package models

type Student struct {
	Id    int    `json:"Id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"Name" gorm:"type:varchar(100);not null"`
	Sex   string `json:"Sex" gorm:"type:varchar(10);default:'unknown'"`
	Born  string `json:"Born" gorm:"type:date"`
	Birth string `json:"Birth" gorm:"type:varchar(255)"`
}

type User struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	Username         string `json:"username" gorm:"unique;not null"`
	Password         string `json:"password" gorm:"not null"`
	SecurityQuestion string `json:"security_question" gorm:"not null"`
	SecurityAnswer   string `json:"security_answer" gorm:"not null"`
}
