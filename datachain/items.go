package datachain

import (
	"bytes"
	"fmt"

	"github.com/jackpal/bencode-go"
)

// Item содержит объект элемента
type Item struct {
	SelfTxAddress    string   `bencode:"self_tx_address"`
	InputTxAddresses []string `bencode:"input_tx_addresses"`
}

// NewItem создает новый элемент
func NewItem(selfAddr string) *Item {
	return &Item{
		SelfTxAddress: selfAddr,
	}
}

// FromBencode декодирует из бинкода
func FromBencode(data []byte) (*Item, error) {
	var it Item
	r := bytes.NewReader(data)
	err := bencode.Unmarshal(r, &it)
	if err != nil {
		return nil, err
	}
	return &it, nil
}

// PushInput добавляет новый вход (полный адрес ссылающейся странзакции)
func (it *Item) PushInput(fullAddr string) {
	it.InputTxAddresses = append(it.InputTxAddresses, fullAddr)
}

// ToBencode кодирует в бинкод
func (it *Item) ToBencode() []byte {
	buf := new(bytes.Buffer)
	err := bencode.Marshal(buf, *it) // Маршалит по значению только
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return buf.Bytes()
}

// Общие функции для работы с БД

// AddItem добавляет элемент в базу данных

// GetItem получает из базы элемент

// DelItem удаляет элемент из базы
