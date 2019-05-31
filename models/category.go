package models

type Category struct {
	UID      string      `json:"uid,omitempty"`
	Category string      `json:"category"`
	ID       int         `json:"id,omitempty"`
	Name     string      `json:"name,omitempty"`
	Level    int         `json:"level,omitempty"`
	Children []*Category `json:"children,omitempty"`
	Pic      string      `json:"pic,omitempty"`
	Type     string      `json:"type,omitempty"`
}
