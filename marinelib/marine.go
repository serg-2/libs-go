package marinelib

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)


func EncodeType1(mmsi uint32, speed float64, longtitude float64, latitude float64, course float64, ts int) string {
	//Message Type 1
	_type := fmt.Sprintf("%06b", int(1))
	//directive to an AIS transceiver that this message should be rebroadcast
	_repeat := "00"
	//MMSI
	_mmsi := fmt.Sprintf("%030b", mmsi)
	//Status not defined (15)
	// Status 8 = under way sailing
	_status := fmt.Sprintf("%04b", int(8))
	//rate of turn not defined (128)
	_rot := fmt.Sprintf("%08b", int(128))
	// Speed over ground is in 0.1-knot resolution from 0 to 102 knots. value 1023 indicates speed is not available, value 1022 indicates 102.2 knots or higher.
	_speed := fmt.Sprintf("%010b", int(speed*10))
	// 0 - > 10m
	// 1 - < 10m
	_accuracy := "1"

	//NB. We add a mask to tell program how long is our representation (overwise on negative integers, it cannot do the complement 2).
	_long := fmt.Sprintf("%028b", int(longtitude*600000)&0xFFFFFFF)
	_lat := fmt.Sprintf("%027b", int(latitude*600000)&0x7FFFFFF)

	// Course 0.1 resolution. Course over ground will be 3600 (0xE10) if that data is not available.
	_course := fmt.Sprintf("%012b", int(course*10))
	// # 511 (N/A)
	//_trueHeading := "111111111"
	_trueHeading := fmt.Sprintf("%09b", int(course))
	// Second of UTC timestamp
	_ts := fmt.Sprintf("%06b", int(ts))
	// "00": manufactor NaN
	// "000":  spare
	// "1": Raim flag
	_flags := "000001"
	// '11100000000000000110' : Radio status ??
	// _rstatus := "0000000000000000000" Radio status
	_rstatus := fmt.Sprintf("%019b", 49168)

	message := _type + _repeat + _mmsi + _status + _rot + _speed + _accuracy + _long + _lat + _course + _trueHeading + _ts + _flags + _rstatus
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
	return fmt.Sprintf("%02X", sum)
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

func GenerateNMEAGGA(lat, long, alt float64) string {
	var lat_s, lon_s, sign_lat, sign_lon string

	time := fmt.Sprintf(time.Now().UTC().Format("150405.00"))
	heightOfGeoid := 13.2

	lat_min := (lat - float64(int64(lat))) * 60
	lat_d := float64(int64(lat))*100 + lat_min

	lon_min := (long - float64(int64(long))) * 60
	lon_d := float64(int64(long))*100 + lon_min

	if lat_d < 0 {
		lat_s = fmt.Sprintf("%9.4f", -lat_d)
		sign_lat = "S"
	} else {
		lat_s = fmt.Sprintf("%9.4f", lat_d)
		sign_lat = "N"
	}

	if lon_d < 0 {
		lon_s = fmt.Sprintf("%010.4f", -lon_d)
		sign_lon = "W"
	} else {
		lon_s = fmt.Sprintf("%010.4f", lon_d)
		sign_lon = "E"
	}

	alt_s := fmt.Sprintf("%.2f", alt)
	geoid_s := fmt.Sprintf("%.2f", heightOfGeoid)

	//Creating AIS MESSAGE
	// 4 - RTK Quality
	// 20 - satellites
	// 1 - Horizontal dilution of position
	//     (empty field) time in seconds since last DGPS update
	//     (empty field) DGPS station ID number
	nmeaMessage := "$GPGGA," + time + "," + lat_s + "," + sign_lat + "," + lon_s + "," + sign_lon + ",4,20,1," + alt_s + ",M," + geoid_s + ",M,,"
	nmeaMessage = nmeaMessage + "*" + calculateChecksum(nmeaMessage[1:])

	return nmeaMessage
}

func ConvertXYZtoLatLongAlt(x, y, z float64, ref_lat, ref_lon, ref_z float64, ugol float64) (float64, float64, float64) {
	// XY - Right side
	// ugol - Counter clockwise

	T := ugol * math.Pi / 180

	meridian := 111134.861111
	parallel := (40075696 * math.Cos(ref_lat*math.Pi/180)) / 360

	new_x := (x*math.Cos(T)-y*math.Sin(T))/parallel + ref_lon
	new_y := (x*math.Sin(T)+y*math.Cos(T))/meridian + ref_lat
	new_z := z + ref_z
	return new_y, new_x, new_z
}
