package protocol

import (
	"fmt"
	"testing"
)

func TestNewInv(t *testing.T) {
	inv := NewInv(1, 2)
	inv.SetContent("test")
	fmt.Println(inv)
	invJSON, _ := inv.ToJSON()
	fmt.Println(string(invJSON))
	invBencode, _ := inv.ToBencode()
	fmt.Println(string(invBencode))
	fmt.Println("--- from json")
	fj, _ := FromJSON(invJSON)
	fmt.Println(fj)
	fben, _ := FromBencode(invBencode)
	fmt.Println(fben)
}
