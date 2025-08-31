package reaction

import (
	"time"

	"jcourse_go/internal/domain/shared"
)

type Reaction struct {
	ID shared.IDType `json:"id"`

	ReviewID shared.IDType `json:"review_id"`
	UserID   shared.IDType `json:"user_id"`
	Reaction string        `json:"reaction"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
