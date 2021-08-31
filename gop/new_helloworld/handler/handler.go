package handler

// "handler/"是为了解决名称冲突的问题
const HelloServiceName = "handler/HelloService"

type NewHelloService struct {
}

func (h *NewHelloService) Hello(request string, reply *string) error {
	*reply = "get request " + request
	return nil
}
