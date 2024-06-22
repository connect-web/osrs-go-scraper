package main

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/entities"
	"github.com/connect-web/Low-Latency/internal/utils/name"
	"strconv"
)

func GetNames(info name.HiscoreType, page int, iterator *entities.ProxyIterator) (default_response map[string]bool, err error) {
	parameters := map[string]string{
		"table": strconv.Itoa(info.Id),
		"page":  strconv.Itoa(page),
	}
	if info.Minigames != "" {
		parameters["category_type"] = "1"
	}

	response, err := entities.Request("https://secure.runescape.com/m=hiscore_oldschool/overall", parameters, iterator, -999)

	if err == nil {
		page_content := name.NewDataFormat(page)
		name_map := page_content.GetNames(response)
		return name_map, nil
	} else {
		fmt.Println(err.Error())
		return default_response, err
	}
}
