package util_test

import (
	"jcourse_go/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegWord(t *testing.T) {
	err := util.InitSegWord()
	assert.NoError(t, err)
	const txt = "电路理论"
	var target = []string{"电路", "理论"}
	assert.Equal(t, util.SegWord(txt), target)
}
