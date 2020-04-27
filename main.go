package main

import (
	"bytes"
	"eltrade/lib"
	"encoding/hex"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"io"
	"log"
	"time"
)

func main() {

	options := serial.OpenOptions{
		PortName:              "/dev/tty.usbmodem143101", // ls /dev/tty.*
		BaudRate:              115200,                    //https://github.com/ethno2405/eltrade-fiscal-device-protocol/blob/master/EltradeProtocol/EltradeProtocol/EltradeFiscalDeviceDriver.cs
		DataBits:              8,                         // MECeF_MCF_SFE_Protocole_v2.pdf : Page 5
		StopBits:              1,
		ParityMode:            serial.PARITY_NONE,
		InterCharacterTimeout: 500,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	/*msg := bytes.Buffer{}
	msg.Write([]byte{0x01})//soh
	msg.Write([]byte{0x20 + 4 }) //len
	msg.Write([]byte{0x21}) //seq
	msg.Write([]byte{0xC1}) //cmd
	//msg.Write([]byte{}) //data
	msg.Write([]byte{0x05}) //AMB
	msg.Write(lib.CheckSum(0x24  + 0x21 + 0xC1 + 0x05 )) //bcc
	msg.Write([]byte{0x03}) //etx
	fmt.Println("Rq: \n",hex.EncodeToString(msg.Bytes()))*/
	msg := eltrade.NewRequest(eltrade.DEV_STATE)
	b := msg.Build()
	b = msg.Build()
	b = msg.Build()
	fmt.Println("Rq: \n", hex.EncodeToString(b))
	port.Write(b)
	time.Sleep(500 * time.Millisecond)

	//49 50

	//s := "12AB"

	//moved := lib.CheckSum(0x12AB)
	//for i := 0; i < len(moved); i++ {
	//	fmt.Print(" \t",IntToHexChars(49))

	//}
	/*	a := 0x12AB + 0x00
		b := byte(((a & 0xf000) >> 12) + 0x30)
		c := byte(((a & 0xf00) >> 8) + 0x30)
		d := byte(((a & 0xf0) >> 4) + 0x30)
		e := byte(((a & 0xf) >> 0) + 0x30)



		fmt.Println("\n Rq: \n",lib.CheckSum(0x12AB ))
		fmt.Print("\t",b)
		fmt.Print("\t",c)
		fmt.Print("\t",d)
		fmt.Print("\t",e)
	*/
	// Reading the output

	// Make sure to close it later.
	defer port.Close()

	i := 0
	//	s := bytes.Buffe
	cr := byte(0x00)
	r := bytes.Buffer{}
	for ok := true; ok; ok = cr != byte(0x3) {
		buf := make([]byte, 1)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error reading from serial port: ", err)
			}
		} else {
			buf = buf[:n]
			//			if ()
			r.Write(buf)
			i++
			cr = buf[0]
			fmt.Println("Rx: ", i, "", n, "", hex.EncodeToString(buf))
		}
	}

	resp := eltrade.Response{}
	status := resp.Parse(r.Bytes())
	sq, _ := resp.GetSeq()
	println("seq: ", sq, "status", uint16(status))
	dt, _ := resp.GetData()
	println("------my data \n", dt, "\n ----------------")
	//<01><LEN><SEQ><CMD><DATA><04><STATUS><05><BCC><03>
	len := uint8(r.Bytes()[1] - byte(0x26))
	println("len", len)
	d := r.Bytes()[4 : len-7]
	println("data", d)
	bs, err := hex.DecodeString(hex.EncodeToString(d))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}

func getResp(port io.ReadWriteCloser) {

}

func dec2hex(dec int) string {
	color := dec * 255 / 100
	return fmt.Sprintf("%02x", color)
}

func IntToHexChars(v uint64) []byte {
	HEX := []byte("0123456789ABCDEF")
	result := make([]byte, 0, 16)
	for v != 0 {
		nibble := v & 0x0F
		c := HEX[nibble]
		// not optimal code here?
		result = append(result, c)
		v = (v >> 4)
	}
	return result
}
