package testing

type Retrieve struct {
}

func (Retrieve) Get(string) string {
	return "this is algorithm testing default return"
}
