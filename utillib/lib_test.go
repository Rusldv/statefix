package utillib

import (
	"fmt"
	"testing"
	"time"
)

/*
func TestUint32ToBytes(t *testing.T) {
	fmt.Println(Uint32ToBytes(5))
}

func TestUint64ToBytes(t *testing.T) {
	fmt.Println(Uint64ToBytes(15))
}

func TestBytesToSHA256(t *testing.T) {
	fmt.Println(BytesToSHA256([]byte("test")))
}

func TestBytesToSHA256Hex(t *testing.T) {
	fmt.Println(BytesToSHA256Hex([]byte("test")))
	fmt.Println(BytesToSHA256Hex([]byte("test2")))
	fmt.Println(BytesToSHA256Hex([]byte("test3")))
}
*/

func TestNewSmartAddress(t *testing.T) {
	fmt.Println(NewSmartAddress(time.Now().Unix()))
}
