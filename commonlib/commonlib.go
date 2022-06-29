package commonlib

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// ChkFatal - check and exit upon Fatal
func ChkFatal(err error) {
	if err != nil {
		switch err.(type) {
		case *net.OpError:
			// Error during querying database
			log.Fatalln("NETWORK ERROR!")
		default:
			// Other errors
			log.Printf("%T\n", err)
			log.Fatal(err)
		}
	}
}

// IntArrayToString - Convert integer array to string with delimiters
func IntArrayToString(a []int32, delimiter string) string {
	b := ""
	for _, v := range a {
		if len(b) > 0 {
			b += delimiter
		}
		b += strconv.Itoa(int(v))
	}
	return b
}

func ChkNonFatal(err error) {
	if err != nil {
		switch err.(type) {
		case *net.OpError:
			// Error during querying database
			log.Println("NETWORK ERROR!")
		case *pq.Error:
			pqerr, _ := err.(*pq.Error)
			if pqerr.Code == "57014" {
				log.Println("(TIMEOUT) Request to DB was canceled.")
				return
			} else if pqerr.Code == "23505" {
				log.Println("(CONSTRAINT) Unique violation constraint!")
				return
			} else if pqerr.Code == "42P01" {
				log.Println("(CONSTRAINT) No such table!")
				return
			} else {
				log.Printf("PQ Unknown Error Code: %v\n", pqerr.Code)
				log.Println(err)
			}
		default:
			if errors.Is(err, context.DeadlineExceeded) {
				log.Println("(TIMEOUT)Context Deadline exceeded.")
				return
			}
			// Other errors
			log.Printf("%T\n", err)
			log.Println(err)
		}
	}
}

// CheckStringInSlice - check string in slice
func CheckStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetIntString(s string, lengthString int, negativePossible bool, maxValue int) (int, bool) {
	if len(s) > lengthString {
		return 0, false
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	if !negativePossible {
		if result < 0 {
			return 0, false
		}
	}
	if result > maxValue {
		return 0, false
	}

	return result, true
}

// GetStringWithCheck - check string for length, regex(optional, "" - no regex), with change register(optional)
func GetStringWithCheck(s string, lengthString int, regexPatternForString string, toLower bool, toUpper bool) (string, bool) {
	// Check appropriate run parameters
	if toUpper && toLower {
		return "", false
	}
	// Check length
	if len(s) > lengthString {
		return "", false
	}
	if regexPatternForString != "" {
		matched, err := regexp.MatchString(regexPatternForString, s)
		// BAD REGEXP
		if err != nil {
			return "", false
		}
		// Not Matched
		if !matched {
			return "", false
		}
	}

	if toLower {
		return strings.ToLower(s), true
	}
	if toUpper {
		return strings.ToUpper(s), true
	}

	return s, true
}
