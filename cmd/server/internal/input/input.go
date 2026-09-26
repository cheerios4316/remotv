package input

import "flag"

type UserInput struct {
	TcpPort  int
	HttpPort int
}

func ParseFlags() UserInput {
	var tcpPort, httpPort int

	flag.IntVar(&tcpPort, "tcp", 9000, "Port to listen for a TCP connection on")
	flag.IntVar(&tcpPort, "T", 9000, "Port to listen for a TCP connection on (shorthand)")

	flag.IntVar(&httpPort, "http", 8080, "Port to listen for HTTP requests on")
	flag.IntVar(&httpPort, "H", 8080, "Port to listen for HTTP requests on (shorthand)")

	flag.Parse()

	return UserInput{
		TcpPort:  tcpPort,
		HttpPort: httpPort,
	}
}
