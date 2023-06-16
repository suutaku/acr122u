package acr122u

const (
	ISO14443Part3    = 0x03
	ISO14443Part3Str = "ISO14443 Part 3"
)

var (
	MIFAREClassic1K  = []byte{0x01}
	MIFAREClassic4K  = []byte{0x02}
	MIFAREUltralight = []byte{0x03}
	MIFAREMini       = []byte{0x26}
	TopazAndJewel    = []byte{0xf0, 0x04}
	FeliCa212K       = []byte{0xf0, 0x11}
	FeliCa424K       = []byte{0xf0, 0x12}
	UNDEFINED        = []byte{0xff, 0x00}
)

const (
	MIFAREClassic1KStr  = "MIFARE Classic 1K"
	MIFAREClassic4KStr  = "MIFARE Classic 4K"
	MIFAREUltralightStr = "MIFARE Ultralight"
	MIFAREMiniStr       = "MIFARE Mini"
)
