package eltrade

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

type Status int

const (
	// NAK is like HTTP 400 badRequest,
	// This code is sent by the Slave, if it detects an error in the checksum or
	// the format of the received message.When the Host receives a NAK, it needs to resend a message with the same sequential number
	NAK Status = 0x15
	// SYN code is sent by the Slave, when it receives a command requiring longer execution time.
	// SYN is sentevery 100 ms, until the packed response message is ready
	SYN Status = 0x16
	// OK represent successful response
	OK        Status = 0x01
	NOT_READY Status = -1
	INVALID   Status = 0x0
	cmdIndex         = 4
	seqIndex         = 3
)

type Response struct {
	// Response format : <01><LEN><SEQ><CMD><DATA><04><STATUS><05><BCC><03>
	raw         bytes.Buffer
	readOnlyRaw bytes.Buffer
	status      Status
	seq         uint8
	cmd         uint8
}

func (r *Response) Parse(rawBytes []byte) Status {
	Logger.Debugf("fn:Response.Parse -- Received bytes %s ", hex.EncodeToString(rawBytes))
	rawBytes = clean(rawBytes)
	Logger.Debugf("fn:Response.Parse -- Cleared bytes %s ", hex.EncodeToString(rawBytes))
	r.raw.Write(rawBytes)
	r.readOnlyRaw.Write(rawBytes)
	header, err := r.raw.ReadByte()
	if err != nil {
		Logger.Errorf("fn:Response.Parse -- %s", err.Error())
		return INVALID
	}
	r.status = Status(header)
	return Status(header)
}

func (r *Response) GetSeq() (uint8, error) {
	if r.seq != 0 {
		return r.seq, nil
	}
	if r.status != OK || r.readOnlyRaw.Len() < seqIndex {
		err := errors.New(fmt.Sprintf("Bad Response. Status : %s", hex.EncodeToString([]byte{byte(r.status)})))
		Logger.Errorf("fn:Response.GetSeq -- %s", err.Error())
		return 0, err
	}
	r.seq = r.readOnlyRaw.Next(seqIndex)[seqIndex-1]
	return r.seq, nil
}

func (r *Response) GetCmd() (uint8, error) {
	if r.cmd != 0 {
		return r.cmd, nil
	}
	if r.status != OK || r.readOnlyRaw.Len() < cmdIndex {
		err := errors.New(fmt.Sprintf("Bad Response. Status : %s", r.status))
		Logger.Errorf("fn:Response.GetCmd -- %s", err.Error())
		return 0, err
	}
	r.cmd = r.readOnlyRaw.Next(cmdIndex)[cmdIndex-1]
	return r.seq, nil
}

func (r *Response) GetData() (string, error) {
	if r.status != OK {
		return "", errors.New(fmt.Sprintf("Bad Response. Status : %s", r.status))
	}
	// shift LEN SEQ and CMD
	r.raw.ReadByte()
	r.raw.ReadByte()
	r.raw.ReadByte()
	next, _ := r.raw.ReadByte()
	data := bytes.Buffer{}
	for next != BRK {

		data.Write([]byte{next})
		next, _ = r.raw.ReadByte()
	}
	return string(data.Bytes()), nil
}

func clean(rawBytes []byte) []byte {
	for ok := false; ok; ok = len(rawBytes) > 1 && rawBytes[len(rawBytes)-1] != ETX {
		rawBytes = rawBytes[:len(rawBytes)-1]
	}
	for ok := len(rawBytes) > 1 && rawBytes[0] != SOH; ok; ok = len(rawBytes) > 1 && rawBytes[0] != SOH {
		rawBytes = rawBytes[1:]
	}
	return rawBytes
}
