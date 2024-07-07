package dto

type TeacherDTO struct {
	ID         int64  `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Department string `json:"department"`
	Title      string `json:"title"`
}

type OfferedCourseDTO struct {
	ID           int64        `json:"id"`
	Semester     string       `json:"semester"`
	Grade        []string     `json:"grade"`
	Language     string       `json:"language"`
	TeacherGroup []TeacherDTO `json:"teacher_group"`
}

type CourseDetailDTO struct {
	ID            int64              `json:"id"`
	Code          string             `json:"code"`
	Name          string             `json:"name"`
	Credit        float64            `json:"credit"`
	MainTeacher   TeacherDTO         `json:"main_teacher"`
	OfferedCourse []OfferedCourseDTO `json:"offered_courses"`
}

type CourseDetailRequest struct {
	CourseID int64 `uri:"courseID" binding:"required"`
}

type CourseListItemDTO struct {
	ID          int64      `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Credit      float64    `json:"credit"`
	MainTeacher TeacherDTO `json:"main_teacher"`
	Categories  []string   `json:"categories"`
	Department  string     `json:"department"`
}

type CourseListRequest struct {
	Page        int64  `json:"page" form:"page"`
	PageSize    int64  `json:"page_size" form:"page_size"`
	Departments string `json:"departments" form:"departments"`
	Categories  string `json:"categories" form:"categories"`
	Credits     string `json:"credits" form:"credits"`
}

type CourseListResponse = BasePaginateResponse[CourseListItemDTO]
