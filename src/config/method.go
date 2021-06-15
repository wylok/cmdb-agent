package config

import "encoding/json"

func (co CollectionData) CollectionDataString() string {
	s, _ := json.Marshal(co)
	return string(s)
}
