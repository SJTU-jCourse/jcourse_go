package subscribe

import "jcourse_go/internal/domain/shared"

type CourseSubscribeCommand struct {
	CourseID shared.IDType
	Level    SubscribeLevel
}
