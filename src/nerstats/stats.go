package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/sbl/ner"
	"gitlab.logzero.in/arelangi/mlog"
)

var streamChan chan []byte
var ext *ner.Extractor

func init() {
	streamChan = make(chan []byte)
	var err error
	ext, err = ner.NewExtractor("/repos/MITIE/MITIE-models/english/ner_model.dat")
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	data := "/repos/mldata/ner100.csv"
	csv_file, _ := os.Open(data)
	r := csv.NewReader(csv_file)

	count := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record[0])
		r, err := handleTweet(record[0])
		if err != nil {
			fmt.Println("*****************")
			continue
		}
		fmt.Println("-------------------------")
		fmt.Println(r.Predictions.LstmPrediction)
		fmt.Println(r.Predictions.TranslatedText)
		fmt.Println(r.Entities)
		count++

	}
}

func handleTweet(input string) (resp Response, err error) {
	var predictions Predictions
	var mitieVals, proseVals []string
	msg := RTTweet{OriginalText: input, NormalizedText: input}

	predictions, err = callMLEngine(msg)
	if err != nil {
		return
	}

	if predictions.TranslatedText == "" {
		mitieVals = mitieEntities(msg.NormalizedText)
		proseVals = proseEntities(msg.NormalizedText)
	} else {
		mitieVals = mitieEntities(predictions.TranslatedText)
		proseVals = proseEntities(predictions.TranslatedText)
	}

	entities := make(map[string][]string)
	entities["mitie"] = mitieVals
	entities["prose"] = proseVals
	entities["spacy"] = predictions.SpacyEntities
	entities["nltk"] = predictions.NLTKEntities

	resp = Response{Tweet: msg, Predictions: predictions, Entities: entities}
	return

}

func callMLEngine(request RTTweet) (predictions Predictions, err error) {
	url := "http://localhost:7070/predict"
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(request)

	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest("POST", url, b)
	if err != nil {
		mlog.Error("Creating POST Request to ML Engine failed", mlog.Items{"error": err})
		return
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		mlog.Error("Request to ML Engine failed", mlog.Items{"error": err})
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&predictions)
	if err != nil {
		//		mlog.Error("JSON Decoding error", mlog.Items{"error": err})
		return
	}

	//Now, make a call to the NER libraries

	return
}

type RTTweet struct {
	OriginalText   string `json:"original_text"`
	NormalizedText string `json:"normalized_text"`
}

type Predictions struct {
	FtPrediction      string   `json:"ft_prediction"`
	LstmPrediction    string   `json:"lstm_prediction"`
	NbPrediction      string   `json:"nb_prediction"`
	SvmPrediction     string   `json:"svm_prediction"`
	TranslatedText    string   `json:"translated_text"`
	NonEnglishPercent float64  `json:"non_english_percent"`
	SpacyEntities     []string `json:"spacy_entities"`
	NLTKEntities      []string `json:"nltk_entities"`
}

type Response struct {
	Tweet       RTTweet             `json:"tweet"`
	Predictions Predictions         `json:"predictions"`
	Entities    map[string][]string `json:"entities"`
}

