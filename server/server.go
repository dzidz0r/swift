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
	"time"
)

type server struct {
	conn          []net.Conn
	hostname      string
	serverPort    int
	listener      net.Listener
	serverTimeout time.Duration
	timer         *time.Timer
}

func NewServer() *server {
	name, _ := os.Hostname()
	return &server{
		hostname:      name,
		serverPort:    3000,
		serverTimeout: time.Second * 10,
	}
}
func (s *server) Start() {
	go func() {
		s.Broadcast()
	}()

	s.timer = time.AfterFunc(s.serverTimeout, func() {
		fmt.Println("Server timeout... shutting down")
		s.Shutdown()
		os.Exit(1)
	})
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.serverPort))
	if err != nil {
		log.Fatal(err)
	}
	s.listener = listener

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal("accepting connection err: ", err)
		}
		s.timer.Stop()
		log.Println("Connection made: ", conn)
		s.conn = append(s.conn, conn)
		conn.Write([]byte("Welcome to this "))
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

	for _, conn := range s.conn {
		// send file size
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
	s.listener.Close()
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
