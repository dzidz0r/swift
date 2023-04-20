package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const (
	_PORT = 5050
)

type server struct {
	conn net.Conn
}

func NewServer() *server {
	return &server{}
}

func (s *server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", _PORT))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Connection established: %v", conn)
			s.conn = conn
			break
		}
	}
}

func (s *server) Send(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// send file size
	err = binary.Write(s.conn, binary.LittleEndian, int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	// send file
	i, err := io.CopyN(s.conn, bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Written %d/%d bytes to the connection", i, len(data))

	return nil
}

func (s *server) Receive() {

	// receive file size
	var dataSize int64
	err := binary.Read(s.conn, binary.LittleEndian, &dataSize)
	if err != nil {
		log.Fatal(err)
	}

	// recieve data prefixed with filename
	data := new(bytes.Buffer)
	i, err := io.CopyN(data, s.conn, dataSize)
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
