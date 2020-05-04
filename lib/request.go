package eltrade

import (
	"bytes"
	"encoding/hex"
	"errors"
)

type Request struct {
	// Request format : <01><LEN><SEQ><CMD><DATA><05><BCC><03>

	value bytes.Buffer
	Data  []byte
	cmd   Command
	// Seq start from 20h to FFh
	Seq uint8
}

// NewRequest returns instance of Eltrade CC300 command
// data can be null depend on command
//
// MECeF_MCF_SFE_Protocole_v2.pdf
func NewRequest(cmd Command) *Request {
	req := Request{value: bytes.Buffer{}}
	req.value.Write([]byte{SOH})
	req.cmd = cmd
	req.Seq = MIN_SEQ
	return &req
}

func (r *Request) Body(data string) (*Request, error) {
	bytes := []byte(clear(data))
	if len(bytes) > 200 {
		return nil, errors.New("OutOfBound : Data should be less than 200 bytes")
	}
	//TODO: Control escape chars
	r.Data = bytes
	return r, nil
}

func (r *Request) Build() []byte {
	// control if request has more than head bit
	// if yes clear it to allow multiple builds
	if len(r.value.Bytes()) > 1 {
		r.value.Reset()
		r.value.Write([]byte{SOH})
	}
	// Build request using the format
	// Format : <01><LEN><SEQ><CMD><DATA><05><BCC><03>
	payloadLength := LEN_OFFSET + uint8(len(r.Data))
	r.value.Write([]byte{payloadLength, r.Seq, r.cmd.val()})
	r.value.Write(r.Data)
	r.value.Write([]byte{AMB})
	r.value.Write(bcc(r.value.Bytes(), int(payloadLength)))
	r.value.Write([]byte{ETX})
	r.Seq++
	if r.Seq == MAX_SEQ {
		r.Seq = MIN_SEQ
	}
	Logger.Debugf("fn:request.Build -- %s", hex.EncodeToString(r.value.Bytes()))
	return r.value.Bytes()
}
