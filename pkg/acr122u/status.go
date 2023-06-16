package acr122u

import "github.com/ebfe/scard"

type Atr struct {
	raw []byte
}

func (atr Atr) T0() uint8 {
	return atr.raw[1]
}

func (atr Atr) TD1() uint8 {
	return atr.raw[2]
}

func (atr Atr) TD2() uint8 {
	return atr.raw[3]
}

func (atr Atr) T1() uint8 {
	return atr.raw[4]
}

func (atr Atr) TK() []byte {
	if atr.T1() == 0x80 { // part 3
		len := atr.raw[6]
		return atr.raw[5 : 2+len]
	}
	// part 4
	return atr.raw[5:9]
}

// Status contains the status of a card
type Status struct {
	Reader         string
	State          uint32
	ActiveProtocol uint32
	Atr            Atr
}

func newStatus(scs *scard.CardStatus) (Status, error) {
	if scs == nil {
		return Status{}, scard.ErrUnknownCard
	}

	return Status{
		Reader:         scs.Reader,
		State:          uint32(scs.State),
		ActiveProtocol: uint32(scs.ActiveProtocol),
		Atr:            Atr{scs.Atr},
	}, nil
}
