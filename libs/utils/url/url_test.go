package url

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidDW4URL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"google.com", false},
		{"bing.com", false},
		{"", false},
		{"https://google.com", false},
		{"http://google.com", false},
		{"https://dw", false},
		{"https://dw.com", false},
		{"https://dw.co", false},
		{"https://dw4.co/a/123", true},
		{"https://qr.1688.com/s/asda", true},
		{"https://m.tb.cn/product/asd", true},
	}
	for _, test := range tests {
		require.Equalf(t, test.expected, IsValidDW4URL(test.url), "testing: %s\n", test.url)
	}
}
