package connector

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rusldv/statefix/utillib"

	"github.com/rusldv/statefix/block"
)

// Start shedule
func Start(cfg *utillib.Config) {
	// Синхронизируем блокчейн по хешу его блока генезиса
	fmt.Println("Синхронизация блокчейна", cfg.GenesisHash)
	if cfg.Peers == "" {
		fmt.Println("В файле конфигурации отсутствуют инициализирующие пиры")
		os.Exit(1)
	}
	peers := strings.Split(cfg.Peers, ",")
	//fmt.Println(peers)
	// Коннектимся к пирам, заданным в конфиге
	for _, peer := range peers {
		conn := connectTo(peer)
		fmt.Println(conn)
		// TODO тут мы вызываем функции из connect_lib.go
		// передаем им conn и они начинают обмен инвентарем
		// Отвечать можно прямо функциями из пакета server
		// в этом же пакете можно взять функцию NewInv для создания инвентаря с запросами
	}

}

// Broadcast разсылает блок соседям
func Broadcast(b *block.Block, peers []string) error {
	fmt.Println(peers)
	fmt.Println("Рассылка блока", b.GetHash(), "...")
	//blockBencode := b.ToBencode()
	//fmt.Println(string(blockBencode))
	for _, peer := range peers {
		conn := connectTo(peer)
		fmt.Println(conn)
	}

	return nil
}

func connectTo(peer string) net.Conn {
	fmt.Println("connect to", peer)
	conn, err := net.Dial("tcp", peer)
	if err != nil {
		fmt.Printf("Error connection to %s\n", peer)
		return nil
	}
	return conn
}
