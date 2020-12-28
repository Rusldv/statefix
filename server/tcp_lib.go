package server

import (
	"fmt"
	"net"
	"strconv"

	"github.com/rusldv/statefix/datachain"

	"github.com/rusldv/statefix/protocol"
	"github.com/rusldv/statefix/utillib"
)

const version = 1

// NewInv создает объект инвентаря
func NewInv(t int, content string) []byte {
	inv := protocol.NewInv(version, t)
	inv.SetContent(content)
	invBencode, err := inv.ToBencode()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return invBencode
}

// GetVersion send version
func GetVersion(conn net.Conn) {
	fmt.Println("Called Getversion")
	// Создаем инвентарь
	invBencode := NewInv(protocol.Version, "1.0.0")
	// Отправляем соединению
	fmt.Println("Sending:", string(invBencode))
	utillib.TCPSendObject(conn, invBencode)
}

// GetState возвращает соединению состояние ноды через двоеточие
func GetState(conn net.Conn) {
	fmt.Println("Called GetState")

	state := datachain.GetState()
	if state == nil {
		fmt.Println("Пустой запрос состояния блокчейна [datachain.GetState] в server/tcp_lib.go")
		return
	}
	stateJoin := strconv.FormatInt(state.ChainID, 10) + ":" + strconv.FormatInt(state.LastHeight, 10) + ":" + state.LastHash
	// resp
	invBencode := NewInv(protocol.State, stateJoin)
	fmt.Println(string(invBencode))
	utillib.TCPSendObject(conn, invBencode)
}

// GetHeight возвращает высоту блокчейна
func GetHeight(conn net.Conn) {
	fmt.Println("Called GetHeight")
	state := datachain.GetState()
	//fmt.Println(state.LastHeight)
	invBencode := NewInv(protocol.State, strconv.FormatInt(state.LastHeight, 10))
	fmt.Println(string(invBencode))
	utillib.TCPSendObject(conn, invBencode)
}

// GetLastHash возвращает хеш последнего блока
func GetLastHash(conn net.Conn) {
	fmt.Println("Called GetLastHash")
	hash, err := datachain.LastBlockHash()
	if err != nil {
		fmt.Println(err)
		return
	}
	invBencode := NewInv(protocol.State, hash)
	fmt.Println(string(invBencode))
	utillib.TCPSendObject(conn, invBencode)
}
