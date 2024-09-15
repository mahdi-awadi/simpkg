package types

import (
	"encoding/xml"
	"strconv"
)

type RoundedFloat float64

// MarshalJSON rounded float
func (r RoundedFloat) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(r), 'f', 2, 64)), nil
}

// MarshalXMLAttr rounded float
func (r RoundedFloat) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	s := strconv.FormatFloat(float64(r), 'f', 2, 64)
	return xml.Attr{Name: name, Value: s}, nil
}
