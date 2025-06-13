package converter

import (
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"
)

func SubjectToResponse(subject *entity.Subject) *model.SubjectResponse {
	return &model.SubjectResponse{
		ID:          &subject.ID,
		SubjectName: subject.SubjectName,
		CreatedAt:   &subject.CreatedAt,
		UpdatedAt:   &subject.UpdatedAt,
	}
}
