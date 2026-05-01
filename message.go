package main

import (
	"bufio"
	"fmt"
)

const (
	ECHO  = 1 // echo
	P_REG = 2 // producer register
	// other message type
	// response (simply add 100 to the initial code)
	R_ECHO  = 101 // echo response
	R_P_REG = 102 // producer register response
)

type Message struct {
	ECHO  *string
	P_REG *string
	// response
	R_ECHO  *string
	R_P_REG *byte
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
		st := string(data[1:])
		return &Message{P_REG: &st}
	case R_P_REG:
		st := data[1]
		return &Message{R_P_REG: &st}
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
		if err := writeDataToStreamWithType(streamReadWrite, *msg.P_REG, P_REG); err != nil {
			return err
		}
	} else if msg.R_P_REG != nil {
		data := fmt.Sprintf("%d", *msg.R_P_REG)
		if err := writeDataToStreamWithType(streamReadWrite, data, R_P_REG); err != nil {
			return err
		}
	}
	return nil
}
