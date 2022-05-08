package validator

import "BroadcastService/common"

// Validator validates input
type Validator interface {
	// validates the genuineness of the data
	Validate(data *common.Data) bool
}
