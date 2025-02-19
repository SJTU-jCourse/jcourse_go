package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"jcourse_go/constant"
)

func TestGenerateVerifyCode(t *testing.T) {
	code, err := generateVerifyCode()
	assert.Equal(t, len(code), constant.AuthVerifyCodeLen)
	assert.Nil(t, err)
}
