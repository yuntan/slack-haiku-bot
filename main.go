package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/mattn/go-haiku"
)

type slackResponse struct {
	Text string `json:"text"`
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.FormValue("token") != TOKEN {
		log.Println("error: invalid token")
		return
	}
	text := req.FormValue("text")
	log.Printf("new message \"%s\" in #%s\n", text, req.FormValue("channel_name"))
	haikus := haiku.Find(text, []int{5, 7, 5})
	if len(haikus) == 0 {
		log.Printf("found no haiku")
		return
	}

	log.Printf("found %d haiku\n", len(haikus))
	for _, h := range haikus {
		log.Println(h)
	}

	resp := map[string]string{
		"text": strings.Join(haikus, "\n"),
	}
	b, err := json.Marshal(resp)
	if err != nil {
		log.Println("error: failed to marshal map to json: ", err)
		return
	}
	w.Write(b)
}

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatalln(err)
	}
}
