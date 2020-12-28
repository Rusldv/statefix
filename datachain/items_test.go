package datachain

import (
	"fmt"
	"testing"
)

func TestNewItems(t *testing.T) {
	it := NewItem("ff")
	it.PushInput("f1")
	it.PushInput("f2")
	it.PushInput("f3")

	fmt.Println(it)

	itBen := it.ToBencode()
	fmt.Println(string(itBen))

	it2, _ := FromBencode(itBen)
	fmt.Println(it2)
}
