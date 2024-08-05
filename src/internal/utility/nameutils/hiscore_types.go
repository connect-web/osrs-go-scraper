package nameutils

type HiscoreType struct {
	Skill     string `json:"skill"`
	Minigames string `json:"minigames"`
	Id        int    `json:"id"`
	Limit     int    `json:"limit"`
}
