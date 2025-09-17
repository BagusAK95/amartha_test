package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

const InsertBatchSize = 100

func Value[M any](model M) (driver.Value, error) {
	bytes, err := json.Marshal(model)

	return string(bytes), err
}

func Scan[M any](model M, value any) error {
	var bytes []byte

	switch x := value.(type) {
	case []byte:
		bytes = x
	case string:
		bytes = []byte(x)
	default:
		return errors.New("invalid type for JSONB")
	}

	return json.Unmarshal(bytes, &model)
}

type JSONB map[string]any

func (j JSONB) Value() (driver.Value, error) {
	return Value(j)
}

func (j *JSONB) Scan(value any) error {
	return Scan(j, value)
}
