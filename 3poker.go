package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

const nominals_count = 13

type card struct {
	nominal int
	suit    int
	onHand  bool
}

func cardFromCode(code string) (card, error) {
	c := card{}
	if len(code) != 2 {
		return c, errors.New("wrong code")
	}

	nominals := map[string]int{
		"2": 0,
		"3": 1,
		"4": 2,
		"5": 3,
		"6": 4,
		"7": 5,
		"8": 6,
		"9": 7,
		"T": 8,
		"J": 9,
		"Q": 10,
		"K": 11,
		"A": 12,
	}
	n, found := nominals[code[:1]]
	if !found {
		return c, errors.New("wrong nominal")
	}
	c.nominal = n

	suits := map[string]int{
		"S": 0,
		"C": 1,
		"D": 2,
		"H": 3,
	}

	s, found := suits[code[1:]]
	if !found {
		return c, errors.New("wrong suit")
	}

	c.suit = s

	return c, nil
}

type playerHand [2]card

func readDataFile(filePath string) ([][]playerHand, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")
	countGames, err := strconv.Atoi(lines[0])
	if err != nil {
		return nil, err
	}
	//fmt.Println(countGames)

	games := [][]playerHand{}
	lineCount := 0
	for i := 0; i < countGames; i++ {
		lineCount++
		countPlayers, err := strconv.Atoi(lines[lineCount])
		if err != nil {
			return nil, err
		}
		//fmt.Println(countPlayers)
		gamePlayers := []playerHand{}
		for j := 0; j < countPlayers; j++ {
			lineCount++
			//fmt.Println(lines[lineCount])
			playerData := strings.Split(lines[lineCount], " ")
			player := playerHand{}
			for k := 0; k < len(player); k++ {
				player[k], err = cardFromCode(playerData[k])
				if err != nil {
					return nil, err
				}
			}
			gamePlayers = append(gamePlayers, player)
		}
		games = append(games, gamePlayers)
	}

	return games, nil
}

func makeDeck() (deck []card) {
	for suit := 0; suit < 4; suit++ {
		for nom := 0; nom < nominals_count; nom++ {
			deck = append(deck, card{nominal: nom, suit: suit})
		}
	}

	return deck
}

func dealCards(players []playerHand) (deck []card) {
	deck = makeDeck()
	for _, p := range players {
		for cid := range p {
			id := p[cid].nominal + (p[cid].suit * nominals_count)
			deck[id].onHand = true
		}
	}

	return deck
}

func maxValue(val ...int) int {
	max := val[0]
	for i := 1; i < len(val); i++ {
		if val[i] > max {
			max = val[i]
		}
	}

	return max
}

func calcCost(player playerHand, tableCard card) int {
	setDelta := nominals_count
	switch {
	case (player[0].nominal == tableCard.nominal) && (player[1].nominal == tableCard.nominal):
		return 2*setDelta + tableCard.nominal
	case player[0].nominal == player[1].nominal:
		return setDelta + player[0].nominal
	case player[0].nominal == tableCard.nominal:
		return setDelta + player[0].nominal
	case player[1].nominal == tableCard.nominal:
		return setDelta + player[1].nominal
	default:
		return maxValue(player[0].nominal, player[1].nominal, tableCard.nominal)
	}
}

func getWinnaleCards(players []playerHand) (cards []card) {
	deck := dealCards(players)
	for id := range deck {
		if deck[id].onHand {
			continue
		}
		playerPoints := calcCost(players[0], deck[id])
		winnable := true
		for pid := 1; pid < len(players); pid++ {
			oppPoints := calcCost(players[pid], deck[id])
			if oppPoints > playerPoints {
				winnable = false
				break
			}
		}
		if winnable {
			cards = append(cards, deck[id])
		}
	}

	return cards
}

func printRes(res []card) {
	fmt.Println(len(res))
	for _, c := range res {
		nom := ""
		switch {
		case c.nominal < 8:
			nom = strconv.Itoa(c.nominal + 2)
		case c.nominal == 8:
			nom = "T"
		case c.nominal == 9:
			nom = "J"
		case c.nominal == 10:
			nom = "Q"
		case c.nominal == 11:
			nom = "K"
		case c.nominal == 12:
			nom = "A"
		}

		suit := ""
		switch c.suit {
		case 0:
			suit = "S"
		case 1:
			suit = "C"
		case 2:
			suit = "D"
		case 3:
			suit = "H"
		}

		fmt.Println(nom + suit)
	}
}

func main() {
	filePath := "./test/2"
	data, err := readDataFile(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, d := range data {
		cards := getWinnaleCards(d)
		printRes(cards)
	}
}
