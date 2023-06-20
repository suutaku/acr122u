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
	Id            []byte
	Payload       []byte
}

func NewMessage(t, msg []byte) *Message {
	ret := &Message{}
	if t == nil || msg == nil {
		return ret
	}
	ret.TNF = 0xd1
	ret.TypeLength = uint8(len(t) - 1)
	ret.PayloadLength = uint8(len(msg))
	ret.Type = t[:ret.TypeLength]
	ret.Id = t[ret.TypeLength:]
	ret.Payload = msg
	return ret
}

func (msg *Message) Unmarshal(buf []byte) error {
	if len(buf) < 3 {
		return fmt.Errorf("invalid message buf with length: %v\n", len(buf))
	}

	msg.TNF = buf[0]
	msg.TypeLength = buf[1]
	msg.PayloadLength = buf[2]
	msg.Type = buf[3 : 3+msg.TypeLength]
	msg.Id = buf[3+msg.TypeLength : 4+msg.TypeLength]
	msg.Payload = buf[4+msg.TypeLength:]
	return nil
}

func (msg *Message) Marshal() ([]byte, error) {
	expLen := uint8(3)
	expLen += msg.TypeLength + msg.PayloadLength + 1 // type len + id len + payload len

	buf := make([]byte, 3)
	buf[0] = msg.TNF
	buf[1] = msg.TypeLength
	buf[2] = msg.PayloadLength
	buf = append(buf, msg.Type...)
	buf = append(buf, msg.Id...)
	buf = append(buf, msg.Payload...)
	return buf, nil

}

func NewTag(t []byte, m string) *Tag {
	msg := NewMessage(t, []byte(m))
	b, err := msg.Marshal()
	if err != nil {
		return nil
	}

	return &Tag{Message: msg, Length: uint8(len(b))}
}

func (tag *Tag) Unmarshal(buf []byte) error {
	if len(buf) < 3 {
		return fmt.Errorf("invalid tag buf with length: %v\n", len(buf))
	}

	// get lenth
	tag.Length = buf[1]
	msg := buf[2 : len(buf)-1]
	tag.Message = NewMessage(nil, nil)
	return tag.Message.Unmarshal(msg)
}

func (tag *Tag) Marshal() ([]byte, error) {

	buf := make([]byte, 1)
	buf[0] = NDEF
	buf = append(buf, tag.Length)
	msg, err := tag.Message.Marshal()
	if err != nil {
		return nil, fmt.Errorf("tag unmarshal message: %v", err)
	}
	buf = append(buf, msg...)
	buf = append(buf, ENDFILE)

	return buf, nil
}

func PutTag(c Card, cmds CommandsSet, tag *Tag, key ...byte) error {
	_, err := c.Transmit(cmds.CmdLoadAuth(key...))
	if err != nil {
		return fmt.Errorf("load auth: %v", err)
	}

	block := DEFAULT_BLOCK_NUM
	//append start and stop byte
	memory, err := tag.Marshal()
	if err != nil {
		return err
	}
	memory = append([]byte{0x00, 0x00}, memory...)
	for len(memory)%16 != 0 {
		memory = append(memory, []byte{0x00}...)
	}
	fmt.Printf("put memory: [%x]\n", memory)
	// spit to blocks
	for i := 0; i < len(memory); i += 16 {

		_, err = c.Transmit(cmds.CmdAuthBlock([]byte{block}))
		if err != nil {
			return err
		}
		_, err = c.Transmit(cmds.CmdWriteBlock(memory[i:i+16], []byte{block}))
		if err != nil {
			return fmt.Errorf("write at %d: %v", 0x04, err)
		}
		block++
	}

	return err
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
			t := NewTag(nil, "")
			fmt.Printf("get %x\n", buf)
			if err := t.Unmarshal(buf); err == nil {
				ret = append(ret, t)
			} else {
				fmt.Printf("%v\n", err)
			}
		}
	}
	return ret, nil

}
