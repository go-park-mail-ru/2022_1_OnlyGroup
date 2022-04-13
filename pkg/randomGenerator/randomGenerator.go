package randomGenerator

type RandomGenerator interface {
	String(size int) (string, error)
}
