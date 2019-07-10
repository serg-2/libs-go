package seriallib

import (
	"bufio"
	"github.com/tarm/serial"
	"log"
	"strconv"
	"strings"
	"time"
)

func parseString(f string) float64 {
	result, _ := strconv.ParseFloat(f, 64)
	return result
}

func GetPosition(message_type string, port string, baudrate int, stop_after_fail bool) (float64, float64, bool) {
	//Defining message_types
	var message_length byte
	var positions [2]byte
	var reply []byte
	var err_reader error
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
		return 300, 300, false
	}
	c := &serial.Config{Name: port, Baud: baudrate}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	//      n, err := s.Write([]byte("test"))
	//      if err != nil {
	//              log.Fatal(err)
	//      }

	reader := bufio.NewReader(s)

	//FIRST READ TO SKIP HALF MESSAGES
	reply, err_reader = reader.ReadBytes('\n')
	if err != nil {
		reply, err_reader = reader.ReadBytes('\n')
		if err_reader != nil {
			panic(err_reader)
		}
	}

	a := make([]string, message_length)
	message_flag := false

	for i := 0; !((a[0] == "$GP"+message_type) && (check_valid[a[positions[1]]])); i++ {
		if i > 25 && !message_flag {
			log.Printf("NO GP" + message_type + " MESSAGES\n")
			return 300, 300, false
		}
		reply, err = reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		// DEBUG BLOCK
		// reply = []byte("$GPGLL,,,,,123317.000,V,N*78")
		// reply = []byte("$GPGLL,5539.2373,N,03745.0928,E,123318.000,A,A*5A")
		// reply = []byte("$GPGLL,5539.2373,S,03745.0928,W,123318.000,A,A*5A")
		a = strings.Split(string(reply), ",")
		if a[0] == "$GP"+message_type {
			message_flag = true
			if !check_valid[a[positions[1]]] {
				log.Printf("Trying to get coordinates.GP"+message_type+" Message:  %v\n", strings.TrimSuffix(string(reply), "\n"))
				if stop_after_fail {
					return 300, 300, false
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

	return convertCoordinate(a[positions[0]]), convertCoordinate(a[positions[0]+2]), true
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

