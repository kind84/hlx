package models

type Psycho struct {
	UID     string      `json:"uid,omitempty"`
	Psycho  string      `json:"psycho"`
	Values  []*Psycho   `json:"values,omitempty"`
	Value   string      `json:"value,omitempty"`
	ID      interface{} `json:"id,omitempty"`
	AddonID string      `json:"addonId,omitempty"`
	Label   string      `json:"label,omitempty"`
	Ico     string      `json:"ico,omitempty"`
	Pic     string      `json:"pic,omitempty"`
	Sources []*Source   `json:"sources,omitempty"`
}

type Source struct {
	ID          *int   `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}
