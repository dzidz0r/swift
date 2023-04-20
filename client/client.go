package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
)

type client struct {
	conn net.Conn
	Conn net.Conn
}

func NewClient() *client {
	return &client{}
}

func (c *client) Connect(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	c.Conn = conn
	c.conn = conn
}

func (c *client) Send(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// prefix data with filename
	filename := []byte(fmt.Sprintf("%v$$$$", path.Base(filePath)))
	data = append(filename, data...)

	// send file size
	err = binary.Write(c.conn, binary.LittleEndian, int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	// send data
	i, err := io.CopyN(c.conn, bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("File sent successfully: %d / %d bytes written", i, len(data))

	return nil
}

func (c *client) Receive() {

	// receive file size
	var dataSize int64
	err := binary.Read(c.conn, binary.LittleEndian, &dataSize)
	if err != nil {
		log.Fatal(err)
	}

	// recieve data prefixed with filename
	data := new(bytes.Buffer)
	i, err := io.CopyN(data, c.conn, dataSize)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Received %d bytes from connection", i)

	// seperate data from filename
	filename, fileContent, ok := bytes.Cut(data.Bytes(), []byte("$$$$"))
	if !ok {
		log.Println("Unable to parse file... ")
	}

	err = os.WriteFile(fmt.Sprintf("./1%s", filename), fileContent, os.ModePerm)

	if err != nil {
		log.Fatal(err)
	}
}

func (c *client) Disconnect() error {
	err := c.conn.Close()
	if err != nil {
		log.Println(err)
	}
	return err
}
