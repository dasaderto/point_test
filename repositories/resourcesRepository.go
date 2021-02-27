package repositories

import (
	"gorm.io/gorm"
	"news/db/models"
)

type IResourcesRepository interface {
	GetAll() ([]models.NewsResource, error)
	Create(resource *models.NewsResource) error
	Destroy(pk int) error
}

type ResourcesRepository struct {
	Db *gorm.DB
}

func NewResourcesRepository(conn *gorm.DB) IResourcesRepository {
	var repos IResourcesRepository
	repos = ResourcesRepository{
		Db: conn,
	}
	return repos
}

func (r ResourcesRepository) GetAll() ([]models.NewsResource, error) {
	var resources []models.NewsResource
	err := r.Db.Model(models.NewsResource{}).Find(&resources).Error

	return resources, err
}

func (r ResourcesRepository) Create(resource *models.NewsResource) error {
	return r.Db.FirstOrCreate(resource, models.NewsResource{Link: resource.Link}).Error
}

func (r ResourcesRepository) Destroy(pk int) error {
	return r.Db.Delete(&models.NewsResource{ID: pk}).Error
}