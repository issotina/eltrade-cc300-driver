package main

import (
	"fmt"
	eltrade "github.com/geeckmc/eltrade-cc300-driver/lib"
	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"
	"os"
)

func main() {
	loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr))

	//v := hex.EncodeToString([]byte("	"))

	//fmt.Printf("%s",v)
	exec()

}

func exec() {
	dev, err := eltrade.Open()
	if err != nil {
		panic(fmt.Sprintf("Failed to open dev: %s", err))
	}
	msg := eltrade.NewRequest(eltrade.DEV_STATE)
	resp := dev.Send(msg)

	sq, _ := resp.GetSeq()
	println("seq: ", sq)
	dt, _ := resp.GetData()
	println("-----state \n", dt, "\n ----------------")

	msg = eltrade.NewRequest(eltrade.START_BILL)
	msg.Body("2,Moudilou,3201910768821,0.00,18.00,0.00,18.00,FV,0201810241722,shadai,AIB1")
	resp = dev.Send(msg)

	sq, _ = resp.GetSeq()
	println("seq: ", sq)
	dt, _ = resp.GetData()
	println("------init bill \n", dt, "\n ----------------")

	msg = eltrade.NewRequest(eltrade.ADD_BILL_ITEM)
	msg.Body("Piles plates\tA500*5")
	resp = dev.Send(msg)

	sq, _ = resp.GetSeq()
	println("seq: ", sq)
	dt, _ = resp.GetData()
	println("-----add product \n", dt, "\n ----------------")

	msg = eltrade.NewRequest(eltrade.GET_BILL_TOTAL)
	msg.Body("M2525")
	resp = dev.Send(msg)

	sq, _ = resp.GetSeq()
	println("seq: ", sq)
	dt, _ = resp.GetData()
	println("------total \n", dt, "\n ----------------")

	msg = eltrade.NewRequest(eltrade.END_BILL)
	resp = dev.Send(msg)

	sq, _ = resp.GetSeq()
	println("seq: ", sq)
	dt, _ = resp.GetData()
	println("------fin facture \n", dt, "\n ----------------")
}
