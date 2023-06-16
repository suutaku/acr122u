package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-vgo/robotgo"
	"github.com/suutaku/acr122u/pkg/acr122u"
)

func main() {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}

	h := &handler{log.New(os.Stdout, "", 0)}

	ctx.Serve(h)
}

type handler struct {
	acr122u.Logger
}

const block = 0x04

func (h *handler) ServeCard(c acr122u.Card) {
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
	tags, err := acr122u.GetTags(c, cmdSet)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(tags); i++ {
		fmt.Printf("%x:%#s\n", tags[i].Message.Type, tags[i].Message.Payload)
		robotgo.TypeStr(string(tags[i].Message.Payload))
	}
}
