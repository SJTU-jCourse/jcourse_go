package main

import "time"

type DepartmentV1 struct {
	ID   int64
	Name string
}

func (p *DepartmentV1) TableName() string {
	return "jcourse_api_department"
}

type SemesterV1 struct {
	ID   int64
	Name string
}

func (p *SemesterV1) TableName() string {
	return "jcourse_api_semester"
}

type CategoryV1 struct {
	ID   int64
	Name string
}

func (p *CategoryV1) TableName() string {
	return "jcourse_api_categories"
}

type CourseV1 struct {
	ID            int64
	Code          string
	Name          string
	Credit        float64
	DepartmentID  int64
	MainTeacherID int64
}

func (p *CourseV1) TableName() string {
	return "jcourse_api_course"
}

type TeacherV1 struct {
	ID           int64
	Tid          string
	Name         string
	Title        string
	Pinyin       string
	AbbrPinyin   string
	DepartmentID int64
}

func (p *TeacherV1) TableName() string {
	return "jcourse_api_teacher"
}

type ReviewV1 struct {
	ID         int64
	CourseID   int64
	UserID     int64
	SemesterID int64
	Comment    string
	Score      string
	Rating     int64
	CreatedAt  time.Time
	ModifiedAt time.Time
}

func (p *ReviewV1) TableName() string {
	return "jcourse_api_review"
}

type ReviewRevisionV1 struct {
	ID         int64
	CourseID   int64
	UserID     int64
	ReviewID   int64
	SemesterID int64
	Comment    string
	Score      string
	Rating     int64
	CreatedAt  time.Time
}

func (p *ReviewRevisionV1) TableName() string {
	return "jcourse_api_reviewrevision"
}

type UserV1 struct {
	ID         int64
	Password   string
	Username   string
	DateJoined time.Time
	LastLogin  time.Time
}

func (p *UserV1) TableName() string {
	return "auth_user"
}

type ReviewReactionV1 struct {
	ID         int64
	ReviewID   int64
	UserID     int64
	Reaction   int64
	ModifiedAt time.Time
}

func (p *ReviewReactionV1) TableName() string {
	return "jcourse_api_reviewreaction"
}

type UserPointV1 struct {
	ID          int64
	UserID      int64
	Value       int64
	Description string
	Time        time.Time
}

func (p *UserPointV1) TableName() string {
	return "jcourse_api_userpoint"
}

type UserProfileV1 struct {
	ID            int64
	UserType      string
	UserID        int64
	LowerCase     bool
	SuspendedTill time.Time
	LastSeenAt    time.Time
}

func (p *UserProfileV1) TableName() string {
	return "oauth_userprofile"
}
