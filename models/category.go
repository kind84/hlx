package models

type Category struct {
	Category string      `json:"category"`
	Uid      string      `json:"uid,omitempty"`
	ID       int         `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	Level    int         `json:"level,omitempty"`
	Children []*Category `json:"children,omitempty"`
}
