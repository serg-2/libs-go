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
