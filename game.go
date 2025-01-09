package main

import (
	"go/types"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	priory            uint8
	suit, name, class string
}

type CardToCard struct {
	winC, beatC *Card
}

const delay = time.Second * 1

type TypeMove string

const (
	pickUpC TypeMove = "x"
	passC   TypeMove = "p"
	throwC  TypeMove = "s"
)

const (
	botId = iota
	plId
)

func stackCardsShuffle() []Card {
	stack := []Card{
		{1, "Бубий", "Бубий 6", "6"},
		{2, "Бубий", "Бубий 7", "7"},
		{3, "Бубий", "Бубий 8", "8"},
		{4, "Бубий", "Бубий 9", "9"},
		{5, "Бубий", "Бубий 10", "10"},
		{6, "Бубий", "Валет бубий", "Валет"},
		{7, "Бубий", "Дама бубий", "Дама"},
		{8, "Бубий", "Кароль бубий", "Кароль"},
		{9, "Бубий", "Туз бубий", "Туз"},
		{1, "Пики", "Пики 6", "6"},
		{2, "Пики", "Пики 7", "7"},
		{3, "Пики", "Пики 8", "8"},
		{4, "Пики", "Пики 9", "9"},
		{5, "Пики", "Пики 10", "10"},
		{6, "Пики", "Валет Пики", "Валет"},
		{7, "Пики", "Дама Пики", "Дама"},
		{8, "Пики", "Кароль Пики", "Кароль"},
		{9, "Пики", "Туз Пики", "Туз"},
		{1, "Черви", "Черви 6", "6"},
		{2, "Черви", "Черви 7", "7"},
		{3, "Черви", "Черви 8", "8"},
		{4, "Черви", "Черви 9", "9"},
		{5, "Черви", "Черви 10", "10"},
		{6, "Черви", "Валет Черви", "Валет"},
		{7, "Черви", "Дама Черви", "Дама"},
		{8, "Черви", "Кароль Пики", "Кароль"},
		{9, "Черви", "Туз Черви", "Туз"},
		{1, "Крести", "Крести 6", "6"},
		{2, "Крести", "Крести 7", "7"},
		{3, "Крести", "Крести 8", "8"},
		{4, "Крести", "Крести 9", "9"},
		{5, "Крести", "Крести 10", "10"},
		{6, "Крести", "Валет Крести", "Валет"},
		{7, "Крести", "Дама Крести", "Дама"},
		{8, "Крести", "Кароль Крести", "Кароль"},
		{9, "Крести", "Туз Крести", "Туз"},
	}
	for i := 0; i < len(stack); i++ {
		j := rand.Intn(i + 1)
		stack[i], stack[j] = stack[j], stack[i]

	}

	return stack
}

