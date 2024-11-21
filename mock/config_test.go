package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"syreclabs.com/go/faker"
)

func TestSetSeed(t *testing.T) {
	t.Run("SetSeed", func(t *testing.T) {
		SetSeed(1)
		avatar1 := faker.Avatar()
		name1 := faker.Name()
		SetSeed(1)
		avatar2 := faker.Avatar()
		name2 := faker.Name()
		assert.Equal(t, avatar1, avatar2)
		assert.Equal(t, name1, name2)
	})
}
