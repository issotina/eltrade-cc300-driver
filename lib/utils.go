package eltrade

import (
	"encoding/hex"
	"fmt"
	"strings"
)

var (
	charsToEscape = map[byte]string{
		'\r': "^xa;",
		'\n': "^xd;",
		',':  "^x2c;",
		'<':  "^lt;",
		'>':  "^gt;",
		'&':  "^amp;",
	}
	replacer = strings.NewReplacer("\r", "^xa;",
		"\n", "^xd;",
		//",", "^x2c;",
		"<", "^lt;",
		">", "^gt;",
		"&", "^amp;")
)

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

func clear(str string) string {
	return replacer.Replace(str)
}

func bcc(data []byte, limit int) []byte {
	sum := 0x00
	bcc := make([]byte, 4)

	for i, b := range data {
		if i == 0 {
			continue
		}
		if i == limit {
			break
		}
		sum += int(b)
	}
	bcc[0] = byte(((sum & 0xf000) >> 12) + 0x30)
	bcc[1] = byte(((sum & 0xf00) >> 8) + 0x30)
	bcc[2] = byte(((sum & 0xf0) >> 4) + 0x30)
	bcc[3] = byte(((sum & 0xf) >> 0) + 0x30)
	//TODO : add loop instead of step by step ascii conversion
	return bcc
}

type EltradeString struct {
	Val string
}

func (e *EltradeString) Append(str string) *EltradeString {
	if strings.TrimSpace(str) != "" {
		e.AppendWD(str, ",")
	}
	return e
}
func (e *EltradeString) AppendWD(str string, delimiter string) *EltradeString {
	e.Val += fmt.Sprintf("%s%s", delimiter, strings.TrimSpace(str))
	return e
}
