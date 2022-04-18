package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {

	routes := getRoutes()
	ts := httptest.NewTLSServer(routes) // create a test server of type http
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL + "/") // add ts.URL because Go will create a server on different routes each time it is run
	if err != nil {
		t.Log(err)
		t.Fatal(err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("for home page, expected status 200 but got %d", resp.StatusCode)
	}

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(bodyText), "awesome") {
		cel.TakeScreenShot(ts.URL+"/", "Hometest", 1500, 1000)
		t.Error("did not find 'awesome'")
	}
}

// this test includes sessions
func TestHome2(t *testing.T) {
	// build request manually
	req, _ := http.NewRequest("GET", "/", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx) // now we have a request that is aware of the session

	rr := httptest.NewRecorder() // create a response request/recorder
	// which takes the place of a ResponseWriter()

	// put something in the session
	// so that we can test our session
	cel.Session.Put(ctx, "test_key", "Hello, world")

	// instead of spinning up a test server as above
	// we will create our own
	h := http.HandlerFunc(testHandlers.Home)
	h.ServeHTTP(rr, req)

	// test to see if it contains the correct status codes
	// this is functionally the same as the the previous test
	if rr.Code != 200 {
		t.Errorf("returned wrong response code; got %d but expected 200", rr.Code)
	}

	// But this way we can also test our session
	if cel.Session.GetString(ctx, "test_key") != "Hello, world" {
		t.Error("did not get the correct value from session")
	}
}
