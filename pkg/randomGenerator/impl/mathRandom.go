package impl

import "math/rand"

const secretRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

type mathRandomGenerator struct {
	runes string
}

func NewMathRandomGenerator() *mathRandomGenerator {
	return &mathRandomGenerator{runes: secretRunes}
}

func (generator *mathRandomGenerator) String(size int) (string, error) {
	result := ""
	for i := 0; i < size; i++ {
		result += string(generator.runes[rand.Intn(len(generator.runes))])
	}
	return result, nil
}
