package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {

	tokFile := "env/cred.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil
	}
	return config.Client(context.Background(), tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	
	err = json.NewDecoder(f).Decode(tok)

	b, err := os.ReadFile("env/cred.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	ctx := context.Background()

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}

	fmt.Println("Labels:")
	for _, l := range r.Labels {
		fmt.Printf("- %s\n", l.Name)
	}

	return tok, err
}


func Filter(id string) {
	filterquery := gmail.FilterCriteria{
		Subject: id,
	}

	filter := gmail.Filter{
		Id: "subject-filter",
		Criteria: &filterquery,
	}

	x, err := filter.MarshalJSON()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(x)

}