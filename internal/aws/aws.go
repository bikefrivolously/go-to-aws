package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"gopkg.in/ini.v1"
)

const (
	LoginUrlDuration time.Duration = time.Duration(15) * time.Minute
	DefaultRegion                  = "us-east-1"
)

type SigninToken struct {
	SigninToken string
}

type Profile struct {
	Name           string
	AccountId      string
	RoleName       string
	Region         string
	LoginUrl       string
	LoginUrlExpiry time.Time
}

func NewProfileFromConfig(name string, cfg ini.Section) *Profile {
	p := Profile{Name: name}
	p.AccountId = cfg.Key("sso_account_id").String()
	p.RoleName = cfg.Key("sso_role_name").String()
	p.Region = cfg.Key("region").String()
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

func (p *Profile) GetDefaultLoginUrl() error {
	return p.GetLoginUrl(DefaultRegion)
}

func (p *Profile) GetLoginUrl(region string) error {
	destination_url := url.QueryEscape("https://" + region + ".console.aws.amazon.com")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(p.Name),
	)
	if err != nil {
		return err
	}

	creds, err := cfg.Credentials.Retrieve(context.TODO())
	if err != nil {
		return err
	}

	temp_creds := make(map[string]string)
	temp_creds["sessionId"] = creds.AccessKeyID
	temp_creds["sessionKey"] = creds.SecretAccessKey
	temp_creds["sessionToken"] = creds.SessionToken

	json_creds, err := json.Marshal(temp_creds)
	if err != nil {
		return err
	}

	get_signin_url := "https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=" + url.QueryEscape(string(json_creds))
	resp, err := http.Get(get_signin_url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	var token SigninToken
	json.Unmarshal(body, &token)

	p.LoginUrl = "https://signin.aws.amazon.com/federation?Action=login&Issuer=&Destination=" + destination_url + "&SigninToken=" + token.SigninToken
	p.LoginUrlExpiry = time.Now().Add(LoginUrlDuration)
	return nil
}
