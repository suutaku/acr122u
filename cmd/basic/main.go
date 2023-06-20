package main

import (
	"net/http"

	"github.com/suutaku/acr122u/pkg/acr122u"
)

var h *NFCHandler

func readHandler(w http.ResponseWriter, req *http.Request) {
	buf := h.Get()
	w.Write(buf)
}

func writeHandler(w http.ResponseWriter, req *http.Request) {
	did := req.URL.Query().Get("did")
	if did == "" {
		w.Write([]byte("parameter did was empty"))
	}
	h.writeBuf <- did
}

func main() {
	ctx, err := acr122u.EstablishContext()
	if err != nil {
		panic(err)
	}
	h = NewNFCHander()
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/write", writeHandler)
	go ctx.Serve(h)
	http.ListenAndServe(":8090", nil)
}
