package mocks

type FakeGame struct {
	hash string
}

func (f *FakeGame) clone() (result *FakeGame) {
	result = &FakeGame{}
	*result = *f
	return result
}
