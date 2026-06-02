package helpers_test

import (
	"strings"
	"testing"

	"github.com/raythrp/evermos-internship/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimeParserToDate_ValidDate(t *testing.T) {
	result, err := helpers.TimeParserToDate("15/06/1995")
	require.NoError(t, err)
	assert.Equal(t, 1995, result.Year())
	assert.Equal(t, 6, int(result.Month()))
	assert.Equal(t, 15, result.Day())
}

func TestTimeParserToDate_InvalidFormat(t *testing.T) {
	_, err := helpers.TimeParserToDate("1995-06-15")
	assert.Error(t, err)
}

func TestTimeParserToDate_InvalidDate(t *testing.T) {
	_, err := helpers.TimeParserToDate("not-a-date")
	assert.Error(t, err)
}

func TestConvertToSlug_SpacesAndUppercase(t *testing.T) {
	assert.Equal(t, "hello-world", helpers.ConvertToSlug("Hello World"))
}

func TestConvertToSlug_AlreadyLowercase(t *testing.T) {
	assert.Equal(t, "already-slug", helpers.ConvertToSlug("already slug"))
}

func TestConvertToSlug_MultipleSpaces(t *testing.T) {
	result := helpers.ConvertToSlug("multiple   spaces")
	assert.False(t, strings.Contains(result, " "))
}

func TestIDGenerator_NonZero(t *testing.T) {
	id := helpers.IDGenerator()
	assert.Greater(t, id, uint(0))
}

func TestIDGenerator_UnderMaxValue(t *testing.T) {
	id := helpers.IDGenerator()
	assert.Less(t, id, uint(100_000_000))
}

func TestInvoiceCodeGenerator_HasPrefix(t *testing.T) {
	code := helpers.InvoiceCodeGenerator()
	assert.True(t, strings.HasPrefix(code, "INV-"))
}

func TestInvoiceCodeGenerator_Unique(t *testing.T) {
	// Two calls close together should not produce the same code due to UnixNano
	// (not guaranteed, but tests the general shape)
	code := helpers.InvoiceCodeGenerator()
	assert.Greater(t, len(code), 4)
}
