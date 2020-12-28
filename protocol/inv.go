package protocol

import (
	"bytes"
	"encoding/json"

	"github.com/jackpal/bencode-go"
)

// Inv protocol message
type Inv struct {
	ObjectID int    `json:"object_id" bencode:"object_id"`
	Version  int    `json:"version" bencode:"version"`
	Type     int    `json:"type" bencode:"type"`
	Content  string `json:"content" bencode:"content"`
}

// NewInv created Inv object
func NewInv(v, t int) *Inv {
	return &Inv{
		ObjectID: 1,
		Version:  v,
		Type:     t,
	}
}

// FromJSON декодирует и возвращает указатель на инвентарь из JSON.
func FromJSON(data []byte) (*Inv, error) {
	var inv Inv
	err := json.Unmarshal(data, &inv)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	return &inv, nil
}

// FromBencode декодирует и возвращает указатель оъект.
func FromBencode(data []byte) (*Inv, error) {
	var inv Inv
	r := bytes.NewReader(data)
	err := bencode.Unmarshal(r, &inv)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	return &inv, nil
}

// SetContent setting content for Inv object
func (inv *Inv) SetContent(content string) {
	inv.Content = content
}

// ToJSON кодирует инвентарь в JSON.
func (inv *Inv) ToJSON() ([]byte, error) {
	data, err := json.Marshal(inv)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ToBencode кодирует инвентарь в Bencode.
func (inv *Inv) ToBencode() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := bencode.Marshal(buf, *inv)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
