package acr122u

import (
	"bytes"
)

// Card represents a ACR122U card
type Card interface {
	// Reader returns the name of the reader used
	Reader() string

	// Status returns the card status
	Status() (Status, error)

	// UID returns the UID for the card
	UID() []byte

	Transmit(cmd []byte) ([]byte, error)
	Standard() string
	Name() string
}

type card struct {
	uid    []byte
	reader string
	scard  scardCard
}

func newCard(reader string, sc scardCard) *card {
	return &card{reader: reader, scard: sc}
}

func (c *card) Reader() string {
	return c.reader
}

func (c *card) Status() (Status, error) {
	scs, err := c.scard.Status()
	if err != nil {
		return Status{}, err
	}

	return newStatus(scs)
}

func (c *card) UID() []byte {
	return c.uid
}

// transmit raw command to underlying scardCard
func (c *card) transmit(cmd []byte) ([]byte, error) {
	resp, err := c.scard.Transmit(cmd)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(resp, rcOperationFailed) {
		return nil, ErrOperationFailed
	}

	if bytes.HasSuffix(resp, rcOperationSuccess) {
		return bytes.TrimSuffix(resp, rcOperationSuccess), nil
	}

	return resp, nil
}

// getUID returns the UID for the card
func (c *card) getUID() ([]byte, error) {
	return c.transmit(cmdGetUID)
}

func (c *card) Transmit(cmd []byte) ([]byte, error) {
	return c.transmit(cmd)
}

func (c *card) Standard() string {
	// stat,err := c.Status()
	// if err != nil {
	// 	return ""
	// }
	// st := ISO14443Part3Str
	// tk := stat.Atr.TK()
	// switch tk[3] {
	// case ISO14443Part3:
	// 	st = ISO14443Part3Str
	// 	// TODO

	// }
	return ISO14443Part3Str
}
func (c *card) Name() string {
	stat, err := c.Status()
	if err != nil {
		return ""
	}
	tk := stat.Atr.TK()

	if bytes.Equal(tk[7:8], MIFAREClassic1K) {
		return MIFAREClassic1KStr
	}
	if bytes.Equal(tk[7:8], MIFAREClassic4K) {
		return MIFAREClassic4KStr
	}
	if bytes.Equal(tk[7:8], MIFAREUltralight) {
		return MIFAREUltralightStr
	}
	if bytes.Equal(tk[7:8], MIFAREMini) {
		return MIFAREMiniStr
	}

	return ""
}
