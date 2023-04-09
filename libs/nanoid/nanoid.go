package nanoid

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	random = "A1B2C3D4E5F6G7H8"
	size   = 5
)

// GenerateNanoID will generate 6 char sequence from english alphabet and numbers
func GenerateNanoID() string {
	v, _ := gonanoid.Generate(random, size)
	return v
}
