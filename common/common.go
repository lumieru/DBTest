package common

import (
	"net"
	"log"
	"time"
	"math/rand"
)

var (
	RowsInserted uint32 = 0
)

func GetListener() net.Listener {
	ln, err := net.Listen("tcp4", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	return ln
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImprSrc(n int) []byte {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandBool() bool {
	if rand.Uint32() % 2 == 0 {
		return true
	} else {
		return false
	}
}

func RandInt32() int32 {
	return rand.Int31()
}

func RandInt32Array() []int32 {
	res := make([]int32, rand.Uint32() % 32)
	for i:=0; i<len(res); i++ {
		res[i] = rand.Int31()
	}

	return res
}

func RandID() int32 {
	if RowsInserted == 0 {
		return 0
	} else {
		return (int32)(rand.Uint32()%RowsInserted + 1)
	}
}