package uaa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dave-malone/oauth"
)

func TestNewClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
      "access_token" : "4a53a3331b2445cfaca43c9af00439e8",
      "token_type" : "bearer",
      "expires_in" : 43199,
      "scope" : "clients.read emails.write scim.userids password.write idps.write notifications.write oauth.login scim.write critical_notifications.write",
      "jti" : "4a53a3331b2445cfaca43c9af00439e8"
    }`)
	}))
	defer ts.Close()

	if _, err := getUaac(ts.URL); err != nil {
		t.Errorf("Failed to initialize uaa client; %v", err)
		return
	}
}

func TestGetServerInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
      "app" : {
        "version" : "3.5.0-SNAPSHOT"
      },
      "links" : {
        "uaa" : "http://localhost:8080/uaa",
        "passwd" : "/forgot_password",
        "login" : "http://localhost:8080/uaa",
        "register" : "/create_account"
      },
      "zone_name" : "uaa",
      "entityID" : "cloudfoundry-saml-login",
      "commit_id" : "6681d65",
      "idpDefinitions" : {
        "SAML" : "http://localhost:8080/uaa/saml/discovery?returnIDParam=idp&entityID=cloudfoundry-saml-login&idp=SAML&isPassive=true"
      },
      "prompts" : {
        "username" : [ "text", "Email" ],
        "password" : [ "password", "Password" ]
      },
      "timestamp" : "2016-05-18T18:20:54+0000"
    }`)
	}))
	defer ts.Close()

	uaac, err := getUaac(ts.URL)
	if err != nil {
		t.Errorf("Failed to get uaa client; %v", err)
		return
	}

	serverInfo, err := uaac.GetServerInfo()
	if err != nil {
		t.Errorf("Failed to get ServerInfo: %v", err)
	}

	if len(serverInfo.Version()) == 0 {
		t.Error("ServerInfo.Version() was blank")
	}
}

func TestGetListIdentityZones(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{
        "id": "dummy-id",
        "subdomain": "dummy-subdomain",
        "name": "dummy-name",
        "version": 1,
        "description": "Dummy Description",
        "created": 946710000000,
        "last_modified": 946710000000
    }]`)
	}))
	defer ts.Close()

	uaac, err := getUaac(ts.URL)
	if err != nil {
		t.Errorf("Failed to get uaa client; %v", err)
		return
	}

	zones, err := uaac.ListIdentityZones()
	if err != nil {
		t.Errorf("Failed to get Zones: %v", err)
	}

	if len(zones) == 0 {
		t.Error("[]Zone was empty")
	}
}

func TestGetListUsersWithUaa20Model(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("./testdata/2.0.0/GET.Users.json")

		if err != nil {
			panic("Failed to read ./testdata/2.0.0/GET.Users.json: " + err.Error())
		}

		w.Write(data)
	}))
	defer ts.Close()

	uaac, err := getUaac(ts.URL)
	if err != nil {
		t.Errorf("Failed to get uaa client; %v", err)
		return
	}

	users, err := uaac.ListUsers()
	if err != nil {
		t.Errorf("Failed to get Users: %v", err)
	}

	if len(users.Users) == 0 {
		t.Error("[]Users was empty")
	}
}

func getUaac(apiUrl string) (Client, error) {
	clientConfig := &oauth.ClientConfig{
		ApiAddress:   apiUrl,
		ClientID:     "fake-client",
		ClientSecret: "big-secret",
	}

	oauthClient, err := oauth.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize oauth client; %v", err)
	}

	uaac := &uaaClient{oauthClient}
	return uaac, nil
}
