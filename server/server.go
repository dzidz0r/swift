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
	conn       []net.Conn
	hostname   string
	listener   net.Listener
	serverPort int
}

func NewServer() *server {
	name, _ := os.Hostname()
	return &server{
		hostname:   name,
		serverPort: _PORT,
	}
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
		log.Println("Connection made: ", conn)
		s.conn = append(s.conn, conn)
		go s.readLoop(conn)
	}
}

func (s *server) readLoop(conn net.Conn) {
	defer conn.Close()

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
		log.Printf("Received %d bytes from %v", i, conn.RemoteAddr())

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

func (s *server) Send(filePath string) error {
	fmt.Println("sending")
	fmt.Println(s.conn)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Println(err)
		return err
	}

	// prefix data with filename
	filename := []byte(fmt.Sprintf("%v$$$$", path.Base(filePath)))
	data = append(filename, data...)

	// send file size
	for _, conn := range s.conn {
		go func(conn net.Conn) {
			err = binary.Write(conn, binary.LittleEndian, int64(len(data)))
			if err != nil {
				log.Println(err)
				return
			}

			// send data
			i, err := io.CopyN(conn, bytes.NewReader(data), int64(len(data)))
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("File sent successfully: %d / %d bytes written", i, len(data))
		}(conn)
	}

	return nil
}

func (s *server) Shutdown() {
	defer fmt.Println("all connections closed")
	for _, conn := range s.conn {
		err := conn.Close()
		if err != nil {
			continue
		}
	}
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
