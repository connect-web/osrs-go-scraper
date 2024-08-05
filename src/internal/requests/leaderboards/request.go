package leaderboards

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/requests"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/nameutils"
	"strconv"
)

func GetNames(info nameutils.HiscoreType, page int, iterator *entities.ProxyIterator) (defaultResponse map[string]bool, err error) {
	parameters := map[string]string{
		"table": strconv.Itoa(info.Id),
		"page":  strconv.Itoa(page),
	}
	if info.Minigames != "" {
		parameters["category_type"] = "1"
	}

	response, err := requests.Request("https://secure.runescape.com/m=hiscore_oldschool/overall", parameters, iterator, -999)

	if err == nil {
		pageContent := nameutils.NewHtmlHiscore(page)
		nameMap := pageContent.GetUsernamesFromHtmlString(response)
		return nameMap, nil
	} else {
		fmt.Println(err.Error())
		return defaultResponse, err
	}
}
