package presenter

import "time"

type Category struct {
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	ParentID         string     `json:"parentId,omitempty"`
	Path             string     `json:"path"`
	Level            int        `json:"level"`
	ParentCategories []Category `json:"parents,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt,omitempty"`
}
