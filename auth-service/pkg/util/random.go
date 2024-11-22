package util

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const qwerty = "qwertyuiopasdfghjklzxcvbnm"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(qwerty)

	for i := 0; i < n; i++ {
		c := qwerty[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomCpfCnpj() string {
	n := RandomInt(0, 99999999999999)

	return strconv.FormatInt(n, 10)
}

func RandomMoney() int64 {
	return RandomInt(0, 10000.00*100)
}
