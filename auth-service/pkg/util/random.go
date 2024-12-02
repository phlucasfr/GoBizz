package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const qwerty = "qwertyuiopasdfghjklzxcvbnm"

var randGen *rand.Rand

func init() {
	randSource := rand.NewSource(time.Now().UnixNano())
	randGen = rand.New(randSource)
}

func RandomInt(min, max int64) int64 {
	return min + randGen.Int63n(max-min+1)
}

func GenerateResetToken(email string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%x", timestamp) + ":" + email
}

func GenerateRandomCode() string {
	code := randGen.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(qwerty)

	for i := 0; i < n; i++ {
		sb.WriteByte(qwerty[randGen.Intn(k)])
	}

	return sb.String()
}

func RandomCpfCnpj() string {
	n := RandomInt(0, 99999999999999)
	return strconv.FormatInt(n, 10)
}

func RandomMoney() int64 {
	return RandomInt(0, 1000000)
}
