package datachain

// IBlockEncoder encode block to any formats.
type IBlockEncoder interface {
	ToBencode() []byte
	ToJSON() []byte
}

// IBlock block interface.
type IBlock interface {
	IBlockEncoder
	GetHeight() int64
	GetHash() string
	GetPrevHash() string
	Check() bool
}

// ChainState load and setting chain state from leveldb.
type ChainState struct {
	ChainID    int64
	LastHeight int64
	LastHash   string
}

const lastHash = "1"
