package acr122u

type CommandsSet interface {
	CmdLoadAuth(key ...byte) []byte
	CmdAuthBlock(b []byte) []byte
	CmdReadBlock(b []byte) []byte
	CmdWriteBlock(data, b []byte) []byte
}

type MifareClassic1kCommands struct {
}

func (mcc MifareClassic1kCommands) CmdLoadAuth(key ...byte) []byte {
	if len(key) == 0 {
		key = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	}
	return append([]byte{0xff, 0x82, 0x00, 0x00, 0x06}, key...)
}

func (mcc MifareClassic1kCommands) CmdAuthBlock(b []byte) []byte {
	tmp := append([]byte{0xff, 0x86, 0x00, 0x00, 0x05, 0x01, 0x00}, b...)
	return append(tmp, []byte{0x60, 0x00}...)
}

func (mcc MifareClassic1kCommands) CmdReadBlock(b []byte) []byte {
	tmp := append([]byte{0xff, 0xb0, 0x00}, b...)
	return append(tmp, 0x10)
}

func (mcc MifareClassic1kCommands) CmdWriteBlock(data, b []byte) []byte {
	tmp := append([]byte{0xff, 0xd6, 0x00}, b...)
	tmp = append(tmp, 0x10)
	return append(tmp, data...)
}
