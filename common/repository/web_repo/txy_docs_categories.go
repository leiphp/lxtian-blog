package web_repo

import (
	"lxtian-blog/common/model"
	"lxtian-blog/common/repository"

	"gorm.io/gorm"
)

type TxyDocsCategoriesRepository interface {
	repository.BaseRepository[model.TxyDocsCategory]
}

type txyDocsCategoriesRepository struct {
	*repository.TransactionalBaseRepository[model.TxyDocsCategory]
}

func NewTxyDocsCategoriesRepository(db *gorm.DB) TxyDocsCategoriesRepository {
	return &txyDocsCategoriesRepository{
		TransactionalBaseRepository: repository.NewTransactionalBaseRepository[model.TxyDocsCategory](db),
	}
}
