package api

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// ToJSON converts a slice of strings to datatypes.JSON
func ToJSON(slice []string) datatypes.JSON {
	if len(slice) == 0 {
		return datatypes.JSON([]byte("[]"))
	}
	b, _ := json.Marshal(slice)
	return datatypes.JSON(b)
}
