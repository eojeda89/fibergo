package security

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	id := "132165468798"
	token, err := NewToken(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken(t *testing.T) {
	id := "132165468798"
	token, err := NewToken(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, payload.Id, id)
	assert.Equal(t, payload.Issuer, id)
	assert.Equal(t, time.Unix(payload.IssuedAt, 0).Year(), time.Now().Year())
	assert.Equal(t, time.Unix(payload.IssuedAt, 0).Month(), time.Now().Month())
	assert.Equal(t, time.Unix(payload.IssuedAt, 0).Day(), time.Now().Day())
}
