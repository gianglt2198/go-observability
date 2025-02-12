package utils

import "encoding/json"

func TransformToType[T any](payload []byte) (*T, error) {

	var res T
	err := json.Unmarshal(payload, &res)

	return &res, err
}
