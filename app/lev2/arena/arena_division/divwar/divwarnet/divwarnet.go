package div_war_net

import (
	"fmt"
	. "wartank/app/lev0/types"
	"wartank/app/lev2/arena/arena_net"
)

/*
	Автоматически воюет в сражении
*/

// DivWarNet -- танкует в сражении
type DivWarNet struct {
	*arena_net.АренаСеть
	bot ИБот
}

// NewDivWarNet -- возвращает новый *DivWarNet
func NewDivWarNet(bot ИБот) (*DivWarNet, error) {
	if bot == nil {
		return nil, fmt.Errorf("NewDivWarNet(): IServBpt == nil")
	}

	сам := &DivWarNet{
		// SectionNet: section_net.NewSectionNet(server, bot, ..., "https://wartank.ru/bitva"),
		bot: bot,
	}
	return сам, nil
}
