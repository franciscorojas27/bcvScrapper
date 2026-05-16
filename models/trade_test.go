package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONBValue(t *testing.T) {
	var empty JSONB

	value, err := empty.Value()
	require.NoError(t, err)
	assert.Equal(t, "[]", value)

	data := JSONB{"one", "two"}
	value, err = data.Value()
	require.NoError(t, err)

	encoded, ok := value.([]byte)
	require.True(t, ok)
	assert.JSONEq(t, `["one","two"]`, string(encoded))
}

func TestJSONBScan(t *testing.T) {
	var data JSONB
	err := data.Scan([]byte(`["one","two"]`))
	require.NoError(t, err)
	assert.Equal(t, JSONB{"one", "two"}, data)
}

func TestJSONBScanRejectsInvalidType(t *testing.T) {
	var data JSONB
	err := data.Scan("invalid")
	require.Error(t, err)
	assert.Equal(t, "invalid scan type for JSONB", err.Error())

	_, _ = json.Marshal(data)
}
