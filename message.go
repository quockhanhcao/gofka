package main

import (
	"bufio"
	"fmt"
)

const (
	ECHO  = 1 // echo
	P_REG = 2 // producer register
	PCM   = 3 // producer to consumer msg
	// other message type
	// response (simply add 100 to the initial code)
	R_ECHO  = 101 // echo response
	R_P_REG = 102 // producer register response
	R_PCM   = 103
)

type Message struct {
	ECHO  *string
	P_REG *ProducerRegisterMessage
	PCM   []byte
	// response
	R_ECHO  *string
	R_P_REG *byte
	R_PCM   *byte
}

type ProducerRegisterMessage struct {
	port    uint16
	topicID uint16
}

func (m *ProducerRegisterMessage) fromByte(data []byte) {
	// first 2 bytes: port
	// next 2 bytes: topicID
	m.port = uint16(data[0])<<8 + uint16(data[1])
	m.topicID = uint16(data[2])<<8 + uint16(data[3])
}

func (m *ProducerRegisterMessage) toByte() []byte {
	var data [4]byte
	// first 2 bytes: port
	// next 2 bytes: topicID
	data[0] = byte(m.port >> 8)
	data[1] = byte(m.port % 255)

	data[2] = byte(m.topicID >> 8)
	data[3] = byte(m.topicID % 255)

	return data[:4]
}

// message format
// stream[0]: message size
// stream[1:]: the message
func readFromStream(streamReadWrite *bufio.ReadWriter) ([]byte, error) {
	length, err := streamReadWrite.ReadByte()
	if err != nil {
		return nil, err
	}
	data, err := streamReadWrite.Peek(int(length))
	if err != nil {
		return nil, err
	}
	_, err = streamReadWrite.Discard(int(length))
	if err != nil {
		return nil, err
	}

	return data, nil
}

func readMessageFromStream(streamReadWrite *bufio.ReadWriter) (*Message, error) {
	data, err := readFromStream(streamReadWrite)
	if err != nil {
		return nil, err
	}
	msg := parseMessage(data)
	return msg, nil
}

func parseMessage(data []byte) *Message {
	switch data[0] {
	case ECHO:
		st := string(data[1:])
		return &Message{ECHO: &st}
	case R_ECHO:
		st := string(data[1:])
		return &Message{R_ECHO: &st}
	case P_REG:
		st := data[1:]
		parsedMsg := ProducerRegisterMessage{}
		parsedMsg.fromByte(st)
		return &Message{P_REG: &parsedMsg}
	case R_P_REG:
		st := data[1]
		return &Message{R_P_REG: &st}
	case PCM:
		st := data[1:]
		return &Message{PCM: st}
	case R_PCM:
		st := data[1]
		return &Message{R_PCM: &st}
	default:
		return nil
	}
}

func writeDataToStreamWithType(streamReadWrite *bufio.ReadWriter, data string, msgType byte) error {
	// write length
	err := streamReadWrite.WriteByte(byte(len(data) + 1))
	if err != nil {
		return err
	}
	// write type
	err = streamReadWrite.WriteByte(msgType)
	if err != nil {
		return err
	}
	// write msg
	_, err = streamReadWrite.WriteString(data)
	if err != nil {
		return err
	}
	err = streamReadWrite.Flush()
	if err != nil {
		return err
	}

	return nil
}

func writeMessageToStream(streamReadWrite *bufio.ReadWriter, msg Message) error {
	if msg.ECHO != nil {
		if err := writeDataToStreamWithType(streamReadWrite, *msg.ECHO, ECHO); err != nil {
			return err
		}
	} else if msg.R_ECHO != nil {
		if err := writeDataToStreamWithType(streamReadWrite, *msg.R_ECHO, R_ECHO); err != nil {
			return err
		}
	} else if msg.P_REG != nil {
		data := string(msg.P_REG.toByte())
		if err := writeDataToStreamWithType(streamReadWrite, data, P_REG); err != nil {
			return err
		}
	} else if msg.R_P_REG != nil {
		data := fmt.Sprintf("%d", *msg.R_P_REG)
		if err := writeDataToStreamWithType(streamReadWrite, data, R_P_REG); err != nil {
			return err
		}
	} else if msg.PCM != nil {
		if err := writeDataToStreamWithType(streamReadWrite, string(msg.PCM), PCM); err != nil {
			return err
		}
	} else if msg.R_PCM != nil {
		if err := writeDataToStreamWithType(streamReadWrite, string(*msg.R_PCM), R_PCM); err != nil {
			return err
		}
	}
	return nil
}
