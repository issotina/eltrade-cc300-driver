package eltrade

import "encoding/hex"

func GetBytes(hexString string) []byte {
	value, _ := hex.DecodeString(hexString)
	return value
}

type ByteArray struct {
	bytes []byte
}

func NewByteArray() *ByteArray {
	return &ByteArray{[]byte{}}
}

func (byteArray *ByteArray) Append(b []byte) *ByteArray {
	byteArray.bytes = append(byteArray.bytes, b...)
	return byteArray
}

func (byteArray *ByteArray) AppendHex(hex string) *ByteArray {
	byteArray.bytes = append(byteArray.bytes, GetBytes(hex)...)
	return byteArray
}

func (byteArray *ByteArray) Build() []byte {
	return byteArray.bytes
}
