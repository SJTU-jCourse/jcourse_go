package olddto

import (
	"jcourse_go/internal/domain/user"
)

type CommonInfoResponse struct {
	User     user.UserDetail        `json:"user"`
	Settings map[string]interface{} `json:"settings"`
}
