package words

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type randUnit struct {
	nounIndex, adjectiveIndex int
}

var adjectives, nouns []string
var wordChan chan randUnit

func init() {
	nouns = read("words/nouns.txt")
	adjectives = read("words/adjectives.txt")
	wordChan = make(chan randUnit, 2)
	go generate()
}

func AwesomeUsername() string {
	unit := <-wordChan
	return adjectives[unit.adjectiveIndex] + " " + nouns[unit.nounIndex]
}

func read(what string) []string {
	var words []string
	file, err := os.Open(what)
	if err != nil {
		log.Panicln("[words.read] error opening file", what, ":", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Println("[words.read] error closing file", what, ":", err)
		}
	}()
	sc := bufio.NewScanner(file)
	sc.Split(bufio.ScanLines)
	for sc.Scan() {
		word := strings.Title(strings.TrimSpace(sc.Text()))
		if len(word) != 0 {
			words = append(words, word)
		}
	}
	return words
}

func generate() {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	lenNouns := len(nouns)
	lenAdjectives := len(adjectives)
	for {
		wordChan <- randUnit{
			nounIndex:      r.Intn(lenNouns),
			adjectiveIndex: r.Intn(lenAdjectives),
		}
	}
}
