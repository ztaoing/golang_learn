package main

import "fmt"

//observer模式（观察者模式）：大家都关注一个特定信息，当这个信息发生改变时，所有人都会收到通知
//例如:redis中pub/sub
type Observer interface {
	Update(event string)
}

type EventSource struct {
	observers []Observer
}

func (e *EventSource) AddObserver(o Observer) {
	e.observers = append(e.observers, o)
}

func (e *EventSource) Publish(event string) {
	for _, o := range e.observers {
		o.Update(event)
	}
}

type AObserver struct {
}

func (a *AObserver) Update(event string) {
	fmt.Printf("A received event %s\n", event)
}

type BObserver struct {
}

func (b *BObserver) Update(event string) {
	fmt.Printf("B received event %s\n", event)
}

type CObserver struct {
}

func (c *CObserver) Update(event string) {
	fmt.Printf("C received event %s\n", event)
}

func main() {
	eventSource := EventSource{}
	eventSource.AddObserver(&AObserver{})
	eventSource.AddObserver(&BObserver{})
	eventSource.AddObserver(&CObserver{})

	eventSource.Publish("whoop")
}
