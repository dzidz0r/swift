package server

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

const (
	_PORT = 5050
)

type server struct {
	addr     net.Addr
	hostname string
	listener net.Listener
}

func NewServer() *server {
	name, _ := os.Hostname()
	return &server{
		hostname: name,
	}
}

func (s *server) Broadcast() {

}

func (s *server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", _PORT))
	if err != nil {
		log.Fatal(err)
	}
	s.listener = listener

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal("accepting connection err: ", err)
		}
		go s.readLoop(conn)
	}
}

func (s *server) readLoop(conn net.Conn) {
	for {

		// receive file size
		var dataSize int64
		err := binary.Read(conn, binary.LittleEndian, &dataSize)
		if err != nil {
			log.Fatal(err)
		}

		// recieve data prefixed with filename
		data := new(bytes.Buffer)
		i, err := io.CopyN(data, conn, dataSize)
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

}

func (s *server) Send(conn net.Conn, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// prefix data with filename
	filename := []byte(fmt.Sprintf("%v$$$$", path.Base(filePath)))
	data = append(filename, data...)

	// send file size
	err = binary.Write(conn, binary.LittleEndian, int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	// send data
	i, err := io.CopyN(conn, bytes.NewReader(data), int64(len(data)))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("File sent successfully: %d / %d bytes written", i, len(data))

	return nil
}

// func (s *server) Receive() {

// 	// receive file size
// 	var dataSize int64
// 	err := binary.Read(s.conn, binary.LittleEndian, &dataSize)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// recieve data prefixed with filename
// 	data := new(bytes.Buffer)
// 	i, err := io.CopyN(data, s.conn, dataSize)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Printf("Received %d bytes from connection", i)

// 	// seperate data from filename
// 	filename, fileContent, ok := bytes.Cut(data.Bytes(), []byte("$$$$"))
// 	if !ok {
// 		log.Println("Unable to parse file... ")
// 	}

// 	err = os.WriteFile(fmt.Sprintf("./1%s", filename), fileContent, os.ModePerm)

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
