package generate_test

import (
	"testing"

	"github.com/justasable/pgmock/internal/generate"
	"github.com/stretchr/testify/assert"
)

func TestInteger(t *testing.T) {
	expected := []int{0, 1, -1, 2147483648, -2147483648}
	got := generate.Integer()

	assert.Equal(t, got, expected)
}
