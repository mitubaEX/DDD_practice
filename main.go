// main is an example web app using Login with Twitter.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/mitubaEX/DDD_practice/service"
)

// Config configures the main ServeMux.
type Config struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
}

// New returns a new ServeMux with app routes.
func New(config *Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", service.Welcom)
	mux.Handle("/profile", service.RequireLogin(http.HandlerFunc(service.ProfileHandler)))
	mux.HandleFunc("/logout", service.Logout)
	// 1. Register Twitter login and callback handlers
	oauth1Config := &oauth1.Config{
		ConsumerKey:    config.TwitterConsumerKey,
		ConsumerSecret: config.TwitterConsumerSecret,
		CallbackURL:    "http://localhost:8080/twitter/callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	mux.Handle("/twitter/login", twitter.LoginHandler(oauth1Config, nil))
	mux.Handle("/twitter/callback", twitter.CallbackHandler(oauth1Config, service.Login(), nil))
	return mux
}

// main creates and starts a Server listening.
func main() {
	const address = "localhost:8080"
	// read credentials from environment variables if available
	config := &Config{
		TwitterConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		TwitterConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	}
	// allow consumer credential flags to override config fields
	consumerKey := flag.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flag.String("consumer-secret", "", "Twitter Consumer Secret")
	flag.Parse()
	if *consumerKey != "" {
		config.TwitterConsumerKey = *consumerKey
	}
	if *consumerSecret != "" {
		config.TwitterConsumerSecret = *consumerSecret
	}
	if config.TwitterConsumerKey == "" {
		log.Fatal("Missing Twitter Consumer Key")
	}
	if config.TwitterConsumerSecret == "" {
		log.Fatal("Missing Twitter Consumer Secret")
	}

	log.Printf("Starting Server listening on %s\n", address)
	err := http.ListenAndServe(address, New(config))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
