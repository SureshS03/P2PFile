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

	tokFile := "env/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)

}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err

}

func OAuth() error {
	ctx := context.Background()
	b, err := os.ReadFile("env/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return err
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return err
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
		return err
	}

	fmt.Println(srv.Users)
	res, err := getMailsBySubject(srv, "subject:Mail.pngchu185752")
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}


func getMailsBySubject(srv *gmail.Service, subject string) ([]*gmail.Message, error) {

	user := "me" // authenticated user
	query := fmt.Sprintf("subject:%s", subject)

	resp, err := srv.Users.Messages.List(user).
		Q(query).
		MaxResults(10).
		Do()

	if err != nil {
		return nil, err
	}

	return resp.Messages, nil
}

func getMessage(srv *gmail.Service, msgID string) (*gmail.Message, error) {
	return srv.Users.Messages.Get("me", msgID).
		Format("full").
		Do()
}
