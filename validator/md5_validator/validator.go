package md5_validator

import (
	"BroadcastService/common"
	"BroadcastService/validator"
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type MD5Validator struct {
}

func (M MD5Validator) Validate(data *common.Data) bool {
	if data.Id == "" || data.Raw == "" || data.Signature == "" {
		return false
	}

	m := md5.New()
	m.Write([]byte(data.Raw))
	digest := m.Sum([]byte(nil))

	return strings.ToLower(hex.EncodeToString(digest)) == strings.ToLower(data.Signature)
}

func New() validator.Validator {
	return &MD5Validator{}
}
