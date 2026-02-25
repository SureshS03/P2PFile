package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

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

func OAuth() (*gmail.Service, error) {
	ctx := context.Background()
	b, err := os.ReadFile("env/credentials.json")
	if err != nil {
		ErrPrinter(err)
		return nil, err
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		ErrPrinter(err)
		return nil, err
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		ErrPrinter(err)
		return nil, err
	}

	CrrPrinter("OAuth Done")
	return srv, nil
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

func ReadMailAndDownloadAttachments(
	srv *gmail.Service,
	subject string,
) ([]string ,error) {

	query := fmt.Sprintf("subject:%s", subject)

	resp, err := srv.Users.Messages.List("me").
		Q(query).
		MaxResults(5).
		Do()
	if err != nil {
		return nil, err
	}

	var paths []string

	for _, m := range resp.Messages {

		msg, err := srv.Users.Messages.Get("me", m.Id).
			Format("full").
			Do()
		if err != nil {
			continue
		}

		paths, err = downloadAttachments(srv, msg, "attachments")
		if err != nil {
			ErrPrinter(err)
			continue
		}

		if len(paths) == 0 {
			fmt.Println("No attachments in mail:", msg.Id)
			continue
		}

		CrrPrinter("Downloaded attachments:")
		for _, p := range paths {
			CrrPrinter(" →" + p)
		}
	}

	return paths, nil
}

func downloadAttachments(
	srv *gmail.Service,
	msg *gmail.Message,
	baseDir string,
) ([]string, error) {

	var savedPaths []string

	// Create dir per message
	msgDir := filepath.Join(baseDir, msg.Id)
	if err := os.MkdirAll(msgDir, 0755); err != nil {
		return nil, err
	}

	var walkParts func(parts []*gmail.MessagePart) error

	walkParts = func(parts []*gmail.MessagePart) error {
		for _, part := range parts {

			if part.Filename != "" && part.Body != nil && part.Body.AttachmentId != "" {

				att, err := srv.Users.Messages.Attachments.
					Get("me", msg.Id, part.Body.AttachmentId).
					Do()
				if err != nil {
					return err
				}

				data, err := base64.URLEncoding.DecodeString(att.Data)
				if err != nil {
					return err
				}

				filePath := filepath.Join(msgDir, part.Filename)
				if err := os.WriteFile(filePath, data, 0644); err != nil {
					return err
				}

				savedPaths = append(savedPaths, filePath)
			}

			if len(part.Parts) > 0 {
				if err := walkParts(part.Parts); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if msg.Payload != nil {
		if err := walkParts(msg.Payload.Parts); err != nil {
			return nil, err
		}
	}

	return savedPaths, nil
}