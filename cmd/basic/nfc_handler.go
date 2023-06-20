package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suutaku/acr122u/pkg/acr122u"
)

type NFCHandler struct {
	acr122u.Logger
	readBuf  chan []byte
	writeBuf chan string
}

func NewNFCHander() *NFCHandler {
	return &NFCHandler{
		Logger:   log.New(os.Stdout, "", 0),
		readBuf:  make(chan []byte, 1),
		writeBuf: make(chan string, 0),
	}
}

const block = 0x04

func (h *NFCHandler) ServeCard(c acr122u.Card) {
	var cmdSet acr122u.CommandsSet
	fmt.Println(c.Name())
	if c.Name() == acr122u.MIFAREClassic1KStr {
		cmdSet = &acr122u.MifareClassic1kCommands{}
	}
	if c.Name() == acr122u.MIFAREClassic4KStr {
		cmdSet = &acr122u.MifareClassic1kCommands{}
	}
	if c.Name() == acr122u.MIFAREUltralightStr {
		cmdSet = &acr122u.MifareClassic1kCommands{}
	}

	select {
	case wm := <-h.writeBuf:
		fmt.Printf("write buf comes[%s]\n", wm)
		tag := acr122u.NewTag([]byte{0x55, 0x04}, wm)
		err := acr122u.PutTag(c, cmdSet, tag)
		if err != nil {
			fmt.Println(err)
			return
		}

	default:
		tags, err := acr122u.GetTags(c, cmdSet)
		if err != nil {
			fmt.Println(err)
			return
		}
		for i := 0; i < len(tags); i++ {
			h.readBuf <- tags[i].Message.Payload
		}
	}
}

func (h *NFCHandler) Get() []byte {
	return <-h.readBuf
}

func (h *NFCHandler) Put(msg string) {
	h.writeBuf <- msg
}
