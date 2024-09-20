package marinelib

func GetRhumb(angle float64) string {
	switch {
	case angle < 11.25*0+5.625:
		return "N"
	case angle < 11.25*1+5.625:
		return "NtE"
	case angle < 11.25*2+5.625:
		return "NNE"
	case angle < 11.25*3+5.625:
		return "NEtN"
	case angle < 11.25*4+5.625:
		return "NE"
	case angle < 11.25*5+5.625:
		return "NEtE"
	case angle < 11.25*6+5.625:
		return "ENE"
	case angle < 11.25*7+5.625:
		return "EtN"
	case angle < 11.25*8+5.625:
		return "E"
	case angle < 11.25*9+5.625:
		return "EtS"
	case angle < 11.25*10+5.625:
		return "ESE"
	case angle < 11.25*11+5.625:
		return "SEtE"
	case angle < 11.25*12+5.625:
		return "SE"
	case angle < 11.25*13+5.625:
		return "SEtS"
	case angle < 11.25*14+5.625:
		return "SSE"
	case angle < 11.25*15+5.625:
		return "StE"
	case angle < 11.25*16+5.625:
		return "S"
	case angle < 11.25*17+5.625:
		return "StW"
	case angle < 11.25*18+5.625:
		return "SSW"
	case angle < 11.25*19+5.625:
		return "SWtS"
	case angle < 11.25*20+5.625:
		return "SW"
	case angle < 11.25*21+5.625:
		return "SWtW"
	case angle < 11.25*22+5.625:
		return "WSW"
	case angle < 11.25*23+5.625:
		return "WtS"
	case angle < 11.25*24+5.625:
		return "W"
	case angle < 11.25*25+5.625:
		return "WtN"
	case angle < 11.25*26+5.625:
		return "WNW"
	case angle < 11.25*27+5.625:
		return "NWtW"
	case angle < 11.25*28+5.625:
		return "NW"
	case angle < 11.25*29+5.625:
		return "NWtN"
	case angle < 11.25*30+5.625:
		return "NNW"
	case angle < 11.25*31+5.625:
		return "NtW"
	default:
		return "N"
	}
}