func normalize(tweet string) string {
	//var emojiRx = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[💪]|[❤️]|[🔥]|[🤗]|[👏]|[💞]|[🌻]|[🌺]|[🕊💐|[💯]|[🤘]|[👌]|[💛]|[🥁]|[🎹]|[🤙]|[🤩]|[👇]|[🌟]|[💜]|[💚]|[🥰]|[💦]|[🤦]|[👈]|[🕺]|[💃]|[💖]|[💕]|[✊]|[🏏]|[🎼]|[🗿]|[🎉]|[📸]|[🎊]|[🗣]|[💙]|[👍]|[💥]| [🤔]| [🥳] |[💟]|[🌪]|[🥛] | [💋] | [💏]| [🏃]|[⭐]|[👊]|[🤛]|[💥]|[🎂]|[✌🏻]|[🤣]|[👉]|[🤕]|[💓]|[🤝]|[🇸🇻]|[✅]|[🤮]|[🍬]|[🌹]|[✍🏻]|[✨]|[🎬]|[❓]|[👐]|[🍺]|[🗝]|[👃]|[🖕]|[🔀]|[]🏽]|[📢]|[]🏾]|[🤬]|[🤟]|[👅]|[🌼]|[🤪]|[🤷]|[🌸]|[🧡]|[🤧]|[💝]|[🤫]|[]🏼]|[💫]|[💗]|[🍌]|[🍼]|[👸]|[🍻]|[💰]|[💵]|[🤯]|[🤐]|[📣]|[🌃]|[💑]|[🌃]|[👑]|[🤭]|[🎤]|[🖌]|[🦁]|[👀]|[🤡]|[🌬]|[🤞]|[🤒]|[🏥]|[👩]|[👨]|[💉]|[💊]|[👨‍👨‍👧‍👦]|[🍉]|[🇲🇾]|[💔]|[👫]|[🥺]|[🎶]|[▶]|[•]|[💘]|[🤑]|[🥳]|[🇮🇳]|[🤤]|[🤔]|[🏃]|[🎻]|[🎸]|[📯]|[🎺]|[🎵]|[🌞]|[🥛]|[💎]|[🌄]|[🤓]|[➡]|[✖]|[►]|[❣]|[👋]|[📷]|[💨]|[🔟]|[🎈]|[🤴]|[🎯]|[🌷]|[🌷]|[👪]|[ 💤]|[🌙]|[🌠]|[🤢]|[👄]|[💻]|[💿]|[]🏿]|[👽]|[🌋]|[❌]|[🐉]|[💆]|[👼]|[🍀]|[⏳]|[🎁]|[🎇]|[🏆]|[🥂]|[🤥]|[💁]|[🔯]|[👎]|[🎷]|[🥵]|[🕯]|[👖]|[👱]|[👥]|[🌧]|[👗]|[👬]|[🎆]|[🆚]|[🏄]|[🔊]|[🎙]|[✋]|[🔔]|[🔯]|[👁]|[🏆]|[🍫]|[🥃]|[🏠]|[🍨]|[🍧]|[🥧]|[🍰]|[🐒]|[🐿]|[👼]|[🇧🇩]|[🐕]|[👿]|[🌊]|[🏪]|[🌿]|[🌱]|[🖐]|[🤜]|[🎋]|[🍾]|[🍷]|[🤳]|[🍴]|[🗑]|[🥴]|[👆]|[🔇]|[🐑]|[🇦🇺]|[🐵]|[🐶]|[🦋]|[🔞]|[🔃]|[💸]|[🤸]|[🍭]|[🌊]|[📦]|[🍷]|[🍸]|[🍹]|[🔱]|[🔄]|[🎽]|[🏀]|[💀]|[🎃]|[🏡]|[🐖]|[🐷]|[🐯]|[🐍]|[🌲]|[🐝]|[💧]|[🏋]|[🤺]|[🖐]|[👶]|[🌤]|[🅰]|[🔞]|[🤚]|[🐻]|[💋]|[📽]|[🏭]|[🌁]|[🏗]|[🏪]|[🏩]|[🇹🇹]|[🔫]|[💭]|[🍆]|[🍦]|[🔸]|[📌]|[🔨]|[🥚]|[🤲]|[💭]|[🥇]|[🤨]|[🏟]`)
	//var emojiRx = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[💪]|[❤️]|[🔥]|[🤗]|[👏]|[💞]|[🌻]|[🌺]|[🕊💐|[💯]|[🤘]|[👌]|[💛]|[🥁]|[🎹]|[🤙]|[🤩]|[👇]|[🌟]|[💜]|[💚]|[🥰]|[💦]|[🤦]|[👈]|[🕺]|[💃]|[💖]|[💕]|[✊]|[🏏]|[🎼]|[🗿]|[🎉]|[📸]|[🎊]|[🗣]|[💙]|[👍]|[💥]| [🤔]| [🥳] |[💟]|[🌪]|[🥛] | [💋] | [💏]| [🏃]|[⭐]|[👊]|[🤛]|[💥]|[🎂]|[✌🏻]|[🤣]|[👉]|[🤕]|[💓]|[🤝]|[🇸🇻]|[✅]|[🤮]|[🍬]|[🌹]|[✍🏻]|[✨]|[🎬]|[❓]|[👐]|[🍺]|[🗝]|[👃]|[🖕]|[🔀]|[]🏽]|[📢]|[]🏾]|[🤬]|[🤟]|[👅]|[🌼]|[🤪]|[🤷]|[🌸]|[🧡]|[🤧]|[💝]|[🤫]|[]🏼]|[💫]|[💗]|[🍌]|[🍼]|[👸]|[🍻]|[💰]|[💵]|[🤯]|[🤐]|[📣]|[🌃]|[💑]|[🌃]|[👑]|[🤭]|[🎤]|[🖌]|[🦁]|[👀]|[🤡]|[🌬]|[🤞]|[🤒]|[🏥]|[👩]|[👨]|[💉]|[💊]|[👨‍👨‍👧‍👦]|[🍉]|[🇲🇾]|[💔]|[👫]|[🥺]|[🎶]|[▶]|[•]|[💘]|[🤑]|[🥳]|[🇮🇳]|[🤤]|[🤔]|[🏃]|[🎻]|[🎸]|[📯]|[🎺]|[🎵]|[🌞]|[🥛]|[💎]|[🌄]|[🤓]|[➡]|[✖]|[►]|[❣]|[👋]|[📷]|[💨]|[🔟]|[🎈]|[🤴]|[🎯]|[🌷]|[🌷]|[👪]|[ 💤]|[🌙]|[🌠]|[🤢]|[👄]|[💻]|[💿]|[]🏿]|[👽]|[🌋]|[❌]|[🐉]|[💆]|[👼]|[🍀]|[⏳]|[🎁]|[🎇]|[🏆]|[🥂]|[🤥]|[💁]|[🔯]|[👎]|[🎷]|[🥵]|[🕯]|[👖]|[👱]|[👥]|[🌧]|[👗]|[👬]|[🎆]|[🆚]|[🏄]|[🔊]|[🎙]|[✋]|[🔔]|[🔯]|[👁]|[🏆]|[🍫]|[🥃]|[🏠]|[🍨]|[🍧]|[🥧]|[🍰]|[🐒]|[🐿]|[👼]|[🇧🇩]|[🐕]|[👿]|[🌊]|[🏪]|[🌿]|[🌱]|[🖐]|[🤜]|[🎋]|[🍾]|[🍷]|[🤳]|[🍴]|[🗑]|[🥴]|[👆]|[🔇]|[🐑]|[🇦🇺]|[🐵]|[🐶]|[🦋]|[🔞]|[🔃]|[💸]|[🤸]|[🍭]|[🌊]|[📦]|[🍷]|[🍸]|[🍹]|[🔱]|[🔄]|[🎽]|[🏀]|[💀]|[🎃]|[🏡]|[🐖]|[🐷]|[🐯]|[🐍]|[🌲]|[🐝]|[💧]|[🏋]|[🤺]|[🖐]|[👶]|[🌤]|[🅰]|[🔞]|[🤚]|[🐻]|[💋]|[📽]|[🏭]|[🌁]|[🏗]|[🏪]|[🏩]|[🇹🇹]|[🔫]|[💭]|[🍆]|[🍦]|[🔸]|[📌]|[🔨]|[🥚]|[🤲]|[💭]|[🥇]|[🤨]|[🏟]`)
	var emojiRx = regexp.MustCompile(`[\x{1F600}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[🏃]|[💋]|[💏]|[🤔]|[🥳]|[💤]|[]🏼]|[]🏽]|[]🏾]|[]🏿]|[•]|[⏳]|[▶]|[►]|[✅]|[✊]|[✋]|[✌🏻]|[✍🏻]|[✖]|[✨]|[❌]|[❓]|[❣]|[❤️]|[➡]|[⭐]|[🅰]|[🆚]|[🇦🇺]|[🇧🇩]|[🇮🇳]|[🇲🇾]|[🇸🇻]|[🇹🇹]|[🌁]|[🌃]|[🌃]|[🌄]|[🌊]|[🌊]|[🌋]|[🌙]|[🌞]|[🌟]|[🌠]|[🌤]|[🌧]|[🌪]|[🌬]|[🌱]|[🌲]|[🌷]|[🌷]|[🌸]|[🌹]|[🌺]|[🌻]|[🌼]|[🌿]|[🍀]|[🍆]|[🍉]|[🍌]|[🍦]|[🍧]|[🍨]|[🍫]|[🍬]|[🍭]|[🍰]|[🍴]|[🍷]|[🍷]|[🍸]|[🍹]|[🍺]|[🍻]|[🍼]|[🍾]|[🎁]|[🎂]|[🎃]|[🎆]|[🎇]|[🎈]|[🎉]|[🎊]|[🎋]|[🎙]|[🎤]|[🎬]|[🎯]|[🎵]|[🎶]|[🎷]|[🎸]|[🎹]|[🎺]|[🎻]|[🎼]|[🎽]|[🏀]|[🏃]|[🏄]|[🏆]|[🏆]|[🏋]|[🏏]|[🏗]|[🏟]|[🏠]|[🏡]|[🏥]|[🏩]|[🏪]|[🏪]|[🏭]|[🐉]|[🐍]|[🐑]|[🐒]|[🐕]|[🐖]|[🐝]|[🐯]|[🐵]|[🐶]|[🐷]|[🐻]|[🐿]|[👀]|[👁]|[👃]|[👄]|[👅]|[👆]|[👇]|[👈]|[👉]|[👊]|[👋]|[👌]|[👍]|[👎]|[👏]|[👐]|[👑]|[👖]|[👗]|[👥]|[👨]|[👨‍👨‍👧‍👦]|[👩]|[👪]|[👫]|[👬]|[👱]|[👶]|[👸]|[👼]|[👼]|[👽]|[👿]|[💀]|[💁]|[💃]|[💆]|[💉]|[💊]|[💋]|[💎]|[💐]|[💑]|[💓]|[💔]|[💕]|[💖]|[💗]|[💘]|[💙]|[💚]|[💛]|[💜]|[💝]|[💞]|[💟]|[💥]|[💥]|[💦]|[💧]|[💨]|[💪]|[💫]|[💭]|[💭]|[💯]|[💰]|[💵]|[💸]|[💻]|[💿]|[📌]|[📢]|[📣]|[📦]|[📯]|[📷]|[📸]|[📽]|[🔀]|[🔃]|[🔄]|[🔇]|[🔊]|[🔔]|[🔞]|[🔞]|[🔟]|[🔥]|[🔨]|[🔫]|[🔯]|[🔯]|[🔱]|[🔸]|[🕊]|[🕯]|[🕺]|[🖌]|[🖐]|[🖐]|[🖕]|[🗑]|[🗝]|[🗣]|[🗿]|[🤐]|[🤑]|[🤒]|[🤓]|[🤔]|[🤕]|[🤗]|[🤘]|[🤙]|[🤚]|[🤛]|[🤜]|[🤝]|[🤞]|[🤟]|[🤡]|[🤢]|[🤣]|[🤤]|[🤥]|[🤦]|[🤧]|[🤨]|[🤩]|[🤪]|[🤫]|[🤬]|[🤭]|[🤮]|[🤯]|[🤲]|[🤳]|[🤴]|[🤷]|[🤸]|[🤺]|[🥁]|[🥂]|[🥃]|[🥇]|[🥚]|[🥛]|[🥛]|[🥧]|[🥰]|[🥳]|[🥴]|[🥵]|[🥺]|[🦁]|[🦋]|[🧡]|[🦸]`)

	//	urlexp := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)

	//urlexp := regexp.MustCompile(`(?P<url>https?://[^\s]+)`)
	urlexp2 := regexp.MustCompile(`(?P<url>https?.+)`)
	twitterHandleExp := regexp.MustCompile(`@([A-Za-z]+[A-Za-z0-9-_]+)`)
	hashtagHandleExp := regexp.MustCompile(`#([A-Za-z]+[A-Za-z0-9-_]+)`)
	hashtagHandleExp2 := regexp.MustCompile(`@(_*[0-9]*_*[A-Za-z]+[A-Za-z0-9-_]+)`)

	newLineRegExp := regexp.MustCompile(`\n`)

	printString := emojiRx.ReplaceAllString(tweet, `[e]`)
	printString = urlexp2.ReplaceAllString(printString, `[LINK]`)
	printString = twitterHandleExp.ReplaceAllString(printString, `[USERNAME]`)
	printString = hashtagHandleExp.ReplaceAllString(printString, `[HASHTAG]`)
	printString = hashtagHandleExp2.ReplaceAllString(printString, `[HASHTAG]`)
	printString = strings.ToLower(printString)
	printString = newLineRegExp.ReplaceAllString(printString, ``)

	return printString
}
