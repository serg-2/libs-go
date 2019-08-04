package marinelib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const R = 6373000

func CalculateDistance(a [2]float64, b [2]float64) float64 {
	_diffLong := (math.Pi / 180) * math.Abs(b[1]-a[1])
	_diffLat := (math.Pi / 180) * math.Abs(b[0]-a[0])

	_a1 := math.Pow(math.Sin(_diffLat/2), 2)
	_a2 := math.Cos((math.Pi / 180) * b[0])
	_a3 := math.Cos((math.Pi / 180) * a[0])
	_a4 := math.Pow(math.Sin(_diffLong/2), 2)
	_a := _a1 + _a2*_a3*_a4
	_c := 2 * math.Atan2(math.Sqrt(_a), math.Sqrt(1-_a))
	return R * _c
}

func CalculateBearing(a [2]float64, b [2]float64) float64 {
	_lat1 := (math.Pi / 180) * a[0]
	_lat2 := (math.Pi / 180) * b[0]
	_diffLong := (math.Pi / 180) * (b[1] - a[1])

	_x := math.Sin(_diffLong) * math.Cos(_lat2)
	_y := math.Cos(_lat1)*math.Sin(_lat2) - (math.Sin(_lat1) * math.Cos(_lat2) * math.Cos(_diffLong))

	_initial_bearing := math.Atan2(_x, _y)
	_initial_bearing = _initial_bearing * (180 / math.Pi)
	compass_bearing := math.Mod(_initial_bearing+360, 360)
	return compass_bearing
}

func EncodeType1(mmsi uint32, speed float64, longtitude float64, latitude float64, course float64, ts int) string {
	//Message Type 1
	_type := fmt.Sprintf("%06b", int(1))
	//directive to an AIS transceiver that this message should be rebroadcast
	_repeat := "00"
	//MMSI
	_mmsi := fmt.Sprintf("%030b", mmsi)
	//Status not defined (15)
	_status := fmt.Sprintf("%04b", int(15))
	//rate of turn not defined (128)
	_rot := fmt.Sprintf("%08b", int(128))
	// Speed over ground is in 0.1-knot resolution from 0 to 102 knots. value 1023 indicates speed is not available, value 1022 indicates 102.2 knots or higher.
	_speed := fmt.Sprintf("%010b", int(speed*10))
	// > 10m
	_accurancy := "0"

	//NB. We add a mask to tell program how long is our representation (overwise on negative integers, it cannot do the complement 2).
	_long := fmt.Sprintf("%028b", int(longtitude*600000)&0xFFFFFFF)
	_lat := fmt.Sprintf("%027b", int(latitude*600000)&0x7FFFFFF)

	// Course 0.1 resolution. Course over ground will be 3600 (0xE10) if that data is not available.
	_course := fmt.Sprintf("%012b", int(course*10))
	// # 511 (N/A)
	_trueHeading := "111111111"
	// Second of UTC timestamp
	_ts := fmt.Sprintf("%06b", int(ts))
	// "00": manufactor NaN
	// "000":  spare
	// "0": Raim flag
	_flags := "000000"
	// '11100000000000000110' : Radio status ??
	_rstatus := "0000000000000000000"

	message := _type + _repeat + _mmsi + _status + _rot + _speed + _accurancy + _long + _lat + _course + _trueHeading + _ts + _flags + _rstatus
	return message
}

