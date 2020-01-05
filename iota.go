package main

/**
考点：iota
*/
const (
	x = iota
	y
	z = "zz"
	k
	p = iota
)

func main() {
	fmt.Println(x, y, z, k, p)
}

/**
输出结果:
0 1 zz zz 4
*/
