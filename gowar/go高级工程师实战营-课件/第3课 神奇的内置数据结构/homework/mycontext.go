package main

import "context"

type MyContext struct{}

func (m *MyContext) WithValue(parent context.Context, key, val interface{}) context.Context {
	return nil
}

func main() {}
