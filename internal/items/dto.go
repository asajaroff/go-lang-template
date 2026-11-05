package items

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	itemsdb "github.com/asajaroff/go-lang-template/internal/items/db"
)

// createRequest / updateRequest / response are the wire (JSON) shapes —
// kept distinct from itemsdb.Item (the DB row shape) the same way
// SQLAlchemy models and Pydantic schemas are kept distinct in the Python
// sibling template: the API contract and the storage schema are allowed to
// evolve independently.
type createRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func (r createRequest) validate() error {
	if l := len(r.Name); l < 1 || l > 255 {
		return fmt.Errorf("name must be 1-255 characters")
	}
	if r.Description != nil && len(*r.Description) > 1024 {
		return fmt.Errorf("description must be at most 1024 characters")
	}
	return nil
}

// updateRequest fields are pointers so a PATCH body can omit a field
// entirely (leave it unchanged) — distinct from sending an empty string.
type updateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (r updateRequest) validate() error {
	if r.Name != nil {
		if l := len(*r.Name); l < 1 || l > 255 {
			return fmt.Errorf("name must be 1-255 characters")
		}
	}
	if r.Description != nil && len(*r.Description) > 1024 {
		return fmt.Errorf("description must be at most 1024 characters")
	}
	return nil
}

type response struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toResponse(item itemsdb.Item) response {
	var description *string
	if item.Description.Valid {
		description = &item.Description.String
	}
	return response{
		ID:          item.ID,
		Name:        item.Name,
		Description: description,
		CreatedAt:   item.CreatedAt.Time,
		UpdatedAt:   item.UpdatedAt.Time,
	}
}

func toPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *s, Valid: true}
}
