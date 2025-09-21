package main

import "time"

type Department struct {
	ID   int64
	Name string
}

func (p *Department) TableName() string {
	return "jcourse_api_department"
}

type Semester struct {
	ID        int64
	Name      string
	Available bool
}

func (p *Semester) TableName() string {
	return "jcourse_api_semester"
}

type Category struct {
	ID   int64
	Name string
}

func (p *Category) TableName() string {
	return "jcourse_api_category"
}

type Course struct {
	ID             int64
	Code           string
	Name           string
	Credit         float64
	DepartmentID   int64
	Department     *Department
	MainTeacherID  int64
	MainTeacher    *Teacher
	Categories     []Category `gorm:"many2many:jcourse_api_course_categories;"`
	TeacherGroup   []Teacher  `gorm:"many2many:jcourse_api_course_teacher_group;"`
	LastSemesterID int64
	LastSemester   *Semester
}

func (p *Course) TableName() string {
	return "jcourse_api_course"
}

type Teacher struct {
	ID           int64
	Tid          string
	Name         string
	Title        string
	Pinyin       string
	AbbrPinyin   string
	DepartmentID int64
	Department   *Department
}

func (p *Teacher) TableName() string {
	return "jcourse_api_teacher"
}

type Review struct {
	ID         int64
	CourseID   int64
	Course     *Course
	UserID     int64
	User       *User
	SemesterID int64
	Semester   *Semester
	Comment    string
	Score      string
	Rating     int64
	CreatedAt  time.Time
	ModifiedAt time.Time

	Revisions []ReviewRevision
	Reactions []ReviewReaction
}

func (p *Review) TableName() string {
	return "jcourse_api_review"
}

type ReviewRevision struct {
	ID         int64
	CourseID   int64
	Course     *Course
	UserID     int64
	User       *User
	ReviewID   int64
	Review     *Review
	SemesterID int64
	Semester   *Semester
	Comment    string
	Score      string
	Rating     int64
	CreatedAt  time.Time
}

func (p *ReviewRevision) TableName() string {
	return "jcourse_api_reviewrevision"
}

type User struct {
	ID          int64
	Password    string
	Username    string
	UserProfile *UserProfile
	DateJoined  time.Time
	LastLogin   time.Time

	Points []UserPoint
}

func (p *User) TableName() string {
	return "auth_user"
}

type ReviewReaction struct {
	ID         int64
	ReviewID   int64
	Review     *Review
	UserID     int64
	User       *User
	Reaction   int64
	ModifiedAt time.Time
}

func (p *ReviewReaction) TableName() string {
	return "jcourse_api_reviewreaction"
}

type UserPoint struct {
	ID          int64
	UserID      int64
	User        *User
	Value       int64
	Description string
	Time        time.Time
}

func (p *UserPoint) TableName() string {
	return "jcourse_api_userpoint"
}

type UserProfile struct {
	ID            int64
	UserType      string
	UserID        int64
	User          *User
	LowerCase     bool `gorm:"column:lowercase"`
	SuspendedTill *time.Time
	LastSeenAt    time.Time
}

func (p *UserProfile) TableName() string {
	return "oauth_userprofile"
}

type CourseNotificationLevel struct {
	ID                int64
	NotificationLevel int64 // 0 正常，1 关注，2 忽略
	CourseID          int64
	UserID            int64
	ModifiedAt        time.Time
}

func (c *CourseNotificationLevel) TableName() string {
	return "jcourse_api_coursenotificationlevel"
}

type EnrollCourse struct {
	ID         int64
	CourseID   int64
	UserID     int64
	SemesterID int64
	Semester   *Semester
	CreatedAt  time.Time
}

func (p *EnrollCourse) TableName() string {
	return "jcourse_api_enrollcourse"
}
