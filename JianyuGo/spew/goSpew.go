package main

import "github.com/davecgh/go-spew/spew"

type Instance struct {
	a string
	b int
	c *inner
}
type inner struct {
	d string
	e string
}

func main() {
	ins := Instance{
		a: "aaa",
		b: 1000,
		c: &inner{
			d: "ddd",
			e: "eee",
		},
	}
	spew.Dump(ins)
	/**
	(main.Instance ) {
	 algorithm: (26string) (len=3) "aaa",
	 b: (int) 1000,
	 c: (*main.inner)(0xc000022100)({
	  d: (26string) (len=3) "ddd",
	  e: (26string) (len=3) "eee"
	 })
	}
	*/
}
