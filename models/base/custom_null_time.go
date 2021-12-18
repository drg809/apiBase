package base

import (
	"database/sql"
	"encoding/json"
	"errors"
)

type CustomNullTime struct {
	sql.NullTime
}

func (customNullTime CustomNullTime) MarshalJSON() ([]byte, error) {
	if customNullTime.Valid {
		return json.Marshal(customNullTime.Time)
	}
	return json.Marshal(nil)
}

func (customNullTime *CustomNullTime) UnmarshalJSON(bytes []byte) error {
	if string(bytes) == "null" {
		customNullTime.Valid = false
		return errors.New("invalid date")
	}

	err := json.Unmarshal(bytes, &customNullTime.Time)
	if err != nil {
		return err
	}

	customNullTime.Valid = true
	return nil
}
