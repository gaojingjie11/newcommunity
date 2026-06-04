package model

type ServiceArea struct {
	ID     int64  `gorm:"primaryKey" json:"id"`
	Name   string `gorm:"type:varchar(128)" json:"name"`
	Sort   int    `json:"sort"`
	Status int    `gorm:"not null;default:1" json:"status"`
}

func (ServiceArea) TableName() string { return "service_areas" }
