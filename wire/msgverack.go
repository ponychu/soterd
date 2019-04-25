// Copyright (c) 2013-2015 The btcsuite developers
// Copyright (c) 2018-2019 The Soteria Engineering developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"io"
)

// MsgVerAck defines a soter verack message which is used for a peer to
// acknowledge a version message (MsgVersion) after it has used the information
// to negotiate parameters.  It implements the Message interface.
//
// This message has no payload.
type MsgVerAck struct{}

// SotoDecode decodes r using the soter protocol encoding into the receiver.
// This is part of the Message interface implementation.
func (msg *MsgVerAck) SotoDecode(r io.Reader, pver uint32, enc MessageEncoding) error {
	return nil
}

// SotoEncode encodes the receiver to w using the soter protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgVerAck) SotoEncode(w io.Writer, pver uint32, enc MessageEncoding) error {
	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgVerAck) Command() string {
	return CmdVerAck
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgVerAck) MaxPayloadLength(pver uint32) uint32 {
	return 0
}

// NewMsgVerAck returns a new soter verack message that conforms to the
// Message interface.
func NewMsgVerAck() *MsgVerAck {
	return &MsgVerAck{}
}
