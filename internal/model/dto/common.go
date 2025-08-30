package dto

import (
	"jcourse_go/internal/model/model"
)

type CommonInfoResponse struct {
	User     model.UserDetail       `json:"user"`
	Settings map[string]interface{} `json:"settings"`
}
