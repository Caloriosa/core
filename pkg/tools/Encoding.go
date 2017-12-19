package tools

const encodeStd = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func EncodeUInt64(n uint64) string {
	var (
		b   = make([]byte, 0)
		rem uint64
	)

	// Progressively divide by base, store remainder each time
	// Prepend as an additional character is the higher power
	for n > 0 {
		rem = n % 62
		n = n / 62
		b = append([]byte{encodeStd[rem]}, b...)
	}

	s := string(b)

	return s
}

