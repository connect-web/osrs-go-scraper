package name

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/nameutils"
)

func SubmitUsernames(usernames map[string]struct{}) error {
	if len(usernames) == 0 {
		return nil
	}
	dbClient := NewNameClient()
	defer dbClient.Close()
	filteredUsernames := nameutils.Filter(usernames)

	// first make sure the usernames are not in old_players
	var validUsernames = make(map[string]struct{})

	if lowMemory {
		var validNameErr error
		validUsernames, validNameErr = dbClient.getValidUsernames(filteredUsernames)
		if validNameErr != nil {
			fmt.Printf("Failed to validate usernames %s\n", validNameErr.Error())
			return validNameErr
		}
		fmt.Printf("Removed %d old players from username list!\n", len(filteredUsernames)-len(validUsernames))
	} else {
		validUsernames = usernames
		// filter now happens earlier to ensure larger submit payloads.
		// getValidUsernamesFromMemory(filteredUsernames)
	}

	// insert the usernames
	err := dbClient.insertUsernames(validUsernames)
	if err == nil {
		fmt.Println("Successfully saved usernames.")
	}

	return err
}
