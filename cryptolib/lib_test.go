package cryptolib

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	fmt.Println(Generate())
}

func TestGenerateHex(t *testing.T) {
	fmt.Println(GenerateHex())
}

func TestHexToPK(t *testing.T) {
	fmt.Println(HexToPK(GenerateHex()))
	fmt.Println(HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a"))
}

func TestXYToPublicKey(t *testing.T) {
	PK, err := HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y))
}

func TestPublicKeyMarshal(t *testing.T) {
	PK, err := HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a")
	if err != nil {
		t.Error(err)
	}
	pub := XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y)
	fmt.Println(PublicKeyMarshal(pub))
}

func TestPublicKeyMarshalHex(t *testing.T) {
	PK, err := HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a")
	if err != nil {
		t.Error(err)
	}
	pub := XYToPublicKey(PK.PublicKey.X, PK.PublicKey.Y)
	fmt.Println(PublicKeyMarshalHex(pub))
}

func TestPublicKeyUnmarshalHex(t *testing.T) {
	x, y, err := PublicKeyUnmarshalHex("048405a2ca7a137b2dfcd653f9f0c524051e888e361c2602163ce0518cb818c59be758d30709d70efe0474d8db059eba0455c606c98bd048e544d41af93f8e237c")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(x, y)
}

func TestHexToAddress(t *testing.T) {
	address, err := HexToAddress("048405a2ca7a137b2dfcd653f9f0c524051e888e361c2602163ce0518cb818c59be758d30709d70efe0474d8db059eba0455c606c98bd048e544d41af93f8e237c")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(address)
}

func TestSignDataASN1Hex(t *testing.T) {
	PK, _ := HexToPK("71e2d237d1215f76d6c42d6a835e26ba8719da89cb1936b55c62fa1471f6b0c900d4f6823013c90a")
	data := []byte("test")
	sig, err := SignDataASN1Hex(PK, data)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sig)
}

func TestVerifyDataASN1Hex(t *testing.T) {
	x, y, err := PublicKeyUnmarshalHex("048405a2ca7a137b2dfcd653f9f0c524051e888e361c2602163ce0518cb818c59be758d30709d70efe0474d8db059eba0455c606c98bd048e544d41af93f8e237c")
	if err != nil {
		t.Error(err)
	}
	pub := XYToPublicKey(x, y)
	sig := "304502205df188e4e902c1bc216ff777a82c1be3b5e28face639d1197ddaefd5ed57db51022100c33516cc50254bf75bc1d2e09a008024e486d34b214543292e9a760d2a9accbd"
	data := []byte("test")

	fmt.Println(VerifyDataASN1Hex(pub, sig, data))
}
