package instancer

import (
	"crypto/rand"
	"math/big"
)

func (inst *Instancer) GenerateName() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	const length = 10

	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}

	str := inst.Prefix + string(b)
	return str, nil
}
