package players

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/connect-web/Low-Latency/internal/requests"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"net/url"
)

var (
	HiscoreTypes = map[string]string{
		"Normal":           "hiscore_oldschool",
		"Skiller":          "hiscore_oldschool_skiller",
		"1 Defence":        "hiscore_oldschool_skiller_defence",
		"Fresh start":      "hiscore_oldschool_fresh_start",
		"Seasonal Deadman": "hiscore_oldschool_deadman",
		"Leagues":          "hiscore_oldschool_seasonal",
		"Tournament":       "hiscore_oldschool_tournament",

		"Iron":                "hiscore_oldschool_ironman",
		"Ultimate Iron":       "hiscore_oldschool_ultimate",
		"Hardcore Iron":       "hiscore_oldschool_hardcore_ironman",
		"Group Iron":          "hiscore_oldschool_ironman/group-ironman",
		"Hardcore Group Iron": "hiscore_oldschool_hardcore_ironman/group-ironman",
	} // The m= part of the URL
)

func GetPlayerStats(placeholder playerutils.SimplePlayer, iterator *entities.ProxyIterator) (playerutils.SimplePlayer, error) {
	urlEncodedName := url.QueryEscape(placeholder.Username)
	params := map[string]string{"player": urlEncodedName}

	response, err := requests.Request(fmt.Sprintf("https://secure.runescape.com/m=%s/index_lite.json", "hiscore_oldschool"), params, iterator, -99)

	if err == nil {
		var responseJson playerutils.JsonResponse
		unmarshalErr := json.Unmarshal([]byte(response), &responseJson)
		if unmarshalErr != nil {
			fmt.Printf(unmarshalErr.Error())
			return placeholder, errors.New("failed to load json")
		}
		if !placeholder.LoadJson(responseJson) {
			return placeholder, errors.New("failed to load json")

		}
		return placeholder, nil

	}

	if err.Error() != "page not found" {
		fmt.Printf("RequestErr: %s\n", err.Error())
	}

	return placeholder, err

}
