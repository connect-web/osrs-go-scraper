package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/entities"
	"github.com/connect-web/Low-Latency/internal/utils/stats"
	"net/url"
)

func Get_player_stats(name string, player_id int, iterator *entities.ProxyIterator) (stats.SimplePlayer, error) {
	player := *stats.NewSimplePlayer(player_id)

	url_encoded_name := url.QueryEscape(name)
	params := map[string]string{
		"player": url_encoded_name,
	}

	response, err := entities.Request("https://secure.runescape.com/m=hiscore_oldschool/index_lite.json", params, iterator, -99)

	if err == nil {
		var responseJson stats.JsonResponse
		unmarshErr := json.Unmarshal([]byte(response), &responseJson)
		if unmarshErr == nil {
			if player.LoadJson(responseJson) {
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
