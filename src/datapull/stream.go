package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	_ "github.com/lib/pq"
	"gitlab.logzero.in/arelangi/mlog"
)

var (
	_dbUser     = "aditya"
	_dbPassword = "aditya"
	_dbName     = "aditya"
	_dbHost     = "localhost"
	db          *sql.DB
)

func init() {
	var err error
	user := os.Getenv("DBUSER")
	if user != "" {
		_dbUser = user
	}

	password := os.Getenv("DBPASSWORD")
	if password != "" {
		_dbPassword = password
	}

	name := os.Getenv("DBNAME")
	if name != "" {
		_dbName = name
	}

	host := os.Getenv("DBHOST")
	if host != "" {
		_dbHost = host
	}

	logLevel := os.Getenv("logLevel")
	if logLevel != "" {
		mlog.SetLogLevel(logLevel)
	}

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable",
		_dbUser, _dbPassword, _dbName, _dbHost)
	db, err = sql.Open("postgres", dbinfo)
	if err != nil {
		mlog.Error("Failed to connect to the database", mlog.Items{"error": err})
	}

}

func main() {
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

	db.Close()
}

func handleTweet(tweet *twitter.Tweet) {
	fmt.Println(tweet.User.ScreenName, ": ", tweet.Text)
	saveTweet(tweet.Text)

}

func saveTweet(tweet string) {
	fmt.Println(tweet)
	_, err := db.Exec("INSERT INTO masters.tweets(tweet) VALUES ($1);", tweet)
	if err != nil {
		mlog.Error(fmt.Sprintf("The error is %s", err))
	}
}
