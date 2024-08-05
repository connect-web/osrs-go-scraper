package name

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/nameutils"
)

var knownUsernames = map[string]struct{}{}
var lowMemory = true

/*
low_memory mode will query database for names that should be added

normal mode will have a map of known usernames and do filtering via in memory operations.
*/

func LoadKnownUsernames() {
	dbClient := NewNameClient()
	defer dbClient.Close()
	_ = dbClient.getAllUsernames() // no need to handle it will exit.
	lowMemory = false
	fmt.Printf("Successfully loaded %d usernames into memory.", len(knownUsernames))
}

func getValidUsernamesFromMemory(usernames map[string]struct{}) map[string]struct{} {
	validUsernames := map[string]struct{}{}

	for username := range usernames {
		_, exists := knownUsernames[username]
		if !exists {
			validUsernames[username] = struct{}{}
		}
	}

	return validUsernames
}

func IsUsernameKnown(username string) bool {
	_, exists := knownUsernames[nameutils.FilterName(username)]
	return exists
}
