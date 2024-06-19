package entities

import (
	"fmt"
	utils "github.com/connect-web/Low-Latency-Utils"
	"names.go/entities/limits"
	"strconv"
)

func GetNames(info limits.HiscoreType, page int, iterator *utils.ProxyIterator) (default_response map[string]bool, err error) {
	parameters := map[string]string{
		"table": strconv.Itoa(info.Id),
		"page":  strconv.Itoa(page),
	}
	if info.Minigames != "" {
		parameters["category_type"] = "1"
	}

	response, err := utils.Request("https://secure.runescape.com/m=hiscore_oldschool/overall", parameters, iterator, -999)

	if err == nil {
		page_content := NewDataFormat(page)
		name_map := page_content.GetNames(response)
		return name_map, nil
	} else {
		fmt.Println(err.Error())
		return default_response, err
	}
}
