package internal

import (
	"time"

	"anti_bruteforce/config"
)

type LimitItem struct {
	Burst  int
	Period time.Duration
}

type LimitSettings struct {
	Login    LimitItem
	Password LimitItem
	IP       LimitItem
}

func FromConfig(config config.Config) *LimitSettings {
	return &LimitSettings{
		Login: LimitItem{
			Burst:  config.LimitBurstLogin,
			Period: config.LimitPeriodLogin,
		},
		Password: LimitItem{
			Burst:  config.LimitBurstPassword,
			Period: config.LimitPeriodPassword,
		},
		IP: LimitItem{
			Burst:  config.LimitBurstIP,
			Period: config.LimitPeriodIP,
		},
	}
}
