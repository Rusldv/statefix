package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/rusldv/statefix/cryptolib"
	"github.com/rusldv/statefix/utillib"
)

const (
	version = "v1"
)

// Field - поле транзакции включающее base64 данные.
type Field struct {
	N     int    `json:"n" bencode:"n"`
	Value string `json:"value" bencode:"value"`
}

// Transaction - транзакция.
type Transaction struct {
	ObjectID  int     `json:"object_id" bencode:"object_id"`
	Type      byte    `json:"type" bencode:"type"`
	Hash      string  `json:"hash" bencode:"hash"`
	PublicKey string  `json:"public_key" bencode:"public_key"`
	Sig       string  `json:"sig" bencode:"sig"`
	BlockHash string  `json:"block_hash" bencode:"block_hash"`
	Timestamp int64   `json:"timestamp" bencode:"timestamp"`
	Fields    []Field `json:"fields" bencode:"fields"`
}

// NewTransaction возвращает указатель на объект транзакции.
func NewTransaction() *Transaction {
	return &Transaction{
		ObjectID:  2,
		Hash:      version,
		Sig:       version,
		BlockHash: version,
		Timestamp: time.Now().Unix(),
	}
}

// FromJSON декодирует и возвращает указатель на транзакцию из JSON.
func FromJSON(data []byte) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal(data, &tx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tx, nil
}

// FromBencode декодирует и возвращает указатель на транзакцию из Bencode.
func FromBencode(data []byte) (*Transaction, error) {
	var tx Transaction
	r := bytes.NewReader(data)
	err := bencode.Unmarshal(r, &tx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tx, nil
}

// SetType (t) - устанавливает тип транзакции.
func (tx *Transaction) SetType(t byte) {
	tx.Type = t
}

// SetHash (hash) - устанавливает хеш транзакции.
func (tx *Transaction) SetHash(hash string) {
	tx.Hash = hash
}

// SetPublicKey (pub) - устанавливает публичный ключ аккаунта подписи транзакции.
func (tx *Transaction) SetPublicKey(pub string) {
	tx.PublicKey = pub
}

// SetSig (sig) - устанавливает подпись для хеша транзакции.
func (tx *Transaction) SetSig(sig string) {
	tx.Sig = sig
}

// SetBlockHash (hash) - устанавливает хеш заголовка блока, в который добавлена транзакция.
func (tx *Transaction) SetBlockHash(hash string) {
	tx.BlockHash = hash
}

// AddField (value) - добавляет поле с данными в транзакцию.
func (tx *Transaction) AddField(value string) {
	field := Field{
		N:     len(tx.Fields),
		Value: value,
	}
	tx.Fields = append(tx.Fields, field)
}

// Bytes возвращает байтовое представление данных транзакции (для хеширования).
func (tx *Transaction) Bytes() []byte {
	var bfields [][]byte
	for _, v := range tx.Fields {
		bval32, err := utillib.Uint32ToBytes(v.N)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		bfields = append(bfields, bval32)
		bfields = append(bfields, []byte(v.Value))
	}

	bval64, err := utillib.Uint64ToBytes(tx.Timestamp)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return bytes.Join([][]byte{
		[]byte{tx.Type},
		// Поле Hash не учавствует в хешировании
		[]byte(tx.PublicKey),
		// Поле Sig не учавствует в хешировании
		// Поле BlockHash не учавствует в хешировании
		bval64,
		bytes.Join(bfields, []byte{}),
	}, []byte{})
}

// BytesText возвращает байтовое представление транзакции в виде текста
func (tx *Transaction) BytesText() []byte {
	var bhash string
	for _, b := range tx.Bytes() {
		n := (int(b))
		bhash = bhash + strconv.Itoa(n)
	}

	return []byte(bhash)
}

// ToJSON кодирует транзакцию в JSON.
func (tx *Transaction) ToJSON() []byte {
	data, err := json.Marshal(tx)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

// ToBencode кодирует транзакцию в Bencode.
func (tx *Transaction) ToBencode() []byte {
	buf := new(bytes.Buffer)
	err := bencode.Marshal(buf, *tx) // Маршалит по значению только
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return buf.Bytes()
}

// Check - проверяет транзакцию с байтовым хешем.
func (tx *Transaction) Check() bool {
	if tx.Hash != utillib.BytesToSHA256Hex(tx.Bytes()) {
		return false
	}
	// Проверка подписи транзакции
	x, y, err := cryptolib.PublicKeyUnmarshalHex(tx.PublicKey)
	if err != nil {
		fmt.Println(err)
		return false
	}
	pub := cryptolib.XYToPublicKey(x, y)
	vf := cryptolib.VerifyDataASN1Hex(pub, tx.Sig, utillib.BytesToSHA256(tx.Bytes()))
	if vf != true {
		return false
	}

	// TODO другие проверки

	return true
}

// CheckText - проверяет транзакцию c текстовым представлением байтового хеша.
func (tx *Transaction) CheckText() bool {
	//fmt.Println(string(tx.BytesText()))
	if tx.Hash != utillib.BytesToSHA256Hex(tx.BytesText()) {
		fmt.Println("tx.Hash:", tx.Hash)
		fmt.Println("runHash:", utillib.BytesToSHA256Hex(tx.BytesText()))
		fmt.Println("Вычеслиный хеш транзакции не совпадает с установленным")
		return false
	}

	// Проверка подписи транзакции
	x, y, err := cryptolib.PublicKeyUnmarshalHex(tx.PublicKey)
	if err != nil {
		fmt.Println(err)
		return false
	}
	pub := cryptolib.XYToPublicKey(x, y)
	vf := cryptolib.VerifyDataASN1Hex(pub, tx.Sig, utillib.BytesToSHA256(tx.BytesText()))
	if vf != true {
		return false
	}

	// TODO другие проверки

	// miner.Insert(tx)

	return true
}
