package bcrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	password := "super secret password"
	hash, err := BcryptPassword(password)
	assert.NoError(t, err)
	t.Logf("hash = %q", hash)
	err = CompareHashAndPassword(password, hash)
	assert.NoError(t, err)
	err = CompareHashAndPassword("Another secret password", hash)
	assert.ErrorIs(t, err, ErrMismatchedHashAndPassword)
}
