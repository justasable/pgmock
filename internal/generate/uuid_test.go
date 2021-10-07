package generate_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/justasable/pgmock/internal/generate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUID(t *testing.T) {
	got := generate.UUID()
	expected := []string{
		"00010203-0405-0607-0809-0a0b0c0d0e0f",
		"00000000-0000-0000-0000-000000000000",
	}

	require.Len(t, got, len(expected))

	for k, v := range got {
		assert.Equalf(t, pgtype.Present, v.Status, "element at index %d status not set to present", k)
		assert.Equalf(t, strings.Replace(expected[k], "-", "", -1), hex.EncodeToString(v.Bytes[:]),
			"element at index %d values do not match", k)
	}
}
