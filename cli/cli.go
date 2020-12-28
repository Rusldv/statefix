package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"

	"github.com/rusldv/statefix/block"

	"github.com/rusldv/statefix/utillib"

	"github.com/rusldv/kit/fileutil"
)

// Config this
type Config struct {
	TCPPort string `json:"tcp_port"`
	Peers   string `json:"peers"`
}

var host = flag.String("host", "127.0.0.1", "Node address.")
var port = flag.String("port", "11900", "Connected port")
var conf = flag.String("config", "./config.json", "Initial configuration file.")

var cfg Config

func main() {
	flag.Parse()
	addr := *host + ":" + *port
	fmt.Println("Connect to", addr)
	if len(*conf) > 0 {
		f, err := fileutil.ReadFileString(*conf)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(f)

		err = json.Unmarshal([]byte(f), &cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(cfg)
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Error connection to %s\n", addr)
		return
	}
	defer conn.Close()

	/*
		intBytes, _ := utillib.Uint32ToBytes(5)
		conn.Write(intBytes)
		conn.Write([]byte{1, 2, 3, 4, 5})
	*/
	// utillib.TCPSendObject(conn, []byte{11, 12, 13, 14, 15, 200, 255})
	//inv := protocol.NewInv(2, 3)
	//inv.SetContent("hello world")
	//invBencode, _ := inv.ToBencode()
	//fmt.Println(invBencode)
	//utillib.TCPSendObject(conn, invBencode)
	//fmt.Println("OK")

	//fmt.Println("--- connector")
	//peers := strings.Split(cfg.Peers, ",")
	//connector.Run(peers)

	// create inv object
	/*
		inv := protocol.NewInv(1, protocol.GetVersion)
		invBencode, _ := inv.ToBencode()
		fmt.Println(string(invBencode))
		// send
		utillib.TCPSendObject(conn, invBencode)
		// receive
		data, _, err := utillib.TCPReceiveObject(conn)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))
	*/

	// create tx object
	/*
		tx := transaction.NewTransaction()
		txBencode := tx.ToBencode()
		utillib.TCPSendObject(conn, txBencode)
	*/

	// create Block object
	block := block.NewBlock()
	blockBencode := block.ToBencode()
	utillib.TCPSendObject(conn, blockBencode)
}
