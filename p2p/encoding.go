package p2p

import (
	"encoding/gob"
	"io"
	"bytes"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct {}

func (dec GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct {}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf);
	if err != nil {
		return err
	}

	msg.Payload = bytes.Trim(buf[:n], "\x00")

	return nil
}