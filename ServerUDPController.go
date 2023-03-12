package main

import (
	"CSC445_Assignment2/tftp"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

// handleConnectipnUDP handles a single udp "connection"
func (c *TFTPProtocol) handleConnectionUDP() {
	buf := make([]byte, 1024)
	go func() {
		for {
			// read message
			n, raddr, err := c.conn.ReadFromUDP(buf)
			if err != nil {
				log.Println("Error reading message:", err)
				continue
			}
			// decode message
			msg := buf[:n]
			c.handleRequest(raddr, msg)
		}
	}()
}

func RunServerMode() {
	udpServer, err := NewTFTPServer()
	if err != nil {
		log.Println("Error creating server:", err)
		return
	}
	defer udpServer.Close()
	udpServer.handleConnectionUDP() // launch in separate goroutine
	select {}
}

func (c *TFTPProtocol) handleRequest(addr *net.UDPAddr, buf []byte) {
	code := binary.BigEndian.Uint16(buf[:2])
	switch tftp.TFTPOpcode(code) {
	case tftp.TFTPOpcodeRRQ:
		c.handleRRQ(addr, buf)
		break
	case tftp.TFTPOpcodeWRQ:
		log.Println("Received WRQ")
		break
	case tftp.TFTPOpcodeACK:
		log.Println("Received ACK")
		break
	case tftp.TFTPOpcodeTERM:
		log.Println("Received TERM, Terminating Transfer...")
		return
	default:

	}
}

func (c *TFTPProtocol) SetTransferSize(size uint32) {
	c.xferSize = size
}

func (c *TFTPProtocol) handleRRQ(addr *net.UDPAddr, buf []byte) {
	log.Println("Received RRQ")
	var req tftp.Request
	err := req.Parse(buf)
	if err != nil {
		c.sendError(addr, 4, "Illegal TFTP operation")
		return
	}
	log.Printf("Received %d bytes from %s for file %s \n", len(buf), addr, string(req.Filename))
	err, img := IQ.AddNewAndReturnImg(string(req.Filename))
	if err != nil {
		c.sendError(addr, 10, "File not found")
		return
	}
	c.SetProtocolOptions(req.Options, len(img))
	opAck := tftp.NewOpt(req.Options, c.xferSize)
	_, err = c.conn.WriteToUDP(opAck.ToBytes(), addr)
	if err != nil {
		log.Println("Error sending data packet:", err)
		return
	}

	c.dataBlocks, err = tftp.PrepareData(img, int(c.blockSize))
	if err != nil {
		return
	}
	fmt.Sprintf("Sending %d blocks", len(c.dataBlocks))
}

func (c *TFTPProtocol) sendError(addr *net.UDPAddr, errCode uint16, errMsg string) {
	errPack := tftp.NewErr(errCode, []byte(errMsg))
	_, err := c.conn.WriteToUDP(errPack.ToBytes(), addr)
	if err != nil {
		log.Println("Error sending error packet:", err)
		return
	}
}
