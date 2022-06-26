package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/vaughan0/go-ini"
)

func makeClient(conf map[string]string) *twitter.Client {
	config := oauth1.NewConfig(conf["consKey"], conf["consSecret"])
	token := oauth1.NewToken(conf["accToken"], conf["accSecret"])
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func getTweets(conf map[string]map[string]string) ([]string, error) {
	count, err := strconv.Atoi(conf["target"]["sampleNum"])
	if err != nil {
		return nil, err
	}

	client := makeClient(conf["oauth"])
	timeline, res, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName: conf["target"]["name"],
		Count:      count,
	})
	if err != nil || res.StatusCode != 200 {
		fmt.Println(err)
		fmt.Println(res)
		return nil, err
	}

	tw := make([]string, 0, count)
	for _, t := range timeline {
		// なんとなく1行にまとめる あとで解析するのでスペースで区切っておく
		s := strings.Replace(t.Text, "\r\n", " ", -1)
		s = strings.Replace(s, "\n", " ", -1)
		s = strings.Replace(s, "\r", " ", -1)

		tw = append(tw, s)
	}

	return tw, nil
}

func makeChain(tweets []string) map[string]map[string]int {
	tk := tokenizer.New()
	chain := make(map[string]map[string]int)

	for _, tw := range tweets {
		ts := tk.Tokenize(tw)
		for i, t := range ts {
			cur := t.Surface
			if cur == "EOS" {
				break
			}

			next := ts[i+1].Surface
			if _, ok := chain[cur]; ok == false {
				chain[cur] = make(map[string]int)
			}
			chain[cur][next]++
		}
	}

	return chain
}

func genTweet(chain map[string]map[string]int) string {
	rand.Seed(time.Now().UnixNano())

	s := ""
	w := "BOS"
	for w != "EOS" {
		type key struct {
			key   string
			thres int
		}

		var sum int = 0
		var keys []key
		for k, v := range chain[w] {
			sum += v
			keys = append(keys, key{k, sum})
		}

		r := rand.Intn(sum) + 1
		w = keys[0].key
		for _, k := range keys {
			if k.thres > r {
				break
			}
			w = k.key
		}
		if w != "EOS" {
			s += (w + "")
		}
	}

	return s
}

func filterWords(src string) string {
	// メンションしない
	rep := regexp.MustCompile(`@[\w]*`)
	dst := rep.ReplaceAllString(src, "")
	// ハッシュタグつけない
	rep = regexp.MustCompile(`#.*`)
	dst = rep.ReplaceAllString(src, "")
	// URLつぶやかない(これきちんと機能していない...)
	rep = regexp.MustCompile(`http[\w:]`)
	dst = rep.ReplaceAllString(src, "")

	return dst
}

func readConf(iniName string) map[string]map[string]string {
	conf := make(map[string]map[string]string)
	ini, err := ini.LoadFile(iniName)
	if err != nil {
		return conf
	}

	// ini -> conf map
	for k, v := range ini {
		conf[k] = make(map[string]string)
		for kk, vv := range v {
			conf[k][kk] = vv
		}
	}

	return conf
}

func main() {
	conf := readConf("ysbot.ini")
	//fmt.Println(conf)

	tw, err := getTweets(conf)
	if err != nil {
		fmt.Println(err)
	} else {
		//for _, t := range tw {
		//	fmt.Println(t)
		//}
	}

	chain := makeChain(tw)
	//for i, w := range chain {
	//	fmt.Println(i, w)
	//}

	s := genTweet(chain)
	s = filterWords(s)
	//fmt.Println(">>>", s)

	client := makeClient(conf["oauth"])
	client.Statuses.Update(s, nil)
}
