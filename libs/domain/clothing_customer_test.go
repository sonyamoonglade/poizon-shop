package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidPhoneNumber(t *testing.T) {
	tests := []struct {
		description string
		phoneNumber string
		expected    bool
	}{
		{
			description: "valid 11-digit phone number that starts with 7",
			phoneNumber: "79261234567",
			expected:    true,
		},
		{
			description: "valid 11-digit phone number that starts with 8",
			phoneNumber: "89261234567",
			expected:    true,
		},
		{
			description: "valid 11-digit phone number with spaces and dashes",
			phoneNumber: "8 926-123-45-67",
			expected:    false,
		},
		{
			description: "valid 11-digit phone number with parentheses",
			phoneNumber: "8(926)123-45-67",
			expected:    false,
		},
		{
			description: "phone number with less than 11 digits",
			phoneNumber: "7926123456",
			expected:    false,
		},
		{
			description: "phone number with more than 11 digits",
			phoneNumber: "792612345678",
			expected:    false,
		},
		{
			description: "phone number that starts with +7",
			phoneNumber: "+79261234567",
			expected:    false,
		},
		{
			description: "phone number that starts with 6",
			phoneNumber: "69261234567",
			expected:    false,
		},
		{
			description: "phone number with invalid characters",
			phoneNumber: "8a261234567",
			expected:    false,
		},
	}

	for _, test := range tests {
		actual := IsValidPhoneNumber(test.phoneNumber)
		t.Run(test.description, func(t *testing.T) {
			require.Equal(t, test.expected, actual)
		})
	}
}
