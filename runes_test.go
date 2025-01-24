package main

import (
	"fmt"
	"testing"

	. "github.com/serg-2/libs-go/commonlib"
)

func TestRunelib(t *testing.T) {

	// Без вариационного селектора
	var b rune = '💰'
	// var b rune = '⚡'
	
	// Вариационный селектор (Unicode U+FE0F) сообщает, что символ нужно отображать как полноценный эмодзи, а не как обычный текст. Если вариационный селектор отсутствует, система может показать символ в текстовом виде.
	// то есть это 2 руны!
	var c string = `💰️`
	// var c string = `⚡️`
	
	// True
	fmt.Println(GetEmojiWithSelector(b) == c)
	fmt.Println(GetEmojiFromString(c, true)[0] == b)
	fmt.Println(len((GetEmojiFromString(c, false))) == 2)
	
	a := `0 ы 💰 💸💸💸💸💸💸💸💸  With selector: 💰️💰️💰️💰️💰️`
	fmt.Println(len((GetEmojiFromString(a, true))) == 14)
	fmt.Println(len((GetEmojiFromString(a, false))) == 19)

	for _,c := range a {
		fmt.Printf("Symbol: %x Len: %d Mean: %s\n",
			c,
			len([]byte(string(c))),
			string(c),
		)
	}
}
