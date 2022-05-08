package md5_validator

import (
	"BroadcastService/common"
	"crypto/md5"
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestMD5Validator_Validate(t *testing.T) {
	v := New()
	data := &common.Data{
		Id:  "testId",
		Raw: "raw",
	}

	// echo -n [raw] | md5
	m := md5.New()
	m.Write([]byte(data.Raw))
	digest := m.Sum([]byte(nil))
	data.Signature = strings.ToLower(hex.EncodeToString(digest))

	result := v.Validate(data)

	assert.True(t, result)
}
