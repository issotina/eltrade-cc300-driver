package eltrade

import "strings"

/**
 Device MESSAGE SEQUENCE
  <01>
	Preamble
	length: 1byte
	value: 01H

  <LEN>
	Number of bytes from <01> (not included) to <05> (included)
	plus a fixed offset of 20H
	length: 1
	bytevalue: 20H - FFH

  <SEQ>
	Sequential number of the frame
	length: 1 byte
	value: 20H - 7FH
	The fiscal printer enters the same <SEQ> in the response message.
	If the FP receives amessage with the same <SEQ> as the last received message,
	then it takes no action, just repeatsthe last sent message.

  <CMD>
	Command code
	length: 1 byte
	value: 20H - FFH
	The FP enters the same <CMD> in the response message.
	If the printer receives a non-existing code,
	it responds with a packed message with zero length of the data field and sets the
	respective status bit.

  <DATA>
	Data
	length:
		- 0–218 bytes for Host to Printer.
		- 0–213 bytes for Printer to Host.
	value:20H–FFH and, additionally, 09H and 0AH.
	The data area format and length are command-dependent.
	If the command has no data, then the length of this field is zero.
	If there is a syntax error in the data,
	then the respective status bit is set and a packed message
	with zero length of the data field is returned.

  <04>
	Separator (for Printer-to-Host messages only)
	length: 1 byte
	value: 04H
*/
type Command uint8

const (
	SOH                        = 0x1
	AMB                        = 0x5
	ETX                        = 0x3
	BRK                        = 0x04
	LEN_OFFSET                 = 0x24
	MIN_SEQ                    = 0x20
	MAX_SEQ                    = 0xFF
	DEV_STATE          Command = 0xC1
	NETWORK_STATE      Command = 0xC2
	CITIZEN_INFO               = 0x2B
	START_BILL                 = 0xC0
	ADD_BILL_ITEM              = 0x31
	GET_BILL_SUB_TOTAL         = 0x33
	GET_BILL_TOTAL             = 0x35
	END_BILL                   = 0x38
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
		",", "^x2c;",
		"<", "^lt;",
		">", "^gt;",
		"&", "^amp;")
)

//, '\n', ',', '<', '>', '&'
func (c Command) val() uint8 {
	return uint8(c)
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

func clear(str string) string {
	return replacer.Replace(str)
}
