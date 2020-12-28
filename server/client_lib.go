package server

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/rusldv/statefix/datachain"
	"github.com/rusldv/statefix/miner"
	"github.com/rusldv/statefix/transaction"
)

// Тестовая функция
func test(payload string) (*WSRes, error) {
	fmt.Println("wslib: test:", payload)
	// Ответ клиенту
	res := WSRes{
		Status: "ok",
		Result: "Test ok: " + payload,
	}

	return &res, nil
}

// Обработка неизвестного запроса
func any(call string) (*WSRes, error) {
	fmt.Println("wslib: default")
	// Ответ клиенту
	res := WSRes{
		Status: "error",
		Result: "Default: function not found: " + call,
	}

	return &res, nil
}

// Запрос блока по его хешу
func getBlock(hash string) (*WSRes, error) {
	fmt.Println("wslib: getBlock:", hash)
	block, err := datachain.GetBlock(hash)
	if err != nil {
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "getBlock error",
		}
		return &res, nil
	}
	blockJSONstr := string(block.ToJSON())
	//fmt.Println(blockJSONstr)
	// Ответ клиенту
	res := WSRes{
		Status: "ok",
		Result: blockJSONstr,
	}

	return &res, nil
}

// Запрос блока по его номеру
func getBlockNum(num string) (*WSRes, error) {
	//fmt.Println("wslib: getBlockNum:", num)
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "getBlockNum error",
		}
		return &res, nil
	}
	len := int64(len(datachain.BlocksPool) - 1)
	//fmt.Println(len, n, len < n)
	if len < n {
		fmt.Println("LIMIT")
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "getBlockNum error limit block num",
		}
		return &res, nil
	}

	hash := datachain.BlocksPool[n]
	//fmt.Println("wslib: hash:", hash)

	block, err := datachain.GetBlock(hash)
	if err != nil {
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "getBlockNum",
		}
		return &res, nil
	}
	blockJSONstr := string(block.ToJSON())
	//fmt.Println(blockJSONstr)
	// Ответ клиенту

	res := WSRes{
		Status: "ok",
		Result: blockJSONstr,
	}

	return &res, nil
}

// Запрос транзакции по ее адресу
func getTransaction(addr string) (*WSRes, error) {
	var res WSRes
	blockHash := addr[:64]
	fmt.Println("blockHash:", blockHash)
	txHash := addr[64:]
	fmt.Println("txHash:", txHash)
	// загружаем блок
	b, err := datachain.GetBlock(blockHash)
	if err != nil {
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "getTransaction error",
		}
		return &res, nil
	}
	// ищем в нем транзакцию
	for i, v := range b.Transactions {
		if v.Hash == txHash {
			//fmt.Println(txHash, i, "OK")
			txInBlock := b.Transactions[i].ToJSON()
			// Ответ клиенту
			res = WSRes{
				Status: "ok",
				Result: string(txInBlock),
			}
		}
	}
	return &res, nil
}

// Принимаем отправленную транзакцию
func sendTx(txb64 string) (*WSRes, error) {
	//fmt.Println("sending:", txb64)
	// Декодируем из base64 в JSON текстовое представление
	data, err := base64.StdEncoding.DecodeString(txb64)
	if err != nil {
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "sendTx error",
		}
		return &res, nil
	}
	//fmt.Println(string(data))
	// Создаем объект транзакции и загружаем в него полученный JSON
	tx, err := transaction.FromJSON([]byte(data))
	if err != nil {
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "sendTx error",
		}
		return &res, nil
	}
	fmt.Println("wslib: tx object:", tx)
	// Проверка
	if !tx.CheckText() {
		// Ответ клиенту
		res := WSRes{
			Status: "error",
			Result: "non checked",
		}
		return &res, nil
	}
	// Добавляем в пул если майнится блок
	if miner.Mine == true { // TODO засунуть в конфиг
		// Если да, то добавляем в пул
		miner.Push(tx)
		miner.Show() // TODO test
	}

	// Передаем в коннектор для пересылки соседям
	// Ответ клиенту
	res := WSRes{
		Status: "ok",
		Result: "checked",
	}
	return &res, nil
}