func EncodeType24(mmsi uint32, part string, name string, callsign string, vsize string, vtype int) string {
	var message string
	//Type 24
	_type := fmt.Sprintf("%06b", int(24))
	//directive to an AIS transceiver that this message should be rebroadcast.  00 = default; 11 = do not repeat any more
	_repeat := "00"
	//MMSI
	_mmsi := fmt.Sprintf("%030b", mmsi)
	if part == "A" {
		//Part of Message
		_part := "00"
		//Vessel Name. Maximum 120 bit.20 characters 6-bit ASCII. Default = not available = @@@@@@@@@@@@@@@@@@@@
		_name := encodeString(name)
		//Padding. 160 bits per RFC + 8 bits spare.
		npadding := 168 - len(_type) - len(_repeat) - len(_mmsi) - len(_part) - len(_name)
		_padding := strings.Repeat("0", npadding)

		message = _type + _repeat + _mmsi + _part + _name + _padding

	} else {
		//Part of Message
		_part := "01"
		// Vessel Type (Number of Passengers).
		// 0 = not available or no ship = default. 1-99 = as defined in § 3.3.2. 100-199 = reserved, for regional use. 200-255 = reserved, for future use
		_vtype := fmt.Sprintf("%08b", vtype)
		//Vendor ID. Default @@@@@@@. MSB. 18 bits(3 symbols) - manufacture’s mnemonic code. 4 bits - Unit Model Code.20 bits - Serial Number
		_vendorID := strings.Repeat("0", 42)
		// Call Sign. 7 six-bit characters. @@@@@@@ = not available = default. Craft associated with a parent vessel should use “A” followed by the last 6 digits of the MMSI of the parent vessel.
		if len(callsign) > 7 {
			panic("Length of CallSign exceeds 7 characters")
		}
		csign := encodeString(callsign)
		_callsign := csign + strings.Repeat("0", 42-len(csign))
		// AIS antenna in the middle of the boat. Dimensions 30 bit. Reference point for reported position
		_length, err := strconv.Atoi(strings.Split(vsize, "x")[0])
		if err != nil {
			panic("Unable to detect vessel size - length")
		}
		_width, err := strconv.Atoi(strings.Split(vsize, "x")[1])
		if err != nil {
			panic("Unable to detect vessel size - width")
		}
		// NEXT 30 bits - DIMENSTIONS for main craft. OR! For auxiliary craft - MMSI of MotherSHIP
		//First - Dimension to Bow. Second - Dimenstion to Stern.
		_half_length := fmt.Sprintf("%09b", _length/2)
		//First - Dimension to Port. Second - Dimenstion to Startboard.
		_half_width := fmt.Sprintf("%06b", _width/2)
		//Type of electronic position fixing device. 0 -Default. 1 = GPS, 2 = GLONASS, 3 = combined GPS/GLONASS, 4 = Loran-C, 5 = Chayka, 6 = integrated navigation system, 7 = surveyed; 8 = Galileo, 9-14 = not used, 15 = internal GNSS
		// RARE Parameters. Not supported usually
		//_gps_type := "0111"
		_gps_type := "0000"
		//SPARE
		_spare := "00"
		//168 bit
		message = _type + _repeat + _mmsi + _part + _vtype + _vendorID + _callsign + _half_length + _half_length + _half_width + _half_width + _gps_type + _spare

	}
	return message
}

func calculateChecksum(sentence string) string {
	//The NMEA checksum is computed on the entire sentence including the AIVDM/AIVDO tag but excluding the leading "!"
	//The checksum is merely a byte-by-byte XOR of the sentence
	var sum int
	for i := 0; i < len(sentence); i++ {
		sum = sum ^ int(sentence[i])
	}
	return fmt.Sprintf("%X", sum)
}

func encodeString(sentence string) string {
	var encodedString string = ""
	var index int
	const vocabulary = "@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^- !\"#$%&'()*+,-./0123456789:;<=>?"
	for _, char := range strings.ToUpper(sentence) {
		index = strings.Index(vocabulary, string(char))
		encodedString = encodedString + fmt.Sprintf("%06b", index)
	}
	return encodedString
}

func GenerateNMEA(tmp_string string) string {
	in_string := []byte(tmp_string)
	strlen := len(in_string)

	//Check Valid input. But for 24A this could be wrong!!! TODO
	if strlen%6 != 0 {
		panic("Input length not an even multiple of 6...")
	}

	for i := 0; i < strlen; i++ {

		if !(string(in_string[i]) == "0" || string(in_string[i]) == "1") {
			panic("Input contains non-binary value")
		}

	}

	channel := "A"

	var in_buffer []int64
	var armored_message string

	//Getting Armored ASCII
	for i := 0; i < strlen/6; i++ {
		//One character = 6 bits, so get the 6-bit block
		var stringtoparse string
		for j := 0; j < 6; j++ {
			stringtoparse = stringtoparse + string([]byte(in_string[i*6+j:i*6+j+1]))
		}
		converted, err := strconv.ParseInt(string(in_string[i*6:i*6+6]), 2, 64)
		if err != nil {
			panic("Can't convert to integer")
		}
		in_buffer = append(in_buffer, converted)
		//Convert to Char. Special rules
		if in_buffer[i] > 39 {
			in_buffer[i] = in_buffer[i] + 8
		}
		in_buffer[i] = in_buffer[i] + 48
		//Creating Message
		armored_message = armored_message + string(in_buffer[i])
	}
	//DEBUG: Print Armored message
	//fmt.Println(armored_message)

	//Creating AIS MESSAGE
	nmeaMessage := "!AIVDM,1,1,," + channel + "," + armored_message + ",0"
	nmeaMessage = nmeaMessage + "*" + calculateChecksum(nmeaMessage[1:])

	return nmeaMessage
}
