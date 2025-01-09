package main

import "time"

type Mode uint8

const (
	lite Mode = iota
	mid
	hard
)

type Bot struct {
	pl   Player
	mode Mode
}

// RespondC Hits with the smallest card if possible. Otherwise , it takes away
func (bot *Bot) RespondC(beatC Card, suitTrump string, throw CallbackThrowCard) {
	if len(bot.pl.cards) == 0 {
		panic("Bot is deck of cards is empty!")
	}
	key := bot.FindMinC(suitTrump, func(c Card) bool {
		_, ok := c.isWin(beatC, suitTrump)
		return ok == true
	})

	minC, ok := bot.pl.cards[key]

	if ok {
		err, _ := throw(&minC)
		if err == nil {
			delete(bot.pl.cards, key)
			time.Sleep(delay * 2)
			println("Bot throws the card: " + minC.name)
			return
		}
		panic(err.Error())
	}
	time.Sleep(delay * 2)
	println("bot takes the cards...")
	err, _ := throw(nil)
	if err != nil {
		panic(err.Error())
	}

}

// FindMinC The smallest card including the trump card
func (bot *Bot) FindMinC(trump string, conditions ...func(c Card) bool) int {
	kTrump := make([]int, 0, 3)
	kDef := make([]int, 0, 3)

	appC := func(key int) {
		switch bot.pl.cards[key].suit {
		case trump:
			kTrump = append(kTrump, key)
		default:
			kDef = append(kDef, key)
		}
	}
	chCond := func(key int) bool {
		for _, f := range conditions {
			if !f(bot.pl.cards[key]) {
				return false
			}
		}
		return true
	}
	for k, _ := range bot.pl.cards {
		if chCond(k) {
			appC(k)
		}
	}
	//the smallest game card
	var minCk = -1
	for i, sl := range [2][]int{kDef, kTrump} {
		var priory uint8 = 10
		for iv, v := range sl {
			if iv == 0 && i == 1 && minCk != -1 {
				return minCk
			}
			if bot.pl.cards[v].priory < priory {
				priory = bot.pl.cards[v].priory
				minCk = v
			}
		}
	}
	return minCk
}
