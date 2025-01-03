package nameutils

type PageLimitInfo struct {
	Limit       int
	HiscoreType HiscoreType
}

type PageLimits struct {
	Limits []HiscoreType
}

func NewPageLimitManager() *PageLimits {
	pageLimits := &PageLimits{
		Limits: []HiscoreType{},
	}
	pageLimits.Reset()
	return pageLimits
}
func (p *PageLimits) Add(hiscoreLimit HiscoreType) {
	p.Limits = append(p.Limits, hiscoreLimit)
}

func (p *PageLimits) Reset() {
	// All skills start off with known limits of 80k Pages.
	p.Limits = []HiscoreType{
		{Skill: "Overall", Id: 0, Limit: 80_000},
		{Skill: "Attack", Id: 1, Limit: 80_000},
		{Skill: "Defence", Id: 2, Limit: 80_000},
		{Skill: "Strength", Id: 3, Limit: 80_000},
		{Skill: "Hitpoints", Id: 4, Limit: 80_000},
		{Skill: "Ranged", Id: 5, Limit: 80_000},
		{Skill: "Prayer", Id: 6, Limit: 80_000},
		{Skill: "Magic", Id: 7, Limit: 80_000},
		{Skill: "Cooking", Id: 8, Limit: 80_000},
		{Skill: "Woodcutting", Id: 9, Limit: 80_000},
		{Skill: "Fletching", Id: 10, Limit: 80_000},
		{Skill: "Fishing", Id: 11, Limit: 80_000},
		{Skill: "Firemaking", Id: 12, Limit: 80_000},
		{Skill: "Crafting", Id: 13, Limit: 80_000},
		{Skill: "Smithing", Id: 14, Limit: 80_000},
		{Skill: "Mining", Id: 15, Limit: 80_000},
		{Skill: "Herblore", Id: 16, Limit: 80_000},
		{Skill: "Agility", Id: 17, Limit: 80_000},
		{Skill: "Thieving", Id: 18, Limit: 80_000},
		{Skill: "Slayer", Id: 19, Limit: 80_000},
		{Skill: "Farming", Id: 20, Limit: 80_000},
		{Skill: "Runecraft", Id: 21, Limit: 80_000},
		{Skill: "Hunter", Id: 22, Limit: 80_000},
		{Skill: "Construction", Id: 23, Limit: 80_000},
	}
}

func (p *PageLimits) Reverse() {
	length := len(p.Limits)
	for i := 0; i < length/2; i++ {
		p.Limits[i], p.Limits[length-i-1] = p.Limits[length-i-1], p.Limits[i]
	}
}
