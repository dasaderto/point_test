package models

type NewsResource struct {
	ID               int
	IsRss            bool   `gorm:"not null"`
	Link             string `gorm:"size:255;unique;not null"`
	ItemAttrs        string `gorm:"size:255"`
	PictureAttrs     string `gorm:"size:255"`
	PublishDateAttrs string `gorm:"size:255"`
	TitleAttrs       string `gorm:"size:255"`
	DescriptionAttrs string `gorm:"size:255"`
	LinkAttrs        string `gorm:"size:255"`
}

type News struct {
	ID          int
	ResourceID  int          `gorm:"not null"`
	Resource    NewsResource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Picture     string       `gorm:"size:512"`
	PublishDate string       `gorm:"size:255"`
	Title       string       `gorm:"size:250"`
	Category    string       `gorm:"size:100"`
	Author      string       `gorm:"size:100"`
	Link        string       `gorm:"size:512;unique;not null"`
	Description string
}
