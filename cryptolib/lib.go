// Copyright (c) 2020, Ruslan Dorofeev <rusldv@yandex.ru>
// All rights reserved.
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cryptolib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

const (
	generateSize = 40
)

// ESig сигнатура подписи (для кодирования в ASN1).
type ESig struct {
	R, S *big.Int
}

// Generate генерация случайного числа для приватного ключа.
func Generate() []byte {
	buf := make([]byte, generateSize)
	rand.Reader.Read(buf)
	return buf
}

// GenerateHex генерация случайного числа для приватного ключа в строковом представлении 16-ричного числа.
func GenerateHex() string {
	return hex.EncodeToString(Generate())
}

// HexToPK создание объекта приватного ключа из строкового 16-ричного представления.
func HexToPK(hexStr string) (*ecdsa.PrivateKey, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	d := new(big.Int)
	d.SetBytes(bytes)
	PK := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	PK.PublicKey.Curve = curve
	PK.PublicKey.X, PK.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())
	PK.D = d

	return PK, nil
}

// XYToPublicKey создает объект публичного ключа из X, Y координат.
func XYToPublicKey(x, y *big.Int) *ecdsa.PublicKey {
	pub := new(ecdsa.PublicKey)
	pub.Curve = elliptic.P256()
	pub.X = x
	pub.Y = y

	return pub
}

// PublicKeyMarshal кодирует публичный ключ в байты.
func PublicKeyMarshal(pub *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(pub.Curve, pub.X, pub.Y)
}

// PublicKeyMarshalHex кодирует публичный ключ в строковый хеш.
func PublicKeyMarshalHex(pub *ecdsa.PublicKey) string {
	return hex.EncodeToString(elliptic.Marshal(pub.Curve, pub.X, pub.Y))
}

// PublicKeyUnmarshal декодирует публичный ключ из байтов в X, Y координаты.
func PublicKeyUnmarshal(bytes []byte) (*big.Int, *big.Int) {
	return elliptic.Unmarshal(elliptic.P256(), bytes)
}

// PublicKeyUnmarshalHex декодирует публичный ключ из строкового хеша в X, Y координаты.
func PublicKeyUnmarshalHex(hexStr string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, nil, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), bytes)

	return x, y, nil
}

// HexToAddress преобразует хеш публичного ключа в адрес.
func HexToAddress(hexStr string) (string, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	sha2 := sha256.New()
	sha2.Write(bytes)
	ripe := ripemd160.New()
	ripe.Write(sha2.Sum(nil))
	bhash := ripe.Sum(nil)
	hash := hex.EncodeToString(bhash[:])
	return hash, nil
}

// SignData создает R, S сигнатуру для данных.
func SignData(PK *ecdsa.PrivateKey, data []byte) (*big.Int, *big.Int, error) {
	return ecdsa.Sign(rand.Reader, PK, data)
}

// SignDataASN1 создает R, S сигнатуру для данных в ASN1.
func SignDataASN1(PK *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	r, s, err := SignData(PK, data)
	if err != nil {
		return nil, err
	}
	esig := ESig{
		R: r,
		S: s,
	}
	sig, err := asn1.Marshal(esig)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// SignDataASN1Hex создает R, S сигнатуру для данных в ASN1 и 16-ричном строковом формате.
func SignDataASN1Hex(PK *ecdsa.PrivateKey, data []byte) (string, error) {
	sig, err := SignDataASN1(PK, data)
	if err != nil {
		return "", err
	}
	hexSig := hex.EncodeToString(sig)
	return hexSig, nil
}

// VerifyData верифицирует R, S сигнатуру и данные.
func VerifyData(pub *ecdsa.PublicKey, R, S *big.Int, data []byte) bool {
	//fmt.Println(pub)
	//fmt.Println(R, S)
	//fmt.Println(data)

	return ecdsa.Verify(pub, data, R, S)
}

// VerifyDataASN1Hex верифицирует R, S сигнатуру в 16-ричном строковом формате и данные.
func VerifyDataASN1Hex(pub *ecdsa.PublicKey, sig string, data []byte) bool {
	// Unmarshal sig from ASN1
	bytes, err := hex.DecodeString(sig)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var esig ESig
	asn1.Unmarshal(bytes, &esig)

	return VerifyData(pub, esig.R, esig.S, data)
}
