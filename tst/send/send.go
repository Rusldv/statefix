package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/rusldv/kit/fileutil"

	"github.com/rusldv/statefix/miner"
	"github.com/rusldv/statefix/protocol"
	"github.com/rusldv/statefix/utillib"
)

var host = flag.String("h", "127.0.0.1", "Connected host.")
var port = flag.String("p", "11900", "Connected port.")
var f = flag.String("f", "", "Sending data file.")
var inv = flag.Int("inv", 1, "Inv object code.")
var txtest = flag.Bool("txtest", false, "Sending test transaction.")

func main() {
	flag.Parse()
	addr := *host + ":" + *port
	fmt.Println("send", addr, *f)
	fmt.Println("inv:", *inv)
	// Соединение
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Error connection to %s\n", addr)
		return
	}
	defer conn.Close()

	if *txtest {
		cfg, err := utillib.ReadConfig("config.json")
		if err != nil {
			fmt.Println(err)
			return
		}
		tx := miner.NewOneTx(cfg)
		fmt.Println(tx)
		data := tx.ToBencode()
		err = utillib.TCPSendObject(conn, []byte(data))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("OK")

		return
	}

	if len(*f) > 0 {
		// Для отправки объекта предварительно его нужно сохранить в Bencode, а не JSON
		fmt.Println("Sending data from Bencode")
		// Загрузка объекта из Bencode файла
		data, err := fileutil.ReadFileString(*f)
		if err != nil {
			fmt.Println(err)
			return
		}
		t, err := utillib.TCPObjectID([]byte(data))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Type object:", t)
		err = utillib.TCPSendObject(conn, []byte(data))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("OK")
	} else {
		// Создаем инвентарь
		inv := protocol.NewInv(1, *inv)
		//inv.SetContent(*senddata)
		invBencode, _ := inv.ToBencode()
		fmt.Println(string(invBencode))
		// Send object
		err = utillib.TCPSendObject(conn, invBencode)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Called OK")
		fmt.Println("Response:")
		respData, _, err := utillib.TCPReceiveObject(conn)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(respData))
		// используем protocol.FromBencode(respData)
		respInv, err := protocol.FromBencode(respData)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(respInv.Type, respInv.Content)
	}
}
