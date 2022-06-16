package lookergo

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/k0kubun/pp/v3"
	_ "github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	ctx    = context.TODO()
	client *Client
	server *httptest.Server
)

// < Helpers

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	serverURL, _ := url.Parse(server.URL)
	client.BaseURL = serverURL
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	expected := url.Values{}
	for k, v := range values {
		expected.Add(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatalf("parseForm(): %v", err)
	}

	if !reflect.DeepEqual(expected, r.Form) {
		t.Errorf("Request parameters = %v, expected %v", r.Form, expected)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func rBodyNotExpected(body interface{}, expected interface{}) string {
	return pp.Sprintf("Response body:\n %v\nExpected:\n %v\n", body, expected)
}

func errGotWant(resource string, got interface{}, want interface{}) (string, string) {
	if diff := cmp.Diff(got, want); diff != "" {
		return pp.Sprintf("%s\ngot:\n%v\nwant:\n%v\n", resource, got, want),
			fmt.Sprintf("%s mismatch (-want +got):\n%s", resource, diff)
		// return "", fmt.Sprintf("%s mismatch (-want +got):\n%s", resource, diff)
	}
	return "", ""
}

// !> Helpers

/* func TestNewRequest_get(t *testing.T) {
	c := NewClient(nil)

	reqURL := defaultBaseURL + "foo"
	req, _ := c.NewRequest(ctx, http.MethodGet, reqURL, nil)

} */

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {

		t.Error(rBodyNotExpected(body, expected))
	}
}

func TestAuth3(t *testing.T) {
	_ = os.Setenv("HTTP_PROXY", "http://localhost:9090")

	c := NewClient(nil)

	ctx := context.Background()
	err := c.SetBaseURL("http://dev.l0c4l.host:3001/api")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	err = c.SetOauthCredentials(ctx, "mrClientId", "mrClientSecret")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	type foo struct {
		A string
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	_, err = c.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

}

func TestAuth5(t *testing.T) {
	_ = os.Setenv("HTTP_PROXY", "http://localhost:9090")

	c := NewClient(nil)

	ctx := context.Background()
	err := c.SetBaseURL("http://dev.l0c4l.host:3001/api")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}
	err = c.SetOauthCredentials(ctx, "mrClientId", "mrClientSecret")
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	type foo struct {
		A string
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	_, err = c.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	devClient := NewClient(nil)
	_ = devClient.SetBaseURL(c.BaseURL.String())

	err = devClient.SetOauthStaticToken(ctx, &oauth2.Token{AccessToken: "very-unique", TokenType: "Bearer"})
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	req2, _ := c.NewRequest(ctx, http.MethodGet, "/", nil)
	body2 := new(foo)
	_, err = devClient.Do(context.Background(), req2, body2)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

}
