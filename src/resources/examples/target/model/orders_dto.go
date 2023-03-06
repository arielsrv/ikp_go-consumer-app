package model

type OrderDTO struct {
	ID      int64  `json:"id,omitempty"`
	Details string `json:"details,omitempty"`
}
