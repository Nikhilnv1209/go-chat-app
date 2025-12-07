package models

type Group struct {
	BaseModel
	Name string `gorm:"size:100;not null" json:"name"`
}
