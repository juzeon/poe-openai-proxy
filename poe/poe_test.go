package poe

import (
	"errors"
	"testing"

	poe_api "github.com/lwydyby/poe-api"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			assert.Equal(t, true, errors.As(err.(error), &invalidError))
		}
	}()
	poe_api.NewClient("", nil)
}
