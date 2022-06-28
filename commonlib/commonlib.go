package commonlib

import (
	"context"
	"errors"
	"github.com/lib/pq"
	"log"
	"net"
	"strconv"
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
