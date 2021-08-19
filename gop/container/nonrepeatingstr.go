package main

func lengthOfNonRepeatingSubStr(s string) int {
	// 已经出现过的字符的位置
	// 用空间换时间
	lastOccurred := make(map[rune]int)
	// lastOccurred := make([]int,0xffff)
	// 开始位置
	start := 0
	// 当前最大长度
	maxLength := 0
	//支持中文
	for i, v := range []rune(s) {
		// 当前字符已经出现过，并且在有效范围内
		if lastI, ok := lastOccurred[v]; ok && lastI >= start {
			// start移动
			start = lastI + 1
		}
		// 移动或不移动后的当前长度， 大于，之前的长度
		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}
		// 存储已经访问的字符
		lastOccurred[v] = i
	}
	return maxLength
}
func main() {

}
