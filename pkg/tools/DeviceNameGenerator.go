package DeviceNameGenerator

import "fmt"
import "math/rand"
import "strconv"
import "strings"
import "time"

var increment = 0;

func pad(v string, n int) string {
	if (len(v) >= n) {
		return v
	}
	return strings.Repeat("0", n - len(v)) + v
}

func dec2bin(dec int64) string {
	return strconv.FormatInt(dec, 2);
}

func generate() string {
	loc, _ := time.LoadLocation("UTC")
	t := time.Now().Sub(time.Date(2009, 1, 1, 0, 0, 0, 0, loc))
	if (increment > 1023) {
		increment = 0;
	}
	// TODO: Write Base62 encoding
	return "01" + pad(dec2bin(t.Nanoseconds()/time.Second.Nanoseconds()), 32) + "I" + pad(dec2bin(int64(rand.Intn(16))), 4) + pad(dec2bin(int64(increment + 1)), 10);
}