// distribution of game _cards takes place in turn. The distribution of cards is the same for all players, including the bot.
func distribution(stack []Card, p ...*Player) []Card {
	ac := uint8(6 * len(p))
	fc := func() uint8 {
		var c uint8
		for _, v := range p {
			c += uint8(len(v.cards))
		}
		return c
	}
	key := func(cards map[int]Card) int {
		maxV := 0
		for k, _ := range cards {
			if k >= maxV {
				maxV = k
			}
		}
		return maxV + 1
	}
	for {
		if fc() == ac || len(stack) == 0 {
			return stack
		}
		for _, v := range p {
			if len(v.cards) < 6 {

				v.cards[key(v.cards)] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		}
	}
}

func pickUpCards(cards map[int]Card, stack []CardToCard) {
	key := func() int {
		maxV := 0
		for k, _ := range cards {
			if k >= maxV {
				maxV = k
			}
		}
		return maxV + 1
	}

	for _, v := range stack {
		if v.winC != nil {
			cards[key()] = *v.winC
		}
		if v.beatC != nil {
			cards[key()] = *v.beatC
		}
	}
}
func gameOver(stack []Card, p ...*Player) (int, bool) {
	if len(stack) == 0 {
		for _, pl := range p {
			if len(pl.cards) == 0 {
				return int(pl.id), true
			}
		}
	}
	return -1, false
}

// CallbackThrowCard it is executed on any throwC by any player
type CallbackThrowCard func(card *Card) (error, bool)

// SetBeatC the card that needs to be beaten can be any.If an error is returned, a flowing stack is returned.
func (card Card) SetBeatC(stack []CardToCard) (error, []CardToCard) {
	if len(stack) == 0 {
		return nil, append(stack, CardToCard{beatC: &card})
	}
	for i := len(stack) - 1; i >= 0; i-- {
		if stack[i].winC == nil {
			return types.Error{Msg: "The card is not beaten!"}, stack
		}
		if stack[i].beatC.class == card.class || stack[i].winC.class == card.class {
			return nil, append(stack, CardToCard{beatC: &card})
		}
	}
	return types.Error{Msg: "This card cannot be thrown!"}, stack
}

func (card Card) isWin(beatC Card, suitTrump string) (error, bool) {
	bsCheck := card.suit == beatC.suit && card.priory > beatC.priory

	switch beatC.suit {
	case suitTrump:
		if bsCheck {
			return nil, true
		}

	default:
		switch card.suit {
		case suitTrump:
			return nil, true
		default:
			if bsCheck {
				return nil, true
			}
		}
	}
	return types.Error{Msg: "You can't beat this card."}, false
}

// SetWinC maybe the card is a winning one. Before set, check if it is a winning one.
func (card Card) SetWinC(stack []CardToCard) (error, bool) {
	if len(stack) < 1 {
		return types.Error{Msg: "An empty stack of cards!"}, false
	}
	stack[len(stack)-1].winC = &card
	return nil, true

}
func infoGame(stack []Card, trumpC Card) string {
	return "count: " + strconv.Itoa(len(stack)) + " trump card: " + trumpC.name
}
func infoCurrentBatch(stack []CardToCard) string {
	sn := make([]string, 0, len(stack))
	for _, v := range stack {
		if v.winC == nil {
			sn = append(sn, v.beatC.name)
		} else {
			sn = append(sn, v.winC.name+" -> "+v.beatC.name)
		}
	}
	return "Current batch: " + strings.Join(sn, ",")
}

// start..
func start(stack []Card, pl *Player, bot *Bot, trumpC Card) {
	/*botPl := &bot.pl
		println("Card distribution...")
		time.Sleep(delay / 2)
		stack = distribution(stack, pl, botPl)

		priority := 10
		//The player with the first move
		currIdPl := plId
		for _, c := range [2]*Player{pl, botPl} {
			for _, v := range c.cards {
				if v.suit == trumpC.suit && v.priory < uint8(priority) {
					priority = int(v.priory)
					currIdPl = int(c.id)
				}
			}
		}
		println("The first one goes by card (id):" + strconv.Itoa(currIdPl))

		currStackC := make([]CardToCard, 0, 6)
		//parties
	/*	for {
			currStackC = currStackC[:0]
			println(infoGame(stack, trumpC))
			time.Sleep(delay / 2)
			println(infoPlayer(pl))
			switch currIdPl {
			case plId:
				endBatch := func(stack []CardToCard, iThrow int) bool {
					println("Now the bot is running...")
					time.Sleep(delay * 2)

					cbBotThrow := func(card *Card) (error, bool) {
						if card == nil {
							pickUpCards(bot.pl.cards, stack)
							currIdPl = plId
						} else {
							stack[iThrow].winC = card
							currIdPl = botId
						}

						return nil, true
					}
					bot.RespondC(*stack[iThrow].beatC, trumpC.suit, cbBotThrow)
					return
				}
				pl.StartBatch(currStackC, endBatch)
				panic("Finish game..")
			case botId:
				return
			}
			return

		} */
}
