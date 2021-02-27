package repositories

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"news/db/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(conn *gorm.DB) Repository {
	repos := Repository{
		db: conn,
	}
	return repos
}

func (r Repository) Migrate() {
	err := r.db.AutoMigrate(&models.News{}, &models.NewsResource{})
	if err != nil {
		log.Error("Fail migration" + err.Error())
	}
}