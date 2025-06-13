package repository

import (
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	Repository[entity.Subject]
	Log *logrus.Logger
}

func NewSubjectRepository(log *logrus.Logger) *SubjectRepository {
	return &SubjectRepository{
		Log: log,
	}
}

func (r *SubjectRepository) Search(db *gorm.DB, request *model.SearchSubjectRequest) ([]entity.Subject, int64, error) {
	var subjects []entity.Subject
	if err := db.Scopes(r.FilterSubject(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&subjects).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Subject{}).Scopes(r.FilterSubject(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return subjects, total, nil
}

func (r *SubjectRepository) FilterSubject(request *model.SearchSubjectRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if subjectName := request.SubjectName; subjectName != "" {
			tx = tx.Where("username LIKE ?", "%"+subjectName+"%")
		}
		return tx
	}
}
