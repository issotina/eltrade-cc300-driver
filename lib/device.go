package eltrade

import (
	"bytes"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"github.com/juju/loggo"
	"go.bug.st/serial.v1/enumerator"
	"io"
	"time"
)

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

func (c Command) val() uint8 {
	return uint8(c)
}

const (
	SOH                         = 0x1
	AMB                         = 0x5
	ETX                         = 0x3
	BRK                         = 0x04
	LEN_OFFSET                  = 0x24
	MIN_SEQ                     = 0x20
	MAX_SEQ                     = 0xFF
	DEV_STATE           Command = 0xC1
	NETWORK_STATE       Command = 0xC2
	TAXPAYER_INFO               = 0x2B
	START_BILL                  = 0xC0
	ADD_BILL_ITEM               = 0x31
	GET_BILL_SUB_TOTAL          = 0x33
	GET_BILL_TOTAL              = 0x35
	END_BILL                    = 0x38
	CMD_PROCESSING_TIME         = 100 * time.Millisecond
	RESPONSE_DELIMITER          = 0x2c
	EltradeVID                  = "03EB"
	EltradePID                  = "6119"
)

var (
	Logger = loggo.GetLogger("eltrade.driver")
)

type Device struct {
	serial io.ReadWriteCloser
	IsOpen bool
}

func Open() (*Device, error) {
	Logger.SetLogLevel(loggo.DEBUG)
	dev := Device{IsOpen: false}
	portName, err := getPortName()
	if err != nil {
		Logger.Errorf("fn:eltrade.Open -- %s", err.Error())
		return nil, fmt.Errorf("serial.Open: %v", err)
	}
	options := serial.OpenOptions{
		PortName:              portName,
		BaudRate:              115200, //https://github.com/ethno2405/eltrade-fiscal-device-protocol/blob/master/EltradeProtocol/EltradeProtocol/EltradeFiscalDeviceDriver.cs
		DataBits:              8,      // MECeF_MCF_SFE_Protocole_v2.pdf : Page 5
		StopBits:              1,
		ParityMode:            serial.PARITY_NONE,
		InterCharacterTimeout: 500,
	}
	dev.serial, err = serial.Open(options)
	if err != nil {
		Logger.Errorf("fn:eltrade.Open -- %s", err.Error())
		return nil, fmt.Errorf("serial.Open: %v", err)
	}
	dev.IsOpen = true
	Logger.Debugf("fn:eltrade.Open -- Success ")
	return &dev, nil
}

func getPortName() (string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", err
	}
	if len(ports) == 0 {
		return "", fmt.Errorf("No serial ports found!")
	}
	for _, port := range ports {
		Logger.Debugf("Found port: Name: %s\n VendorID: %s\n ProductId: %s\n", port.Name, port.VID, port.PID)
		if port.IsUSB && port.VID == EltradeVID && port.PID == EltradePID {
			Logger.Infof("Matched port: %s\n", port.Name)
			return port.Name, nil
		}
	}
	return "", fmt.Errorf("No Eltrade devices found on serial ports !")
}

func (dev *Device) Send(req *Request) Response {
	if !dev.IsOpen {
		return Response{status: NOT_READY}
	}
	dev.serial.Write(req.Build())
	rawResponse := bytes.Buffer{}
	r := Response{}
	askResponseCount := 0
	for {
		// Slave -} Host. Slave notifies Host that the command will consume for execution time more than 100 ms.
		// * The Slave have to send {SYN} on every 100ms while processing command and return response with packaged message
		n, err := rawResponse.ReadFrom(dev.serial)
		if err != nil {
			Logger.Errorf("fn:eltrade.Send -- %s", err.Error())
			return Response{status: INVALID}
		}
		Logger.Debugf("fn:eltrade.Send -- %s Bytes read ", n)

		r.Parse(rawResponse.Bytes())
		seq, err := r.GetSeq()
		if err != nil {
			Logger.Errorf("fn:eltrade.Send -- %s", err.Error())
			return Response{status: INVALID}
		}
		if seq != uint8(SYN) {
			break
		}
		askResponseCount++
		Logger.Debugf("fn:eltrade.Send -- wait for response: %s time ", askResponseCount)
		//time.Sleep(CMD_PROCESSING_TIME)
	}

	return r
}

func (dev *Device) Close() {
	if dev.IsOpen {
		dev.serial.Close()
	}
}
