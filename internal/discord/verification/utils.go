package verification

import (
	"fmt"
	"math/rand"
	"time"
)

// 6 Digit code
func genCode6() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(900000) + 100000
	return fmt.Sprintf("%06d", code)
}