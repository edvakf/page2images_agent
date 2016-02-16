package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// API endpoint for page2images
const API string = "http://api.page2images.com/restfullink"

// APIResponse is the page2images API Response struct
type APIResponse struct {
	Status            string `json:"status"`
	EstimatedNeedTime int    `json:"estimated_need_time,omitempty"`
	ImageURL          string `json:"image_url,omitempty"`
	OriginalURL       string `json:"ori_url,omitempty"`
	Duration          int    `json:"duration,omitempty"`
	LeftCalls         string `json:"left_calls,omitempty"`
	ErrNo             int    `json:"errno,omitempty"`
	Msg               string `json:"msg,omitempty"`
}

var key string
var port string
var urlPrefix string

func init() {
	flag.StringVar(&key, "api-key", "", "API Key for page2images.com")
	flag.StringVar(&port, "port", "8080", "http port to listen")
	flag.StringVar(&urlPrefix, "url-prefix", "", "only URLs starting with this prefix are permitted (leave blank to permit any URLs)")
}

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if q.Get("p2i_url") == "" {
		w.WriteHeader(http.StatusNotFound)
	}
	if urlPrefix != "" && !strings.HasPrefix(q.Get("p2i_url"), urlPrefix) {
		w.WriteHeader(http.StatusBadRequest)
	}
	q.Del("p2i_key")
	q.Set("p2i_key", key)

	err := func() error {
		for {
			resp, err := http.Get(API + "?" + q.Encode())
			if err != nil {
				return errors.New("failed to request page2images.com: " + err.Error())
			}
			defer resp.Body.Close()

			var j APIResponse

			var buf bytes.Buffer
			io.Copy(&buf, resp.Body)
			log.Println(buf.String())
			dec := json.NewDecoder(&buf)

			//dec := json.NewDecoder(resp.Body)
			err = dec.Decode(&j)
			if err != nil {
				return errors.New("failed to decode response json: " + err.Error())
			}
			if j.Status == "processing" {
				// sometimes image is ready well before the estimated_need_time, so don't wait for the full length
				time.Sleep(time.Duration(j.EstimatedNeedTime/4) * time.Second)
				continue
			}
			if j.Status == "error" {
				w.WriteHeader(j.ErrNo)
				w.Write([]byte("failed to decode response json: " + j.Msg))
				return nil
			}
			if j.Status == "finished" {
				http.Redirect(w, r, j.ImageURL, http.StatusMovedPermanently)
				return nil
			}
		}
	}()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func main() {
	flag.Parse()
	if key == "" {
		panic("api-key must be set (you must log in to page2images.com and generate one)")
	}

	http.HandleFunc("/", handler)
	log.Println("listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
