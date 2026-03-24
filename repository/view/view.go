package view

import (
	"context"
	"fmt"
	"imchinese/repository/generated"
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

func (r *Repository) Save(ctx context.Context, e models.Existence) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return saveWithoutTransaction(ctx, tx, e)
	})
}

func findViewByName(ctx context.Context, db *gorm.DB, name string) (*models.View, error) {
	views, err := gorm.G[models.View](db).Where(generated.View.Name.Eq(name)).Find(ctx)
	if err != nil {
		return nil, err
	}
	if len(views) > 1 {
		panic(fmt.Errorf("unimplemented View name conflict solution on %s", name))
	}
	var view *models.View
	if len(views) == 1 {
		view = &views[0]
	}
	return view, nil
}

// undefinedModelID provides a Model ID for View s that have not been concentrated.
// The NULL FK in SQL may be another solution, but GORM CLI use zero when we left FK field as empty,
// which eventually leads to the zero FK as a result or breaking the FOREIGN KEY constraint.
func undefinedModelID() int {
	return 0
}

func saveWithoutTransaction(ctx context.Context, tx *gorm.DB, e models.Existence) error {
	view, err := findViewByName(ctx, tx, e.View.Name)
	if err != nil {
		return err
	}

	// I have tried for hours, but failed to leverage the gorm CLI Association to create records in multiple tables.
	// Maybe either the Go and SQL schema shall be aligned with the gorm preferred way, whose document lacks.
	// One more thing, I shall highlight that our goal is to keep simple stable code, not introducing another DSL.

	if view == nil {
		view = &models.View{
			Name:    e.View.Name,
			ModelID: undefinedModelID(),
		}
		if err := gorm.G[models.View](tx).Create(ctx, view); err != nil {
			return err
		}
	}

	// Should have been created lines ago if nil.
	//goland:noinspection GoMaybeNil
	viewID := view.ID
	return gorm.G[models.Existence](tx).Set(
		generated.Existence.ViewID.Set(viewID),
		generated.Existence.Tag.Set(e.Tag),
		generated.Existence.Reason.Set(e.Reason),
		generated.Existence.Source.Set(e.Source),
		generated.Existence.Quote.Set(e.Quote),
	).Create(ctx)
}
