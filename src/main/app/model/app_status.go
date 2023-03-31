package model

import "encoding/json"

// AppStatusDTO  Model
// swagger:model AppStatusDTO
type AppStatusDTO struct {
	Status Status `json:"status,omitempty"`
}

func (a AppStatusDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

type Status string

const (
	Started Status = "started"
	Stopped Status = "stopped"
)
