package utillib

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rusldv/statefix/cryptolib"

	"github.com/rusldv/kit/fileutil"
)

// Config main node loaded in the struct from configurate file.
type Config struct {
	TCPPort       string        `json:"tcp_port"`
	Peers         string        `json:"peers"`
	DataPath      string        `json:"data_path"`
	WSPort        string        `json:"ws_port"`
	HTTPPort      string        `json:"http_port"`
	GenesisHash   string        `json:"genesis_hash"`
	MinerInterval time.Duration `json:"miner_interval"`
	ChainID       int64         `json:"chain_id"`
	Version       int           `json:"version"`
	AccountPath   string        `json:"account_path"`
}

// ReadConfig читает файл конфига и возвращает в виде объекта
func ReadConfig(path string) (*Config, error) {
	f, err := fileutil.ReadFileString(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal([]byte(f), &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetPeers берет строку пиров и возвращает их массивом
func GetPeers(peers string) []string {
	arr := strings.Split(peers, ",")
	return arr
}

// Butes

// Uint32ToBytes кодирует uint32 число в байты в Little Endian.
func Uint32ToBytes(v int) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint32(v))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Uint64ToBytes кодирует uint64 число в байты в Little Endian.
func Uint64ToBytes(v int64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint64(v))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Crypto

// BytesToSHA256 кодирует байтовый массив в байтовый хеш
func BytesToSHA256(b []byte) []byte {
	sha2 := sha256.New()
	sha2.Write(b)
	return sha2.Sum(nil)
}

// BytesToSHA256Hex кодирует байтовый массив в строковый хеш
func BytesToSHA256Hex(b []byte) string {
	return hex.EncodeToString(BytesToSHA256(b))
}

// Net

const metaBufferSize = 4 // int32

// TCPReceive принимает объект данных
func TCPReceive(conn net.Conn, bufsiz int) ([]byte, error) {
	buf := make([]byte, bufsiz)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// TCPReceiveObject принимает объект данных определяя их размер из вервых 4 байт
func TCPReceiveObject(conn net.Conn) ([]byte, int, error) {
	meta, err := TCPReceive(conn, metaBufferSize)
	if err != nil {
		return nil, 0, err
	}
	newBufSize := binary.LittleEndian.Uint32(meta)
	data, err := TCPReceive(conn, int(newBufSize))
	if err != nil {
		return nil, 0, err
	}
	return data, int(newBufSize), nil
}

// TCPSendObject отправляет объект данных
func TCPSendObject(conn net.Conn, data []byte) error {
	len := len(data)
	//fmt.Println(len)
	intBytes, err := Uint32ToBytes(len)
	if err != nil {
		return err
	}
	//fmt.Println(intBytes)
	buf := make([]byte, len+metaBufferSize)
	copy(buf, intBytes)
	copy(buf[4:], data)

	_, err = conn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

// TCPObjectID вычисляет ID объекта
func TCPObjectID(data []byte) (int, error) {
	str := string(data)
	//fmt.Println(str)
	idx := strings.Index(str, ":object_idi")
	//fmt.Println(idx)
	if idx == -1 {
		return 0, errors.New("not found object_id field")

	}
	cut := idx + len(":object_idi")
	num := str[cut : cut+1]
	//fmt.Println(num)
	objtype, err := strconv.Atoi(num)
	if err != nil {
		return 0, err
	}
	return objtype, nil
}

// GetPath путь до директории с данными БЧ
func GetPath(path string) string {
	if path == "default" {
		// For windows: C:\Users\ПК\AppData\Roaming\
		path, exists := os.LookupEnv("APPDATA")
		if exists {
			return path
		}
		return ""
	}
	return path
}

// NewSmartAddress создает новый алрес для нового смарта
func NewSmartAddress(now int64) string {
	t, _ := Uint64ToBytes(now)
	str := BytesToSHA256Hex(t)
	addr, err := cryptolib.HexToAddress(str)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return addr
}
