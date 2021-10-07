package generate_test

import (
	"testing"

	"github.com/justasable/pgmock/internal/generate"
	"github.com/stretchr/testify/assert"
)

func TestText(t *testing.T) {
	expected := []string{
		"hello world",
		"3?!-+@.(\x01)\u00f1\u6c34\ubd88\u30c4\U0001f602",
	}
	got := generate.Text()
	assert.Equal(t, expected, got)
}
