package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	gosxnotifier "github.com/deckarep/gosx-notifier"
	"github.com/hunterlong/shapeshift"
)

type Pair struct {
	From      string
	To        string
	Sign      string
	Threshold float64
	Amount    float64
}

func main() {

	help := flag.Bool("h", false, "show usage")
	popup := flag.Bool("popup", true, "show popup message (tested on Mac only)")
	interval := flag.Int("interval", 30, "poll at this interval (seconds)")
	flag.Parse()
	if *help {
		fmt.Println(`Usage example: shapeshift-notifier -popup=false -interval=32 "snt_bat,>0.75,=100000" "eth_btc<0.01" "rlc_gnt,=150"
Defaults: popup=true, interval=30, args="eth_btc,>0.1,=0" 
Signs: Only > and < are allowed for operations. = indicates the amount to convert.  Only the first part with token codes is mandatory.
`)
		return
	}

	var pairs []Pair
	if len(os.Args) > 1 {
		pairs = parseCmdLinePairs(flag.Args())
	} else {
		snt_bat := Pair{
			From:      "eth",
			To:        "btc",
			Sign:      ">",
			Threshold: 0.1,
			Amount:    0,
		}

		pairs = []Pair{snt_bat}
	}

	fmt.Println("Checking for pairs:")
	for _, p := range pairs {
		fmt.Println("\t", p)
	}
	fmt.Println()

	//run immediately first
	for _, p := range pairs {
		checkAndNotify(p, popup)
	}
	for range time.Tick(time.Second * time.Duration(*interval)) {
		for _, p := range pairs {
			checkAndNotify(p, popup)
		}
	}
}

func checkAndNotify(p Pair, popup *bool) {

	rate := getrate(p.From, p.To)
	fmt.Printf("%s: %s_%s: %f", time.Now().Format("3:04"), p.From, p.To, rate)
	switch p.Sign {
	case ">":
		if rate > p.Threshold {
			fmt.Printf(" | matched %s %f", p.Sign, p.Threshold)
			if *popup {
				notify(p.From, p.To, p.Sign, p.Threshold, rate)
			}
		}
	case "<":
		if rate < p.Threshold {
			fmt.Printf(" | matched %s %f", p.Sign, p.Threshold)
			if *popup {
				notify(p.From, p.To, p.Sign, p.Threshold, rate)
			}
		}
	}

	if p.Amount > 0 {
		fmt.Printf(" | %f %s = %f %s", p.Amount, p.From, p.Amount*rate, p.To)
	}
	fmt.Println()
}

func getrate(from, to string) float64 {
	pair := shapeshift.Pair{fmt.Sprintf("%s_%s", from, to)}

	rate, err := pair.GetRates()

	if err != nil {
		panic(err)
	}

	//fmt.Printf("%s: %s_%s: %f\n", time.Now().Format("3:04"), from, to, rate)
	return rate
}

func notify(from, to, sign string, threshold, rate float64) {
	//At a minimum specifiy a message to display to end-user.
	note := gosxnotifier.NewNotification("shapeshift-notifier threshold crossed.")

	//Optionally, set a title
	//note.Title = "It's money making time 💰"
	note.Title = fmt.Sprintf("%s_%s = %f", from, to, rate)

	//Optionally, set a subtitle
	//note.Subtitle = "My subtitle"
	note.Subtitle = fmt.Sprintf("%s %f", sign, threshold)

	//Optionally, set a sound from a predefined set.
	note.Sound = gosxnotifier.Basso

	//Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
	note.Group = "com.sathishvj.shapeshit-notifier.1"

	//Optionally, set a sender (Notification will now use the Safari icon)
	//note.Sender = "com.apple.Safari"

	//Optionally, specifiy a url or bundleid to open should the notification be
	//clicked.
	//note.Link = "http://www.yahoo.com" //or BundleID like: com.apple.Terminal

	//Optionally, an app icon (10.9+ ONLY)
	//note.AppIcon = "gopher.png"

	//Optionally, a content image (10.9+ ONLY)
	//note.ContentImage = "gopher.png"

	//Then, push the notification
	err := note.Push()

	//If necessary, check error
	if err != nil {
		log.Println("Uh oh!")
	}
}

// should be in the format: a_b>0.1, b_c<0.2
// only < and > are allowed
func parseCmdLinePairs(args []string) []Pair {
	var pairs []Pair
	//var err error

	/*
		for _, v := range args {
			re := regexp.MustCompile("(.*)_(.*)([><])(.*)")
			match := re.FindAllStringSubmatch(v, 1)
			p := Pair{
				From: match[0][1],
				To:   match[0][2],
				Sign: match[0][3],
			}
			p.Threshold, err = strconv.ParseFloat(match[0][4], 64)
			if err != nil {
				log.Fatal(err)
			}

			pairs = append(pairs, p)
		}
	*/

	for _, arg := range args {
		parts := strings.Split(arg, ",")
		re := regexp.MustCompile("(.*)_(.*)")
		match := re.FindAllStringSubmatch(parts[0], 1)
		p := Pair{
			From: match[0][1],
			To:   match[0][2],
		}

		for _, part := range parts[1:] {
			f, err := strconv.ParseFloat(part[1:], 64)
			if err != nil {
				log.Fatalf("Cannot convert float value: %s. %v\n", part[1:], err)
			}
			switch part[0] {
			case '<':
				p.Sign = "<"
				p.Threshold = f
			case '>':
				p.Sign = ">"
				p.Threshold = f
			case '=':
				p.Amount = f
			}
		}

		pairs = append(pairs, p)

		fmt.Printf("%+v\n", p)
	}

	return pairs
}
