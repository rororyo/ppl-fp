package converter

import (
	"blessing-be/internal/entity"
	"blessing-be/internal/model"
)

func UserCourseToResponse(userCourse *entity.UserCourse) *model.UserCourseResponse {
	return &model.UserCourseResponse{
		ID:         userCourse.ID,
		User:       *UserToResponse(&userCourse.User),
		Course:     *CourseToResponse(&userCourse.Course),
		AccessedAt: userCourse.AccessedAt,
	}
}

func UserCourseListToResponse(userCourse *entity.UserCourse) *model.UserCourseListResponse {
	return &model.UserCourseListResponse{
		ID:         userCourse.ID,
		User:       *UserToResponse(&userCourse.User),
		Course:     *CourseToListResponse(&userCourse.Course),
		AccessedAt: userCourse.AccessedAt,
	}
}
