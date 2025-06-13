package converter

import (
	"encoding/json"
	"fp-designpattern/internal/entity"
	"fp-designpattern/internal/model"
)

func CourseToResponse(course *entity.Course) *model.CourseResponse {
	var content []model.ContentBlock
	if err := json.Unmarshal(course.Content, &content); err != nil {
		// fallback to empty slice or handle error as needed
		content = []model.ContentBlock{}
	}
	return &model.CourseResponse{
		ID:         course.ID,
		CourseName: course.CourseName,
		Content:    content,
		GradeLevel: course.GradeLevel,
		CreatedAt:  course.CreatedAt,
		UpdatedAt:  course.UpdatedAt,
		Subject:    *SubjectToResponse(&course.Subject),
	}
}

func CourseToListResponse(course *entity.Course) *model.CourseListResponse {
	return &model.CourseListResponse{
		ID:         course.ID,
		CourseName: course.CourseName,
		GradeLevel: course.GradeLevel,
		CreatedAt:  course.CreatedAt,
		UpdatedAt:  course.UpdatedAt,
		Subject:    *SubjectToResponse(&course.Subject),
	}
}
