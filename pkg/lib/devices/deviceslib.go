package deviceslib

import (
	"core/pkg/tools"
	"math/rand"
	"time"
)

var increment = 1

const OFFSET_INCREMENT = 0
const OFFSET_RANDOM = 10
const OFFSET_DATE = 22
const OFFSET_PREFIX = 63

const PREFIX = 1
const RANDOM_MAX = 4096

var RandomGenerator *rand.Rand

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	RandomGenerator = rand.New(source)
}

func GenerateDeviceUID() uint64 {
	// first part: increment
	output := uint64(increment << OFFSET_INCREMENT)

	// second part: random
	output = output | uint64(RandomGenerator.Intn(RANDOM_MAX)<<OFFSET_RANDOM)

	// third part: date
	loc, _ := time.LoadLocation("UTC")
	t := time.Now().Sub(time.Date(2017, 1, 1, 0, 0, 0, 0, loc))
	output = output | uint64((t.Nanoseconds()/time.Millisecond.Nanoseconds())<<OFFSET_DATE)

	// fourth part: prefix
	output = output | uint64(PREFIX<<OFFSET_PREFIX)

	// increment the increment
	increment = (increment + 1) % 1024

	return output
}

func GenerateDeviceUIDString() string {
	uid := GenerateDeviceUID()
	return tools.EncodeUInt64(uid)
}
