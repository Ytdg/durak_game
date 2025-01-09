package main

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	id    uint8
	cards map[int]Card
}

func g(d []CardToCard) {
	d[0].winC = nil
}
func d(t ...int) {

	println(t[0])
}
func main() {
	/*d := []CardToCard{
		{beatC: &Card{name: "Валет"}},
	}
	println(infoCurrentBatch(d))
	d[0].winC = &Card{name: "Король"}
	println(infoCurrentBatch(d))
	d = append(d, CardToCard{beatC: &Card{name: "Дама"}})
	println(infoCurrentBatch(d))

	/*d := Card{suit: "3", class: "r", name: "3r", priory: 1}
	fmt.Printf("%v", reflect.ValueOf(d).Field().)*/
	println("Shuffle...")
	time.Sleep(delay)
	stack := stackCardsShuffle()
	pl := Player{cards: map[int]Card{}, id: plId}
	bot := Bot{pl: Player{cards: map[int]Card{}, id: botId}, mode: lite}
	start(stack, &pl, &bot, stack[0])
}

func infoPlayer(player *Player) string {
	sn := make([]string, 0, 6)
	for k, v := range player.cards {
		sn = append(sn, "("+strconv.Itoa(k)+")"+v.name)
	}
	return "Your id:" + strconv.Itoa(int(player.id)) + " cards:" + strings.Join(sn, ",")
}

// StartBatch  party when the player started
func (pl *Player) StartBatch(currStackC []CardToCard, endBatch func(curSt []CardToCard, iThrow int) bool) {
	cTrow := 0
	throw := func() bool {
		pl.SelectCard(func(card *Card) (error, bool) {
			err, v := card.SetBeatC(currStackC)
			currStackC = v
			return err, err == nil
		})
		cTrow++
		println(infoCurrentBatch(currStackC) + "\n" + infoPlayer(pl))
		nextM := endBatch(currStackC, len(currStackC)-1)
		println(infoCurrentBatch(currStackC) + "\n" + infoPlayer(pl))
		return nextM
	}
	for {
		if cTrow > 0 {
			for {
				var value TypeMove
				print("\nInput 'p'(pass)/'s'(throw).(You can't change the selection.):")
				_, _ = fmt.Scan(&value)

				switch value {
				case passC:
					return
				case throwC:
					if !throw() {
						return
					}
				default:
					fmt.Printf("\n%v", types.Error{
						Msg: "Incorrect value!",
					})
				}
			}
		} else {
			if !throw() {
				return
			}
		}
	}
}

// SelectCard processing of card number input
func (pl *Player) SelectCard(throw CallbackThrowCard) {
	var value string

	for {
		print("Input num card:")
		_, errInp := fmt.Scan(&value)
		if num, errNum := strconv.Atoi(value); errNum == nil && errInp == nil {
			if card, ok := pl.cards[num]; ok {
				err, _ := throw(&card)
				if err == nil {
					delete(pl.cards, num)
					return
				}
				fmt.Printf("\n%v", err)
			}
		}

		println("\nTry again...")
	}
}
