package deviceslib

import (
	"testing"
)

const TEST_UIDS = 100000

func TestGenerateDeviceUIDString(t *testing.T) {
	found := map[string]int{}

	t.Log("Testing Device UID collision")

	for i := 0; i < TEST_UIDS; i++ {
		if i%1000 == 0 {
			t.Log("Device UID at ", i)
		}
		newuid := GenerateDeviceUIDString()
		if val, ok := found[newuid]; ok {
			//t.Fatal("Got a collision... ", i, " ", newuid, " prev at ", val)
			t.Log("Got a collision... ", i, " ", newuid, " prev at ", val)
		}

		found[newuid] = i
	}
}
