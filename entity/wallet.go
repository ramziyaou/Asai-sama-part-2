package entity

import (
	"sync"
	"time"
)

// Cryptowallet object
type CryptoWallet struct {
	Username string
	Name     string `gorm:"primaryKey"`
	Amount   int64  ///////////////////////hhh`gorm:"default:18"`
	// Stop       chan struct{} `sql:"-"`
	// Stop       bool //sql.NullBool  then when assigning sql.NullBool{}
	// Notstarted bool ////////////////////hhh`gorm:"default:true"`
	sync.RWMutex
}

type StartStopCheck struct {
	Username string
	Name     string
	Stop     bool
	Start    bool
}

func NewStartStop(username, name string) *StartStopCheck {
	return &StartStopCheck{
		username,
		name,
		false,
		false,
	}
}

func NewWallet(name string) *CryptoWallet {
	return &CryptoWallet{
		"",
		name,
		0,
		// make(chan struct{}),
		// false,
		// true,
		sync.RWMutex{},
	}
}

func (c *CryptoWallet) Mine() {
	time.Sleep(10 * time.Second)
	c.Lock()
	c.Amount++
	c.Unlock()
}
