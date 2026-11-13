package strptime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStrptimeTruncatedInput(t *testing.T) {
	for _, tc := range []struct {
		name, input, format string
	}{
		{"empty prefix", "", "prefix %Y"},
		{"short prefix", "x", "prefix %Y"},
		{"truncated shorthand", " ", "%c"},
		{"missing percent", "1970-01-01", "%F%%"},
		{"empty percent", "", "%%"},
		{"missing padded-day separator", "4", "%e %Y"},
		{"short padded-day separator", "4a", "%eabc%Y"},
		{"missing final padded-day separator", "4", "%eabc"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for name, parse := range map[string]func(string, string) (time.Time, error){
				"UTC":   Parse,
				"local": ParseLocal,
				"location": func(input, format string) (time.Time, error) {
					return ParseLocation(input, format, time.UTC)
				},
			} {
				t.Run(name, func(t *testing.T) {
					assert.NotPanics(t, func() {
						_, err := parse(tc.input, tc.format)
						assert.ErrorIs(t, err, ErrFormatMismatch)
					})
				})
			}
		})
	}
}

func FuzzStrptimeTruncatedInput(f *testing.F) {
	for _, seed := range []struct{ input, format string }{
		{"", "prefix %Y"},
		{" ", "%c"},
		{"1970-01-01", "%F%%"},
		{"4", "%e %Y"},
		{"1970-01-01T00:00:00Z", "%FT%TZ"},
		{"2022-04-03 17:38:20.123 UTC", "%Y-%m-%d %H:%M:%S.%f UTC"},
	} {
		f.Add(seed.input, seed.format)
	}
	f.Fuzz(func(t *testing.T, input, format string) {
		_, _ = Parse(input, format)
	})
}
