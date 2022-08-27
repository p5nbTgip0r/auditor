package util

type Region struct {
	ID    string
	Name  string
	Emoji string
	VIP   bool
}

var (
	Regions = map[string]Region{}
)

func generateRegion(id string, name string, emoji string, vip bool) {
	Regions[id] = Region{
		ID:    id,
		Name:  name,
		Emoji: emoji,
		VIP:   vip,
	}
}

func init() {
	generateRegion("brazil", "Brazil", ":flag_br:", false)
	generateRegion("hongkong", "Hong Kong", ":flag_hk:", false)
	generateRegion("india", "India", ":flag_in:", false)
	generateRegion("japan", "Japan", ":flag_jp:", false)
	generateRegion("milan", "Milan", ":flag_it:", false)
	generateRegion("rotterdam", "Rotterdam", ":flag_nl:", false)
	generateRegion("russia", "Russia", ":flag_ru:", false)
	generateRegion("singapore", "Singapore", ":flag_sg:", false)
	generateRegion("southafrica", "South Africa", ":flag_za:", false)
	generateRegion("south-korea", "South Korea", ":flag_kr:", false)
	generateRegion("sydney", "Sydney", ":flag_au:", false)
	generateRegion("us-central", "US Central", ":flag_us:", false)
	generateRegion("us-east", "US East", ":flag_us:", false)
	generateRegion("us-south", "US South", ":flag_us:", false)
	generateRegion("us-west", "US West", ":flag_us:", false)

	generateRegion("automatic", "Automatic", "", false)
}
