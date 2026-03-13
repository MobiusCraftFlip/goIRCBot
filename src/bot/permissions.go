package bot

import (
	"fmt"
	"strings"

	"github.com/lrstanley/girc"
)

func (b *Bot) isModerator(event girc.Event) bool {
	for _, channel := range b.IRC().Channels() {
		fmt.Println("A")
		if channel.Name != b.Config().AdminChan {
			fmt.Println("B")
			continue
		}
		fmt.Println("C")
		fmt.Println(strings.ToLower(event.Source.Name), channel.UserList)
		if channel.UserIn(event.Source.Name) {
			fmt.Println("D")
			return true
		}
	}

	return false
}
