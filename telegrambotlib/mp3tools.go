package telegrambotlib

import (
	"os"

	tag "github.com/dhowden/tag"
	cl "github.com/serg-2/libs-go/commonlib"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// Read MP3 tag by filename
func ReadTag(filename string) (string, string) {
	f, err := os.Open(filename)
	cl.ChkNonFatal(err)
	defer f.Close()

	m, err := tag.ReadFrom(f)
	cl.ChkNonFatal(err)

	// fmt.Println("Format", m.Format())
	//fmt.Printf("RAW: %v\n", m.Raw()["TPE1"])

	// DETECT ENCODING
	// import "github.com/mikkyang/id3-go"
	// mp3File, _ := id3.Open(filename)
	// defer mp3File.Close()
	// fmt.Println("Encoding:", mp3File.Frame("TPE1").(*v2.TextFrame).Encoding())

	// Using Artist and Title field
	if m.Format() == tag.ID3v1 {
		return cp1251ToUtf8(m.Artist()), cp1251ToUtf8(m.Title())
	} else if m.Format() == tag.ID3v2_3 {
		return V23convert(m.Artist()), V23convert(m.Title())
	} else {
		return m.Artist(), m.Title()
	}
}

// MP3 id3v2_3 convert
func V23convert(stringFrom string) string {
	// Transform from LATIN1. Because library TAG already decoded frame as LATIN1, what was wrong
	latin1, _, err := transform.String(charmap.ISO8859_1.NewEncoder(), stringFrom)
	// Check error of encode
	if err != nil {
		latin1 = stringFrom
	}

	// Decode from windows1251
	decoder := charmap.Windows1251.NewDecoder()
	res, _ := decoder.String(latin1)

	// If error res will be ""
	if len(res) == 0 {
		return latin1
	}
	return res
}

// Converting charset cp1251 to UTF8
func cp1251ToUtf8(stringFrom string) string {
	// MAIN
	decoder := charmap.Windows1251.NewDecoder()
	res, _ := decoder.String(stringFrom)
	return res
}
