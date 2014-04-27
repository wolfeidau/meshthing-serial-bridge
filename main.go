package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/wolfeidau/meshthing-serial-bridge/pcap"
	"github.com/wolfeidau/meshthing-serial-bridge/rs232"
)

var portFlag = flag.String("port", "", "optional path to serial device")
var baudFlag = flag.Int("baud", 38400, "optional baud of the serial device")
var versionFlag = flag.Bool("version", false, "print the version information")

func findArduino() string {

	if *portFlag != "" {
		return *portFlag
	}

	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "tty.usbmodem") ||
			strings.Contains(f.Name(), "ttyUSB") {
			return "/dev/" + f.Name()
		}
	}

	// Have not been able to find a USB device that 'looks'
	// like an Arduino.
	return ""
}

func main() {

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Bridge v%s\n", Version)
		fmt.Printf("Bridge commit %s\n", GitCommit)
		os.Exit(0)
	}

	// Find the device that represents the arduino serial
	// connection.
	port, err := rs232.OpenPort(findArduino(), *baudFlag, rs232.S_8N1)

	if err != nil {
		log.Fatalf("Err opening port: %s", err)
	}

	log.Printf("port %v", port)

	defer port.Close()

	err = mkfifo("/tmp/wireshark")

	if err != nil {
		log.Fatalf("Error createing fifo: %s", err)
	}

	f, err := os.OpenFile("/tmp/wireshark", os.O_RDWR, 0666)

	if err != nil {
		log.Fatalf("Error opening fifo: %s", err)
	}

	w, err := pcap.NewWriter(f, pcap.NewFileHeader(pcap.DLT_IEEE802_15_4))

	if err != nil {
		log.Fatalf("Error creating pcap writer: %s", err)
	}

	log.Println("reading from serial")
	r := bufio.NewReader(&port)

	for {
		line, _, err := r.ReadLine()

		if err != nil {
			log.Fatalf("Error reading port: %s", err)
		}

		if string(line[:]) == "Online" {
			log.Println("Sniffer Online")
			break
		}
	}

	log.Println("write pcap file header")

	buf := make([]byte, 1024)

	// main read and write loop which converts the raw packet from the serial port to
	// a pcap packet and writes it out the unix fifo.
	for {
		plen, err := r.ReadByte()

		if err != nil {
			log.Fatalf("Error reading packet length: %s", err)
		}

		_, err = io.ReadFull(r, buf[:plen])

		if err != nil {
			log.Fatalf("Error reading packet: %s", err)
		}

		pkt := pcap.NewPacket(buf[:plen], uint32(plen))

		err = w.Write(pkt)

		if err != nil {
			log.Fatalf("Error writing packet: %s", err)
		}
	}
}

func mkfifo(path string) error {
	err := syscall.Mknod(path, syscall.S_IFIFO|0666, 0)
	if os.IsExist(err) {
		return nil
	}
	return err
}
