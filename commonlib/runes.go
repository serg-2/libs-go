package commonlib

func GetEmojiFromString(s string, skipSelector bool) []rune {
	var result []rune
	for _, r := range s {
		if len([]byte(string(r))) < 3 {
			continue
		}
		if skipSelector && r == 0xFE0F {
			continue
		}
		result = append(result, r)
	}
	return result
}

func GetEmojiWithSelector(r rune) string {
	var result []rune
	result = append(result, r)
	result = append(result, 0xFE0F)
	return string(result)
}

func GetSubstringFromStringTelegram(s string, start int, length int) string {
	// log.Printf("Start %d,End: %d\n", start, start+length)
	resultString := []byte{}
	slider := 0
	var specialCase bool
	for i := 0; i < len(s); i++ {
		if (s[i] & 0xc0) != 0x80 {
			if s[i] >= 0xf0 {
				specialCase = true
			} else {
				specialCase = false
			}
		}
		if slider >= start && slider < start+length {
			if specialCase {
				resultString = append(resultString, s[i])
				resultString = append(resultString, s[i+1])
			} else {
				resultString = append(resultString, s[i])
			}
		}
		if (s[i] & 0xc0) != 0x80 {
			if s[i] >= 0xf0 {
				slider += 2
				i++
			} else {
				slider += 1
			}
		}
	}
	return string(resultString)
}
