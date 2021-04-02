package mediator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallbackType(t *testing.T) {
	callback := NewCallback(After, "")
	if callback != nil {
		t.Error("invalid callback return, maybe is nil for invalid handler")
	}
	assert.Nil(t, callback)
	callback = NewCallback(After, func () {})
	assert.NotNil(t, callback)
	assert.Equal(t, After, callback.Type)
}