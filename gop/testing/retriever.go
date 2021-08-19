package testing

type Retrieve struct {
}

func (Retrieve) Get(string) string {
	return "this is a testing default return"
}
