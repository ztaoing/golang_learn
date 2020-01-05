package main

func Stop(stop <-chan bool) {
	close(stop)
}

/**
考点:close channel
有方向的channel不可被关闭
*/
