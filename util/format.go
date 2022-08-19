package util

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
)

func FullTag(u discord.User) string {
	return fmt.Sprintf("%s (`%s` | `%d`)", u.Mention(), u.Tag(), u.ID)
}
