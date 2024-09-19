// TODO: set up basic udp server to send a simple msg to the user
package server

import (
	"log"
	"net"
	"os"
)

type Server struct {
	Port   string
	Ln     *net.UDPConn
	Quitch chan struct{}
}

type Client struct {
	Address *net.UDPAddr
}

func newClient(addr *net.UDPAddr) *Client {
	return &Client{
		Address: addr,
	}
}

func (s *Server) sendVideo(c *Client, video []byte) error {
	const packetSize = 1472 // Safe size for UDP payloads, considering overhead
	n := len(video)
	//the formula has a -1 which makes it such that the total size if it divisible with packet size gives a value of n+1 but the -1/ps makes it less than n+1, otherwise it will be rounded to n+1
	numPackets := (n + packetSize - 1) / packetSize
	totalBytes := 0
	for i := 0; i < numPackets; i++ {
		start := i * packetSize
		end := start + packetSize
		if end > n {
			end = n
		}

		sentBytes, err := s.Ln.WriteToUDP(video[start:end], c.Address)
		if err != nil {
			log.Println("Error sending packet:", err)
			return err
		}
		totalBytes += sentBytes
	}

	log.Printf("Total bytes sent: %d out of %d bytes", totalBytes, n)
	return nil
}

func CreateServer(addr string) *Server {
	return &Server{
		Port:   addr,
		Quitch: make(chan struct{}),
	}
}

func (s *Server) StartServer() {
	addr, err := net.ResolveUDPAddr("udp4", s.Port)
	if err != nil {
		log.Println("Invalid port, using :8080 instead")
		addr, err = net.ResolveUDPAddr("udp4", ":8080")
		if err != nil {
			log.Fatal("Error resolving address")
		}
	}

	s.Ln, err = net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Ln.Close()
	log.Printf("Server started on %d", addr.Port)
	go s.acceptLoop()
	<-s.Quitch
}

func (s *Server) acceptLoop() {
	//provides header for the request
	req := make([]byte, 1024)

	n, addr, err := s.Ln.ReadFromUDP(req)
	if err != nil {
		log.Println("Error accepting connection...")
	}
	log.Println("Connected to addr", addr.IP)
	log.Println("bytes:", n, "Message received: ", string(req))
	client := newClient(addr)
	// read the mp4 File
	videoBytes, err := ReadVideo("../videos/test.mp4")
	if err != nil {
		log.Println("Error reading video", err)
	}
	err = s.sendVideo(client, videoBytes)
	if err != nil {
		log.Println("Error sending video", err)
	}
	log.Println("Video sent sucessfully")
}

func ReadVideo(filename string) ([]byte, error) {
	videoBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return videoBytes, nil
}
