package main

import (
	"encoding/base32"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	quote         = "The fool don‘t think he is wise, but the wise man knows himself to be a fool."
	trivialString = "h"
)

//nolint:goconst // test
func TestDecode(t *testing.T) {
	t.Run("not encoded", func(t *testing.T) {
		res, err := decode(quote)
		assert.NoError(t, err)
		assert.Equal(t, quote, string(res))
	})

	t.Run("base64 encoded", func(t *testing.T) {
		str := "base64:" + base64.StdEncoding.EncodeToString([]byte(quote))
		res, err := decode(str)
		assert.NoError(t, err)
		assert.Equal(t, quote, string(res))
	})

	t.Run("base64 encoded", func(t *testing.T) {
		str := "base64:" + base64.StdEncoding.EncodeToString([]byte(trivialString))
		res, err := decode(str)
		assert.NoError(t, err)
		assert.Equal(t, trivialString, string(res))
	})

	t.Run("base64 decode fails", func(t *testing.T) {
		str := "base64:" + quote
		_, err := decode(str)
		assert.Error(t, err)
	})

	t.Run("unknown encoding", func(t *testing.T) {
		str := "base32:" + base32.StdEncoding.EncodeToString([]byte(quote))
		res, err := decode(str)
		assert.NoError(t, err)
		assert.Equal(t, str, string(res))
	})
}

func TestSplitLabels(t *testing.T) {
	labels := "s1=batman,s2,s3=superman,s4=,s5,"

	exp := map[string]string{
		"s1": "batman",
		"s3": "superman",
		"s4": "", // guess who? the invisible man.
	}

	res := splitLabels(labels)

	require.True(t, len(exp) == len(res))

	for k, v := range res {
		assert.Equal(t, v, exp[k])
	}
}

func TestMergeLabels(t *testing.T) {
	existing := map[string]string{
		"e1": "batman",
		"e2": "superman",
	}

	configured := map[string]string{
		"c1": "wonder woman",
		"e2": "supergirl",
	}

	exp := map[string]string{
		"e1": "batman",
		"c1": "wonder woman",
		"e2": "supergirl",
	}

	res := mergeLabels(existing, configured)

	require.True(t, len(exp) == len(res))

	for k, v := range res {
		assert.Equal(t, v, exp[k])
	}
}
