package acr122u

import "fmt"

const (
	NDEF    = 0x03
	ENDFILE = 0xfe
)

const DEFAULT_BLOCK_NUM uint8 = 0x04

type Tag struct {
	Message *Message
	Length  uint8
}

type Message struct {
	TNF           uint8
	TypeLength    uint8
	PayloadLength uint8
	Type          []byte
	Payload       []byte
}

func NewMessage(buf []byte) *Message {
	if len(buf) < 3 {
		fmt.Printf("invalid message buf with length: %v\n", len(buf))
		return nil
	}
	ret := &Message{}
	ret.TNF = buf[0]
	ret.TypeLength = buf[1]
	ret.PayloadLength = buf[2]
	exceptLen := uint8(3) // 1 (tnf) + 1 (type len) + 1 (payload len)
	exceptLen += ret.TypeLength + ret.PayloadLength
	if len(buf) < int(exceptLen) {
		fmt.Printf("unexcepted message buf with length: %v, excepted: %v\n", len(buf), exceptLen)
		return nil
	}
	ret.Type = buf[3 : 4+ret.TypeLength]
	ret.Payload = buf[4+ret.TypeLength:]
	return ret
}

func NewTag(buf []byte) *Tag {
	if len(buf) < 3 {
		fmt.Printf("invalid tag buf with length: %v\n", len(buf))
		return nil
	}
	ret := &Tag{}
	// get lenth
	ret.Length = buf[1]
	exceptLen := uint8(3 + ret.Length)
	if len(buf) < int(exceptLen) {
		fmt.Printf("unexcepted tag buf with length: %v, excepted: %v\n", len(buf), exceptLen)
		return nil
	}

	msg := buf[2 : exceptLen-1]
	ret.Message = NewMessage(msg)
	return ret

}

func GetTags(c Card, cmds CommandsSet, key ...byte) ([]*Tag, error) {
	_, err := c.Transmit(cmds.CmdLoadAuth(key...))
	if err != nil {
		return nil, fmt.Errorf("load auth: %v", err)
	}

	block := DEFAULT_BLOCK_NUM
	//read all selectors
	memory := make([]byte, 0)
	for {
		_, err = c.Transmit(cmds.CmdAuthBlock([]byte{block}))
		if err != nil {
			break
		}
		resp, err := c.Transmit(cmds.CmdReadBlock([]byte{block}))
		if err != nil {
			break
		}
		if len(resp) == 0 {
			fmt.Printf("empty data\n")
		}
		memory = append(memory, resp...)
		block++
	}
	// split memory to Tag Datas (according NDEF)
	startP := -1
	ret := make([]*Tag, 0)
	for i := 0; i < len(memory); i++ {
		if startP == -1 && memory[i] == NDEF {
			startP = i

			continue
		}
		if startP != -1 && memory[i] == ENDFILE {

			buf := memory[startP : i+1]
			startP = -1
			t := NewTag(buf)
			if t != nil {
				ret = append(ret, t)
			}
		}
	}
	return ret, nil

}
