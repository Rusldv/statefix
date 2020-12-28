package block

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackpal/bencode-go"
	"github.com/rusldv/statefix/cryptolib"
	tx "github.com/rusldv/statefix/transaction"
	"github.com/rusldv/statefix/utillib"
)

const (
	blockVersion = 1
)

// Block - head and transactions body data.
type Block struct {
	ObjectID     int              `json:"object_id" bencode:"object_id"`
	ChainID      int64            `json:"chain_id" bencode:"chain_id"`
	Height       int64            `json:"height" bencode:"height"`
	Hash         string           `json:"hash" bencode:"hash"`
	PrevHash     string           `json:"prev_hash" bencode:"prev_hash"`
	MerkleRoot   string           `json:"merkle_root" bencode:"merkle_root"`
	PublicKey    string           `json:"public_key" bencode:"public_key"`
	Sig          string           `json:"sig" bencode:"sig"`
	Timestamp    int64            `json:"timestamp" bencode:"timestamp"`
	Version      int              `json:"version" bencode:"version"`
	Transactions []tx.Transaction `json:"transactions" bencode:"transactions"`
}

// NewBlock created new block object.
func NewBlock() *Block {
	return &Block{
		ObjectID:  3,
		Timestamp: time.Now().Unix(),
		Version:   blockVersion,
	}
}

// FromJSON декодирует и возвращает указатель на блок из JSON.
func FromJSON(data []byte) (*Block, error) {
	var b Block
	err := json.Unmarshal(data, &b)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &b, nil
}

// FromBencode декодирует и возвращает указатель наблок из Bencode.
func FromBencode(data []byte) (*Block, error) {
	var b Block
	r := bytes.NewReader(data)
	err := bencode.Unmarshal(r, &b)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &b, nil
}

// SetChainID setting ChainID value.
func (b *Block) SetChainID(id int64) {
	b.ChainID = id
}

// GetChainID getting ChainID value.
func (b *Block) GetChainID() int64 {
	return b.ChainID
}

// SetHeight setting Height value.
func (b *Block) SetHeight(height int64) {
	b.Height = height
}

// GetHeight getting Height value.
func (b *Block) GetHeight() int64 {
	return b.Height
}

// SetHash setting Hash value.
func (b *Block) SetHash(hash string) {
	b.Hash = hash
}

// GetHash getting Hash value.
func (b *Block) GetHash() string {
	return b.Hash
}

// SetPrevHash setting PrevHash value.
func (b *Block) SetPrevHash(hash string) {
	b.PrevHash = hash
}

// GetPrevHash getting PrevHash value.
func (b *Block) GetPrevHash() string {
	return b.PrevHash
}

// SetMerkleRoot setting MerkleRoot value.
func (b *Block) SetMerkleRoot(hash string) {
	b.MerkleRoot = hash
}

// GetMerkleRoot getting MerkleRoot value.
func (b *Block) GetMerkleRoot() string {
	return b.MerkleRoot
}

// SetPublicKey setting PublicKey value.
func (b *Block) SetPublicKey(pub string) {
	b.PublicKey = pub
}

// GetPublicKey getting PublicKey value.
func (b *Block) GetPublicKey() string {
	return b.PublicKey
}

// SetSig setting Sig value.
func (b *Block) SetSig(sig string) {
	b.Sig = sig
}

// GetSig getting Sig value.
func (b *Block) GetSig() string {
	return b.Sig
}

// SetTimestamp setting Timestamp value.
func (b *Block) SetTimestamp(tm int64) {
	b.Timestamp = tm
}

// GetTimestamp getting Timestamp value.
func (b *Block) GetTimestamp() int64 {
	return b.Timestamp
}

// SetVersion setting Version value.
func (b *Block) SetVersion(v int) {
	b.Version = v
}

// GetVersion getting Version value.
func (b *Block) GetVersion() int {
	return b.Version
}

// AddTransaction addition transaction object to this block.
func (b *Block) AddTransaction(tx *tx.Transaction) {
	b.Transactions = append(b.Transactions, *tx)
}

// BytesHeader returned header data of bytes slice.
func (b *Block) BytesHeader() []byte {
	bchid64, err := utillib.Uint64ToBytes(b.ChainID)
	if err != nil {
		fmt.Println(err)
		// TODO return?
	}
	bheight64, err := utillib.Uint64ToBytes(b.Height)
	if err != nil {
		fmt.Println(err)
	}
	btm64, err := utillib.Uint64ToBytes(b.Timestamp)
	if err != nil {
		fmt.Println(err)
	}
	bver32, err := utillib.Uint32ToBytes(b.Version)
	if err != nil {
		fmt.Println(err)
	}
	return bytes.Join([][]byte{
		bchid64,
		bheight64,
		// []byte(b.Hash), // Это поле изменяет текущий хеш
		[]byte(b.PrevHash),
		[]byte(b.MerkleRoot),
		[]byte(b.PublicKey),
		// []byte(b.Sig), // Это поле изменяет текущий хеш
		btm64,
		bver32,
	}, []byte{})
}

// ToJSON кодирует блок в JSON.
func (b *Block) ToJSON() []byte {
	data, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return data
}

// ToBencode кодирует блок в Bencode.
func (b *Block) ToBencode() []byte {
	buf := new(bytes.Buffer)
	err := bencode.Marshal(buf, *b)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return buf.Bytes()
}

// Check проверяет блок на корректность.
func (b *Block) Check() bool {
	// Проверка на корректность хеша
	//println(b.Hash)
	//println(utillib.BytesToSHA256Hex(b.BytesHeader()))
	if b.Hash != utillib.BytesToSHA256Hex(b.BytesHeader()) {
		return false
	}

	// Проверка подписи блока
	x, y, err := cryptolib.PublicKeyUnmarshalHex(b.PublicKey)
	if err != nil {
		fmt.Println(err)
		return false
	}
	pub := cryptolib.XYToPublicKey(x, y)
	vf := cryptolib.VerifyDataASN1Hex(pub, b.Sig, utillib.BytesToSHA256(b.BytesHeader()))
	// return vf
	if vf != true {
		return false
	}

	return true
}
