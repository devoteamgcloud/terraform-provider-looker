package acceptance_tests

import (
	"math/rand"
	"strings"
	"time"
)

// Function to generate a random string of specified length
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyz"
	var builder strings.Builder
	for i := 0; i < length; i++ {
		builder.WriteByte(charset[rand.Intn(len(charset))])
	}
	return builder.String()
}
