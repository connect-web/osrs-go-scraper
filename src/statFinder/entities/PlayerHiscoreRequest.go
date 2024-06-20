package entities

import (
	"encoding/json"
	"errors"
	"fmt"
	utils "github.com/connect-web/Low-Latency-Utils"
	"net/url"
)

func Get_player_stats(name string, player_id int, iterator *utils.ProxyIterator) (SimplePlayer, error) {
	player := *NewSimplePlayer(player_id)

	url_encoded_name := url.QueryEscape(name)
	params := map[string]string{
		"player": url_encoded_name,
	}

	response, err := utils.Request("https://secure.runescape.com/m=hiscore_oldschool/index_lite.json", params, iterator, -99)

	if err == nil {
		var responseJson jsonResponse
		unmarshErr := json.Unmarshal([]byte(response), &responseJson)
		if unmarshErr == nil {
			if player.loadJson(responseJson) {
				player.Username = name
				player.Calculations()
				return player, nil
			}
		} else {
			fmt.Printf(unmarshErr.Error())
		}

		return player, errors.New("Failed to load json")

	} else {
		fmt.Println(err.Error())
		return player, err
	}
}
