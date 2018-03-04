package uaa

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type UaaClient struct {
	uid, pass string
	uaaConfig *oauth2.Config
}

func NewUaaClient(uaaHost, uid, pass string) *UaaClient {
	u := UaaClient{
		uid:  uid,
		pass: pass,
		uaaConfig: &oauth2.Config{
			Scopes:   []string{""},
			ClientID: "cf",
			Endpoint: oauth2.Endpoint{
				//AuthURL:  "https://uaa.system.aws-usw02-pr.ice.predix.io/oauth/authorize",
				//TokenURL: "https://uaa.system.aws-usw02-pr.ice.predix.io/oauth/token",
				AuthURL:  "https://" + uaaHost + "/oauth/authorize",
				TokenURL: "https://" + uaaHost + "/oauth/token",
			},
		},
	}
	return &u
}

func (u *UaaClient) Authenticate() (*http.Client, error) {
	token, err := u.uaaConfig.PasswordCredentialsToken(context.Background(), u.uid, u.pass)
	if err != nil {
		log.Printf("Could not get token for '%v', '%v': %v", u.uid, u.pass, err)
		return nil, err
	}

	return u.uaaConfig.Client(context.Background(), token), nil
}
