package common

import (
	"encoding/json"
)

type GetProjectInput struct {
	ID        string
	Cacheable bool
}

type GetProjectOutput struct {
	ID          string                 `json:"id"`
	DisplayName string                 `json:"displayName"`
	Frozen      bool                   `json:"frozen"`
	Extra       map[string]interface{} `json:"extra"`
}

func (p *GetProjectOutput) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	b, _ := json.Marshal(p)
	json.Unmarshal(b, &m)
	return m
}
