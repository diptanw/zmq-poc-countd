package main

import (
	"bufio"
	"flag"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	zmq "github.com/pebbe/zmq4"
)

var (
	countdUrl string
	debug     bool
	rate      float64
	wordLimit int
)

type Words []string

func (words Words) Sample() string {
	return words[rand.Intn(len(words)-1)]
}

func normalizeWord(w string) string {
	const trimList = "!., "
	w = strings.ToLower(w)
	return strings.Trim(w, trimList)
}

func ReadWordsFromStdin(limit int) Words {
	words := make(Words, 0)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		if len(words) < limit {
			w := normalizeWord(scanner.Text())
			words = append(words, w)
			if len(words) == limit {
				log.Printf("Hit limit %d, discarding further words...")
			}
		}
	}

	log.Printf("Read %d words", len(words))
	return words
}

func main() {
	flag.StringVar(&countdUrl, "countd-url", "tcp://localhost:6500", "")
	flag.BoolVar(&debug, "debug", false, "Write emitted words to stdout instead of the ZMQ socket")
	flag.Float64Var(&rate, "rate", 1, "Rate parameter for emission of words. This is the lambda value for an exponential distribution.")
	flag.IntVar(&wordLimit, "word-limit", 10000, "Maximium no. of words to read from stdin.")
	flag.Parse()

	zOutput, _ := zmq.NewSocket(zmq.PUSH)
	if !debug {
		log.Printf("Connecting to %s", countdUrl)
		defer zOutput.Close()
		if err := zOutput.Connect(countdUrl); err != nil {
			log.Fatal(err)
		}
	}

	rand.Seed(time.Now().UnixNano())

	words := ReadWordsFromStdin(wordLimit)

	for {
		word := words.Sample()

		if debug {
			log.Printf("%s", word)
		} else {
			if _, err := zOutput.SendMessage(word); err != nil {
				log.Fatal(err)
			}
		}

		d := time.Millisecond * time.Duration(1000*rand.ExpFloat64()/rate)
		time.Sleep(d)
	}
}
