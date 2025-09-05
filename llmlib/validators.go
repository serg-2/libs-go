package llmlib

import (
	"log"
	"slices"
	"unicode/utf8"
)

func validateQuestion(question string) bool {
	if utf8.RuneCountInString(question) > 300 {
		log.Printf("Too big user request message!\n")
		return false
	}
	return true
}

func validateSystemMessages(previosMessages []SystemMessages) bool {
	if len(previosMessages) > 30 {
		log.Printf("Too big context!\n")
		return false
	}
	for _, message := range previosMessages {
		// Check role
		if !slices.Contains(availableRoles, message.Role) {
			log.Printf("Role %s unsupported.\n", message.Role)
			return false
		}
		// Check length?
	}
	return true
}
