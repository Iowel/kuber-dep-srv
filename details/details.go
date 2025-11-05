package details

import (
	"net"
	"os"
)

func GetHostName() (string, error) {
	hostname, _ := os.Hostname()
	return hostname, nil
}

func GetIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return nil, err
	}
	
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
