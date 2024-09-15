package random

import (
	crypto "crypto/rand"
	"encoding/hex"
	"math"
	"math/big"
	"math/rand"
)

// IntInRange returns a random number in the range [min, max)
func IntInRange(min, max int) int {
	return min + rand.Intn(max-min)
}

// IntInSize returns random number in size
func IntInSize(size int) int {
	maxLimit := int64(int(math.Pow10(size)) - 1)
	lowLimit := int(math.Pow10(size - 1))

	randomNumber, err := crypto.Int(crypto.Reader, big.NewInt(maxLimit))
	if err != nil {
		return 0
	}
	randomNumberInt := int(randomNumber.Int64())

	// Handling integers between 0, 10^(n-1) .. for n=4, handling cases between (0, 999)
	if randomNumberInt <= lowLimit {
		randomNumberInt += lowLimit
	}

	// Never likely to occur, kust for safe side.
	if randomNumberInt > int(maxLimit) {
		randomNumberInt = int(maxLimit)
	}
	return randomNumberInt
}

// String generate random string
func String(length int) string {
	buff := make([]byte, int(math.Ceil(float64(length)/2)))
	_, _ = crypto.Read(buff)
	str := hex.EncodeToString(buff)
	return str[:length] // strip 1 extra character we get from odd length results
}

// ID generate random string
func ID() string {
	return "_" + String(8)
}
