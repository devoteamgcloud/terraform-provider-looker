package acceptance_tests

import (
	"context"

	"github.com/k0kubun/pp/v3"
	_ "github.com/stretchr/testify/assert"

	"net/http"
	"os"
	"testing"

	lookergo "github.com/devoteamgcloud/terraform-provider-looker/pkg/lookergo"
)

// export LOOKER_API_CLIENT_ID="dFNmbP4jzXqp" LOOKER_API_CLIENT_SECRET="KFwpLkMf393z" LOOKER_API_ENDPOINT="https://dev.l0c4l.host:9443/api"
// In Test env Goland: https_proxy=http://172.31.100.95:9090;LOOKER_API_CLIENT_ID=dFNmbP4jzXqp;LOOKER_API_CLIENT_SECRET=KFwpLkMf393z;LOOKER_API_ENDPOINT=https://dev.l0c4l.host:9443/api/
var clientId, clientSecret, apiEndpoint string

var (
	ctx    = context.TODO()
	ac     = lookergo.ApiConfig{}
	client *lookergo.Client
)

func init() {
	clientId = os.Getenv("LOOKER_API_CLIENT_ID")
	clientSecret = os.Getenv("LOOKER_API_CLIENT_SECRET")
	apiEndpoint = os.Getenv("LOOKER_API_ENDPOINT")

	ac = lookergo.ApiConfig{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		BaseURL:      apiEndpoint,
	}
}

func setup() {
	if client == nil {
		client = lookergo.NewFromApiv3Creds(ac)
	}
}

func TestAcceptance_Panic(t *testing.T) {
	panic("bye")
}

func TestAcceptance_Demo(t *testing.T) {
	setup()

	url := ac.BaseURL + "4.0/users/1"
	req, _ := client.NewRequest(ctx, http.MethodGet, url, nil)

	body := new(lookergo.User)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	pp.Print(body)
}
