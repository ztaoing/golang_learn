package main

type MyLock struct {
	// you should have a channel here
}

func (m *MyLock) TryLock() bool {
	return false
}

func main() {}
