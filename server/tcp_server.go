package server

import (
	"fmt"
	"net"

	"github.com/rusldv/statefix/miner"
	"github.com/rusldv/statefix/transaction"

	"github.com/rusldv/statefix/connector"
	"github.com/rusldv/statefix/datachain"

	"github.com/rusldv/statefix/block"

	"github.com/rusldv/statefix/protocol"

	"github.com/rusldv/statefix/utillib"
)

const (
	// Inv инвентарь
	Inv = 1
	// Tx транзакция
	Tx = 2
	// Block блок
	Block = 3
)

type server struct{}

// TCPStart запускает TCP сервер
func TCPStart(cfg *utillib.Config) error {
	port := cfg.TCPPort
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	defer l.Close()
	fmt.Println("TCP сервер запущен на порту", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			conn.Close()
			fmt.Printf("Ошибка получения входящего соединения %v\n", err)
			continue
		}
		fmt.Printf("Установлено входящее соединение с %v\n", conn.RemoteAddr())
		go handle(conn, cfg)
	}
}

func handle(conn net.Conn, cfg *utillib.Config) {
	defer conn.Close()
	data, n, err := utillib.TCPReceiveObject(conn)
	if err != nil {
		fmt.Println(err)
		invError := protocol.NewInv(1, protocol.ErrorState)
		invErrorBencode, _ := invError.ToBencode()
		utillib.TCPSendObject(conn, invErrorBencode)
	}
	fmt.Println("Readed", n, "bytes")
	route(conn, data, cfg)
}

func route(conn net.Conn, data []byte, cfg *utillib.Config) {
	fmt.Println("route")
	//fmt.Println(cfg)
	//str := string(data)
	//fmt.Println(str)
	// Тут с помощью строковой функции выпиливаем object_id
	objid, err := utillib.TCPObjectID(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Проверяем тип полученного объекта
	if objid == Inv {
		inv, err := protocol.FromBencode(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println("inv.Type:", inv.Type, protocol.OKState)
		// Имена констант соответствуют названиям функций в tcp_lib
		switch inv.Type {
		case protocol.GetVersion:
			GetVersion(conn)
		case protocol.OK:
			fmt.Println("OK")
		case protocol.State:
			fmt.Println("State")
		case protocol.ErrorState:
			fmt.Println("ErrorState")
		case protocol.GetBlock:
			fmt.Println("GetBlock")
		case protocol.GetHeight:
			GetHeight(conn)
		case protocol.GetLastHash:
			GetLastHash(conn)
		case protocol.GetState:
			GetState(conn)
		default:
			fmt.Println("Unknown")
		}

	} else if objid == Tx {
		// При получении транзакции проверяем ее и валидируем перед добавлением в блок и пересылкой соседним нодам
		// Также нужно занести ее в хеш таблицу, чтобы заново не добавить в блок
		// Проверяем правильность времени оно не должно быть сбито более чем на 20 сек
		fmt.Println("Tx")
		tx, err := transaction.FromBencode(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !tx.Check() {
			fmt.Println("Принятая транзакция не одобрена")
			return
		}
		//fmt.Println(tx)
		// Включен ли режим майнинга
		if miner.Mine == true { // TODO засунуть в конфиг
			// Если да, то добавляем в пул
			miner.Push(tx)
			miner.Show() // TODO test
		}

		// Рассылка соседним нодам через connector
		//fmt.Println("sending", data)
	} else if objid == Block {
		// При получении блока вызываем функцию консенсуса для его проверки перед принятием и рассылкой дальше
		// Проверять наличие форков цепочки через инвентарь
		fmt.Println("Block")
		// Создаем оъект блока из bencode байтов
		block, err := block.FromBencode(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(block)
		// Выполняем проверку Check
		if !block.Check() {
			fmt.Println("Block not checking")
			return
		}
		// Другие проверки на валидность
		fmt.Println("Block - ok")
		// Проверяем по его хэшу есть ли он уже в БД
		bHash := block.GetHash()
		fmt.Println("bHash:", bHash)
		state := datachain.GetState()
		fmt.Println("ChainID:", state.ChainID)
		fmt.Println("LastHash:", state.LastHash)
		fmt.Println("LastHeight:", state.LastHeight)
		//fmt.Println("state:", state)
		// Проверка соответствия ChainID
		fmt.Println("state.ChainID:", state.ChainID, "block.ChainID:", block.ChainID)
		if block.ChainID != state.ChainID {
			fmt.Println("Ошибка: ID цепочек разные")
			return
		}
		fmt.Println("OK")
		// Выполняем поиск блока в БЧ чтобы его там не было
		/*
			b, _ := datachain.GetBlock(block.GetHash())
			if b != nil {
				fmt.Println("Такой блок уже есть в блокчейне!")
				return
			}
		*/
		// Если нет сохраняем к себе в локальный БЧ с помощью datachain
		/*
			fmt.Println("Добавить блок в блокчейн")
			err = datachain.AddBlock(block)
			if err != nil {
				fmt.Println("Ошибка добавления блока в локальный блокчейн")
				return
			}
			fmt.Println("Блок сохранен в локальный блокчейн.")
		*/
		// И рассылаем соседним нодам с помощью connector
		err = connector.Broadcast(block, utillib.GetPeers(cfg.Peers))
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		fmt.Println("Unknown object id")
		return
	}

	/*
		// Отвечаем что все ок
		invError := protocol.NewInv(1, protocol.OKState)
		invErrorBencode, _ := invError.ToBencode()
		utillib.TCPSendObject(conn, invErrorBencode)
	*/

	return
}
