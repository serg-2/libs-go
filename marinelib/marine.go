package marinelib

import "math"

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

func EncodeType1(mmsi int, speed float64, longtitude float64, latitude float64, course float64, ts int) string {
	//Type 18 (1)
	_typeVessel := fmt.Sprintf("%06b", int(1))
	//directive to an AIS transceiver that this message should be rebroadcast
	_repeat := "00"
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

	message := _typeVessel + _repeat + _mmsi + _status + _rot + _speed + _accurancy + _long + _lat + _course + _trueHeading + _ts + _flags + _rstatus
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

func GenerateNMEA(tmp_string string) string {
	in_string := []byte(tmp_string)
	strlen := len(in_string)

	//Check Valid input
	if strlen%6 != 0 {
		fmt.Println("Input length not an even multiple of 6...")
		os.Exit(1)
	}

	for i := 0; i < strlen; i++ {

		if !(string(in_string[i]) == "0" || string(in_string[i]) == "1") {
			fmt.Println("Input contains non-binary value")
			os.Exit(2)
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
			fmt.Println("Can't convert to integer")
			os.Exit(3)
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
