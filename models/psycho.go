package models

import (
	"encoding/json"
	"reflect"
	"strconv"
)

type Psycho struct {
	UID     string    `json:"uid,omitempty"`
	Psycho  string    `json:"psycho"`
	Values  []*Psycho `json:"values,omitempty"`
	Value   string    `json:"value,omitempty"`
	ID      string    `json:"id,omitempty"`
	AddonID string    `json:"addonId,omitempty"`
	Label   string    `json:"label,omitempty"`
	Ico     string    `json:"ico,omitempty"`
	Pic     string    `json:"pic,omitempty"`
	Sources []*Source `json:"sources,omitempty"`
}

type Source struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
}

func (p *Psycho) UnmarshalJSON(data []byte) error {
	type Alias Psycho
	aux := &struct {
		ID interface{} `json:"id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(aux.ID)
	if v.Kind() == reflect.Float64 {
		f := v.Float()
		i := int(f)
		p.ID = strconv.Itoa(i)
	} else if v.Kind() == reflect.String {
		p.ID = v.String()
	}
	return nil
}

func (s *Source) UnmarshalJSON(data []byte) error {
	type Alias Source
	aux := &struct {
		ID interface{} `json:"id,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(aux.ID)
	if v.Kind() == reflect.Float64 {
		f := v.Float()
		i := int(f)
		s.ID = strconv.Itoa(i)
	} else if v.Kind() == reflect.String {
		s.ID = v.String()
	}
	return nil
}
