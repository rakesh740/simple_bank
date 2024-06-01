package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currency := []string{USD, EUR}

	n := len(currency)

	return currency[rand.Intn(n)]
}

func RandomPassword() string {
	return fmt.Sprintf("%sH%s", RandomString(6), RandomString(6))
}

func RandomName() string {
	return fmt.Sprintf("%s %s", RandomString(6), RandomString(6))
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
