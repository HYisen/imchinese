package view

import (
	"context"
	"imchinese/repository/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*Repository, error) {
	return &Repository{
		db,
	}, nil
}

func (r *Repository) FindAll(ctx context.Context) ([]models.View, error) {
	return gorm.G[models.View](r.db).
		Joins(clause.Has("Model"), nil).
		Find(ctx)
}
