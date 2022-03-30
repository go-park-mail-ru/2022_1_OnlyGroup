package sessionGenerator

import "math/rand"

const secretRunes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

type SessionGenerator interface {
	String(size int) string
}

type randomGenerator struct {
	runes string
}

func NewRandomGenerator() *randomGenerator {
	return &randomGenerator{runes: secretRunes}
}

func (generator *randomGenerator) String(size int) string {
	result := ""
	for i := 0; i < size; i++ {
		result += string(generator.runes[rand.Intn(len(generator.runes))])
	}
	return result
}
