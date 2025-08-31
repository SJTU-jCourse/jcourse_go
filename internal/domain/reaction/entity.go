package reaction

import "time"

type Reaction struct {
	ID       int    `json:"id"`
	ReviewID int    `json:"review_id"`
	Reaction string `json:"reaction"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
