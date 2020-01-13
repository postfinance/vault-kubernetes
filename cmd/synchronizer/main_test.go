package main

import (
	"encoding/base32"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	quote = "The fool don‘t think he is wise, but the wise man knows himself to be a fool."
)

func TestDecode(t *testing.T) {
	t.Run("not encoded", func(t *testing.T) {
		res, err := decode(quote)
		assert.NoError(t, err)
		assert.Equal(t, quote, string(res))
	})

	t.Run("base64 encoded", func(t *testing.T) {
		str := "base64:" + base64.StdEncoding.EncodeToString([]byte(quote))
		res, err := decode(str)
		t.Log(str)
		t.Log(string(res))
		assert.NoError(t, err)
		assert.Equal(t, quote, string(res))
	})

	t.Run("base64 decode fails", func(t *testing.T) {
		str := "base64:" + quote
		res, err := decode(str)
		assert.Error(t, err)
		assert.Equal(t, str, string(res))
	})

	t.Run("unknown encoding", func(t *testing.T) {
		str := "base32:" + base32.StdEncoding.EncodeToString([]byte(quote))
		res, err := decode(str)
		assert.Error(t, err)
		assert.Equal(t, str, string(res))
	})
}
