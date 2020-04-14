package seriallib

import (
	"bufio"
	"encoding/hex"
	"github.com/tarm/serial"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func parseString(f string) float64 {
	result, _ := strconv.ParseFloat(f, 64)
	return result
}

func GetPosition(message_type string, port string, baudrate int, stop_after_fail bool) ([2]float64, bool) {
	//Defining message_types
	var message_length byte
	var positions [2]byte
	var reply []byte
	var err_reader,err error
	check_valid := map[string]bool{
		"A": true,
		"1": true,
		"2": true,
		"3": true,
	}

	switch {
	case message_type == "GGA" || message_type == "1":
		message_length = 15
		message_type = "GGA"
		positions = [2]byte{2, 6}
	case message_type == "RMC" || message_type == "2":
		message_length = 12
		message_type = "RMC"
		positions = [2]byte{3, 2}
	case message_type == "GLL" || message_type == "3":
		message_length = 8
		message_type = "GLL"
		positions = [2]byte{1, 6}
	default:
		log.Printf("Message type not supported\n")
		return [2]float64{300, 300}, false
	}

	reader := GetPortReader(port, baudrate)

	//FIRST READ TO SKIP HALF MESSAGES
	reply, err_reader = reader.ReadBytes('\n')
	if err_reader != nil {
		reply, err_reader = reader.ReadBytes('\n')
		if err_reader != nil {
			panic(err_reader)
		}
	}

	a := make([]string, message_length)
	message_flag := false

	for i := 0; !(((a[0] == "$GP"+message_type) || (a[0] == "$GN"+message_type)) && (check_valid[a[positions[1]]])); i++ {
		if i > 50 && !message_flag {
			log.Printf("NO GP" + message_type + " OR NO GN" + message_type + " MESSAGES\n")
			return [2]float64{300, 300}, false
		}
		reply, err = reader.ReadBytes('\n')

		if err != nil {
			panic(err)
		}
		// DEBUG BLOCK
		// BAD
		// reply = []byte("$GPGLL,,,,,123317.000,V,N*78")
		// reply = []byte("$GPGGA,122934.00,3346.78434,N,01634.270,5,E,1,07,4.24,152.0,M,13.2,M,,*51")
		// GOOD
		// reply = []byte("$GPGLL,5547.1663,N,03246.0837,E,123318.000,A,A*<CHECKSUM>")
		// reply = []byte("$GPGGA,123854.00,5427.87439,N,03333.12340,E,1,06,7.67,160.4,M,13.2,M,,*<CHECKSUM>")
		// reply = []byte("$GPRMC,124402.00,A,5422.27361,N,03531.48171,E,0.020,,280819,,,A*<CHECKSUM>")

		a = strings.Split(string(reply), ",")

		if len(a[0]) < 6 {
			log.Printf("Bad message ID: [%v]\n", a[0])
			continue
		}

		if a[0][3:6] == message_type {
			message_flag = true
			if !check_valid[a[positions[1]]] {
				log.Printf("Trying to get coordinates.G"+string([]byte{a[0][2]})+message_type+" Message:  %v\n", strings.TrimSuffix(string(reply), "\n"))
				if stop_after_fail {
					return [2]float64{300, 300}, false
				}
				time.Sleep(1 * time.Second)
			}
		}
	}
	if a[positions[0]+1] == "S" {
		a[positions[0]] = "-" + a[positions[0]]
	}
	if a[positions[0]+2+1] == "W" {
		a[positions[0]+2] = "-" + a[positions[0]+2]
	}

	return [2]float64{convertCoordinate(a[positions[0]]), convertCoordinate(a[positions[0]+2])}, true
}

func convertCoordinate(x string) float64 {

	coord1 := strings.Split(x, ".")
	gradusy := coord1[0][:len(coord1[0])-2]
	minuty := coord1[0][len(coord1[0])-2:] + "." + coord1[1]
	if string(coord1[0][0]) == "-" {
		minuty = "-" + minuty
	}

	x_final := parseString(gradusy) + parseString(minuty)/60

	answer := x_final
	return answer
}

func GetSkyTraqMessage(reader *bufio.Reader, debug bool) ([]byte, bool) {

	var partMessage byte

	var checkMessage, fullMessage []byte

	begin1, _ := hex.DecodeString("0D")
	begin2, _ := hex.DecodeString("0A")

	begin3, _ := hex.DecodeString("A0")
	begin4, _ := hex.DecodeString("A1")

	frameString,_ := hex.DecodeString("0D0AA0A1")

	// wait for start

	// Accumulating 4 bytes to check
	for len(checkMessage) < 4 {
		partMessage, _ = reader.ReadByte()
		checkMessage = append(checkMessage, partMessage)
		fullMessage = append(fullMessage, partMessage)
	}

	// Looking for separator between packets
	for !(checkMessage[0] == frameString[0] && checkMessage[1] == frameString[1] && checkMessage[2] == frameString[2] && checkMessage[3] == frameString[3]) {
		if debug{
			log.Printf("%v\n", checkMessage)
		}
		checkMessage = checkMessage[1:]
		partMessage, _ = reader.ReadByte()
		checkMessage = append(checkMessage, partMessage)
		fullMessage = append(fullMessage, partMessage)
	}
	//fmt.Println("Packet Border found")

	fullMessage = fullMessage[:len(fullMessage)-4]
	fullMessage = append(begin4, fullMessage...)
	fullMessage = append(begin3, fullMessage...)
	fullMessage = append(fullMessage, begin1[0])
	fullMessage = append(fullMessage, begin2[0])

	// Check CRC
	payload := fullMessage[4:len(fullMessage)-3]
	crc:= calculateCRCSkytraq(payload)
	if crc != fullMessage[len(fullMessage)-3] {
		log.Println("CRC NOT OK!!!")
		return fullMessage, false
	}

	return fullMessage, true
}

func calculateCRCSkytraq(payload []byte) byte {
	var crc byte
	crc = 0
	for i := 0; i < len(payload); i++ {
		crc ^= payload[i]
	}
	//fmt.Printf("Received Payload: %v\n", payload)

	return crc
}

func GetPortReader(nameOfPort string, baudRate int) *bufio.Reader {
	serialPort, err := serial.OpenPort(&serial.Config{Name: nameOfPort, Baud: baudRate})

	// Handle errors
	for err != nil {
		// Check file exists and have access
		_,check1 := err.(*os.PathError)
		if check1 {
			log.Printf("Can't open serial port from config.\nWaiting 10 secs...\n")
			time.Sleep(10 * time.Second)
			continue
		}
		// Check propriate serial device
		_,check2 := err.(syscall.Errno)
		if check2 {
			log.Printf("Serial port from configuration is not serial port\nWaiting 10 secs...\n")
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("Unable to handle error.\n")
		log.Fatal(err)
	}

	// log.Printf("Serial port opened.\n")

	// To Write something
	//      n, err := serialPort.Write([]byte("test"))
	//      if err != nil {
	//              log.Fatal(err)
	//      }

	return bufio.NewReader(serialPort)
}