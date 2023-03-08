package model

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON implements custom unmarshaller that
// always treats the trait value as a string, keeping
// the type consistent with the schema
// display_type can be used for clients to parse into
// numeric types if required
func (attr *ERC721Attribute) UnmarshalJSON(b []byte) error {
	var m struct {
		DisplayType string      `json:"display_type"`
		TraitType   string      `json:"trait_type"`
		Value       interface{} `json:"value"`
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	attr.DisplayType = m.DisplayType
	attr.TraitType = m.TraitType
	switch v := m.Value.(type) {
	case string:
		attr.Value = v
	default:
		attr.Value = fmt.Sprintf("%v", v)
	}
	return nil
}
