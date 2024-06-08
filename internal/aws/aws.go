package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

type SigninToken struct {
	SigninToken string
}

type Profile struct {
	name       string
	account_id string
	role_name  string
}

func NewProfileFromConfig(name string, cfg ini.Section) *Profile {
	p := Profile{name: name}
	p.account_id = cfg.Key("sso_account_id").String()
	p.role_name = cfg.Key("sso_role_name").String()
	return &p
}

func GetAwsProfiles() (profiles []Profile, err error) {
	cfg, err := ini.Load(config.DefaultSharedConfigFilename())
	if err != nil {
		fmt.Println("Error loading config file.")
		return profiles, err
	}
	for _, section := range cfg.Sections() {
		name, is_profile := strings.CutPrefix(section.Name(), "profile ")
		if is_profile {
			profiles = append(profiles, *NewProfileFromConfig(name, *section))
		}
	}
	return profiles, err
}

func GetLoginUrl(profile string, region string) (string, error) {
	var login_url string
	destination_url := url.QueryEscape("https://" + region + ".console.aws.amazon.com")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return login_url, err
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return login_url, err
	}

	temp_creds := make(map[string]string)
	temp_creds["sessionId"] = creds.AccessKeyID
	temp_creds["sessionKey"] = creds.SecretAccessKey
	temp_creds["sessionToken"] = creds.SessionToken

	json_creds, err := json.Marshal(temp_creds)
	if err != nil {
		return login_url, err
	}

	get_signin_url := "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=" + url.QueryEscape(string(json_creds))
	resp, err := http.Get(get_signin_url)
	if err != nil {
		return login_url, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	var token SigninToken
	json.Unmarshal(body, &token)

	login_url = "https://signin.aws.amazon.com/federation?Action=login&Issuer=&Destination=" + destination_url + "&SigninToken=" + token.SigninToken
	return login_url, nil
}
