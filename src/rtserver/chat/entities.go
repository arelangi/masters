package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	sentences "gopkg.in/neurosnap/sentences.v1"
)

var (
	wordTokenizer *sentences.DefaultWordTokenizer
)

func init() {
	punctStrings := sentences.NewPunctStrings()
	wordTokenizer = sentences.NewWordTokenizer(punctStrings)
}

func getTokenStrings(tokens []*sentences.Token) (out []string) {
	for _, token := range tokens {
		out = append(out, token.Tok)
	}
	return
}

func cleanNormalizedTweet(input string) string {
	usernameExp := regexp.MustCompile(`\[username\]`)
	linkExp := regexp.MustCompile(`\[link\]`)
	emoticonExp := regexp.MustCompile(`\[e\]`)
	hashtagExp := regexp.MustCompile(`\[hashtag\]`)
	dotExp := regexp.MustCompile(`\.*`)

	printString := usernameExp.ReplaceAllString(input, ``)
	printString = linkExp.ReplaceAllString(printString, ``)
	printString = emoticonExp.ReplaceAllString(printString, ``)
	printString = hashtagExp.ReplaceAllString(printString, ``)
	printString = dotExp.ReplaceAllString(printString, ``)
	printString = strings.TrimSpace(printString)

	return printString
}

func extractEntities(input string) (out []string) {
	cleanedInput := cleanNormalizedTweet(input)
	tokens := getTokenStrings(wordTokenizer.Tokenize(cleanedInput, false))

	es, err := ext.Extract(tokens)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cleanedInput)
	for _, v := range es {
		out = append(out, fmt.Sprintf("%s: %s", ext.Tags()[v.Tag], v.Name))
	}
	return
}
