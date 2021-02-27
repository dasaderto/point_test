package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"news/db/models"
)

type INewsRepository interface {
	GetAllWithPag(page, limit int, query string) ([]models.News, error)
	Create(news interface{}) error
	GetByLinks(links []string) ([]models.News, error)
}

type NewsRepository struct {
	Db *gorm.DB
}

func NewNewsRepository(conn *gorm.DB) INewsRepository {
	var repos INewsRepository
	repos = NewsRepository{
		Db: conn,
	}
	return repos
}

func (r NewsRepository) GetAllWithPag(page, limit int, searchQuery string) ([]models.News, error){
	var news []models.News
	query := r.Db.Order("id desc").Offset((page * limit) - limit).Limit(limit)
	if len(searchQuery) > 0 {
		query = query.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", searchQuery))
	}
	err := query.Find(&news).Error

	return news, err
}

func (r NewsRepository) Create(news interface{}) error {
	return r.Db.Create(news).Error
}

func (r NewsRepository) GetByLinks(links []string) ([]models.News, error) {
	var news []models.News
	err := r.Db.Model(models.News{}).Where("link IN ?", links).Find(&news).Error

	return news, err
}