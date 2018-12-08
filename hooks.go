package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/go-playground/webhooks.v5/github"
	"os/exec"
)

const (
	path = "/webhooks"
)

func main() {
	f, err := os.OpenFile("hooks.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("Unable to create/open log file %v", err)
	}

	defer f.Close()

	log.SetOutput(f)
	hook, _ := github.New(github.Options.Secret("953507cbb10c25e9284040d4def099f0c57d1813"))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn;t one of the ones asked to be parsed
				log.Printf("There is no event anyway.")
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			pushRequest := payload.(github.PushPayload)
			fmt.Printf("%+v", pushRequest)
		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)
			shellScript := "./script.sh"
			if _, err := exec.Command("/bin/sh", "-c", shellScript).Output(); err != nil {
				log.Printf("Unable to docker compose, %v", err)
			} else {
				log.Printf("A new build was completed.")
			}

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
		default:
			log.Printf("This is not Github hooks request, %v", payload)
			http.Error(w, "Unauthorized access", 401)
		}

	})
	http.ListenAndServe(":3000", nil)
}
