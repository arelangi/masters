package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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

func streamProcessing() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "XtNItqCBKnR7Y5bUkjbwULlyM", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "eyWtcvjhyXDWwNHVv8pzeSS0j6ov35sGIjJZx2Y9QTd4jEKAHp", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "41533711-YRVp67aiA0IWvhpw1u4VGNSRJSR6Rki3e9BhToHLL", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "uUopkkB83HIEiCl8kuvnCMgb0QdX5XOX0FggvyOmBPKxP", "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = handleTweet
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"chakri88", "chaitanya61", "umasudhir", "GabbarSanghi", "etvteluguindia", "priyadarshi_i", "ChaiBisket", "RaviTeja_offl", "SVCCofficial", "MondayReviews", "RRRMovie", "harikraju", "RajivAluri", "JalapathyG", "ganeshravuri", "UV_Creations", "haarikahassine", "123telugu", "ramjowrites", "idlebraindotcom", "idlebrainjeevi", "SKNonline", "vamsikaka", "aacreations9", "telugufilmnagar", "DVVMovies", "AndhraBoxOffice", "Telugu360", "greatandhranews", "smkoneru", "BhuviOfficial", "davidwarner31", "iamyusufpathan", "SunRisers", "IKReddy_Nirmal", "tgten", "IamKodaliNani", "GHMCOnline", "mmkeeravaani", "ZeeTVTelugu", "SakshiTelangana", "ys_sharmila", "VoiceTelangana", "bonthurammohan", "LakshmiManchu", "MaheshhKathi", "hmtvlive", "trsharish", "ntdailyonline", "vennelakishore", "baraju_SuperHit", "harish2you", "Gopimohan", "konavenkat99", "itsRajTarun", "Mee_Sunil", "ganeshbandla", "IamJagguBhai", "themohanbabu", "ramsayz", "MythriOfficial", "chay_akkineni", "TelanganaDGP", "hydcitypolice", "anusuyakhasba", "JanaSainiks", "SriVijayaNagesh", "AadhiHyper", "TV1Telugu", "telanganafacts", "BhakthiTVorg", "etvtelanganaa", "VSReddy_MP", "YSRCPDMO", "etvandhraprades", "kathimahesh", "AP24x7live", "sakshinews", "tv5newsnow", "SakshiHDTV", "NtvteluguHD", "V6News", "abntelugutv", "iamnagarjuna", "trspartyonline", "JaiTDP", "TV9Telugu", "TelanganaCMO", "Loksatta_Party", "YSRCParty", "JP_LOKSATTA", "ysjagan", "RaoKavitha", "KTRTRS", "naralokesh", "thisisysr", "gvprakash", "MsKajalAggarwal", "RanaDaggubati", "ThisIsDSP", "alluarjun", "AlluSirish", "impradeepmachi", "MukhiSree", "i_nivethathomas", "NANDAMURIKALYAN", "JanaSenaParty", "IamSaiDharamTej", "upasanakonidela", "AkhilAkkineni8", "sivakoratala", "Samanthaprabhu2", "urstrulyMahesh", "tarak9999"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()

}

func handleTweet(tweet *twitter.Tweet) {
	msg := RTTweet{OriginalText: tweet.Text, NormalizedText: normalize(tweet.Text)}

	predictions, err := callMLEngine(msg)
	if err != nil {
		return
	}

	entities := extractEntities(msg.NormalizedText)

	resp := Response{Tweet: msg, Predictions: predictions, Entities: entities}

	jsonResp, _ := json.Marshal(resp)
	streamChan <- jsonResp

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
		//mlog.Error("JSON Decoding error", mlog.Items{"error": err})
		return
	}
	return
}

func otherStream() {
	for {
		select {
		case <-streamChan:
			//fmt.Println("----------", message)
		}
	}
}

type RTTweet struct {
	OriginalText   string `json:"original_text"`
	NormalizedText string `json:"normalized_text"`
}

type Predictions struct {
	FtPrediction   string `json:"ft_prediction"`
	LstmPrediction string `json:"lstm_prediction"`
	NbPrediction   string `json:"nb_prediction"`
	SvmPrediction  string `json:"svm_prediction"`
}

type Response struct {
	Tweet       RTTweet     `json:"tweet"`
	Predictions Predictions `json:"predictions"`
	Entities    []string    `json:"entities"`
}
