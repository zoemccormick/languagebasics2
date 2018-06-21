package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

type proverb struct {
	line  string
	chars map[rune]int
}

func (p *proverb) charCount() map[rune]int {
	if p.chars != nil {
		return p.chars
	}

	m := make(map[rune]int, 0)
	for _, c := range p.line {
		m[c] = m[c] + 1
	}
	p.chars = m
	return p.chars
}

func main() {
	path := pathFromFlag()
	if path == "" {
		path = pathFromEnv()
	}

	if path == "" {
		fmt.Println("You must specify one the file path with -f or as FILE environment variable.")
		os.Exit(1)
	}

	proverbs, err := loadProverbs(path)
	if err != nil {
		fmt.Printf("Failed to load proverbs: %s", err)
		os.Exit(1)
	}

	mu := &sync.Mutex{}

	for _, p := range proverbs {
		go func(p *proverb) {
			mu.Lock()
			fmt.Printf("%s\n", p.line)
			for k, v := range p.charCount() {
				fmt.Printf("'%c'=%d, ", k, v)
			}
			fmt.Print("\n\n")
			mu.Unlock()
		}(p)
	}

	//this line is necessary to make sure that the program doesn't finish before
	//the goroutine does - ie. the main method finishes & closes before output from
	//go function is printed
	time.Sleep(5 * time.Second)

}

func pathFromFlag() string {
	path := flag.String("f", "", "file flag")
	flag.Parse()
	return *path
}

func pathFromEnv() string {
	return os.Getenv("FILE")
}

func loadProverbs(path string) ([]*proverb, error) {
	var proverbs []*proverb

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bs), "\n")
	for _, line := range lines {
		p := &proverb{line: line}
		proverbs = append(proverbs, p)
	}

	return proverbs, nil
}
