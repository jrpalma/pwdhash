package rest

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jrpalma/pwdhash/config"
	"github.com/jrpalma/pwdhash/logs"
	"github.com/jrpalma/pwdhash/task"
)

type testShutdown struct {
	wasShutdown bool
	shutdownErr error
}

func (ts *testShutdown) Shutdown() error {
	ts.wasShutdown = true
	return ts.shutdownErr
}

type HandlerHarness struct {
	ts      testShutdown
	conf    config.Config
	log     logs.Logger
	handler *handler
	server  *httptest.Server
}

func newHandlerHarness(api string) (*HandlerHarness, error) {
	hh := &HandlerHarness{}
	log, err := logs.NewStreamLogger(logs.STDERR, logs.INFO)
	if err != nil {
		return nil, err
	}

	hh.log = log
	hh.handler = newHandler(hh.conf, hh.log, hh.ts.Shutdown)

	if api == "newHash" {
		hh.server = httptest.NewServer(http.HandlerFunc(hh.handler.newHash))
	} else if api == "checkHash" {
		hh.server = httptest.NewServer(http.HandlerFunc(hh.handler.checkHash))
	} else if api == "stats" {
		hh.server = httptest.NewServer(http.HandlerFunc(hh.handler.stats))
	} else if api == "shutdown" {
		hh.server = httptest.NewServer(http.HandlerFunc(hh.handler.shutdown))
	} else {
		return nil, fmt.Errorf("Invalid handler api")
	}

	return hh, nil
}

func sendPost(URL string) (task.Result, error) {
	res := task.Result{}
	response, err := http.Post(URL, "", nil)
	if err != nil {
		return res, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}

	res.Code = response.StatusCode
	res.Message = string(body)

	return res, nil
}

func postPassword(field, pwd, method, URL string) (task.Result, error) {
	res := task.Result{}
	form := url.Values{}
	form.Add(field, pwd)

	reqID := fmt.Sprintf("%v", rand.Intn(5000))

	request, err := http.NewRequest(method, URL, strings.NewReader(form.Encode()))
	if err != nil {
		return res, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("X-Request-ID", reqID)
	request.PostForm = form

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return res, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}

	res.Code = response.StatusCode
	res.Message = string(body)

	return res, nil
}

func getRequest(URL string) (task.Result, error) {
	res := task.Result{}
	response, err := http.Get(URL)
	if err != nil {
		return res, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return res, err
	}

	res.Code = response.StatusCode
	res.Message = string(body)

	return res, nil
}

func TestHandler_newHashMethodNotAllowed(t *testing.T) {
	h, err := newHandlerHarness("newHash")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := postPassword("password", "secret", "GET", h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("newHash returned: %v", res.Code)
	}
}

func TestHandler_newInvalidField(t *testing.T) {
	h, err := newHandlerHarness("newHash")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := postPassword("unknown", "secret", "POST", h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusBadRequest {
		t.Errorf("newHash returned: %v", res.Code)
	}
}

func TestHandler_newHash(t *testing.T) {
	h, err := newHandlerHarness("newHash")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := postPassword("password", "secret", "POST", h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusCreated {
		t.Errorf("newHash returned: %+v", res)
	}
}

func TestHandler_checkHashMethodNotAllowed(t *testing.T) {
	h, err := newHandlerHarness("checkHash")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := postPassword("password", "secret", "POST", h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("checkHash returned: %+v", res)
	}
}

func TestHandler_checkHashInvalidHashID(t *testing.T) {
	h, err := newHandlerHarness("checkHash")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := getRequest(h.server.URL + "/invalidInteger")
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusBadRequest {
		t.Errorf("checkHash returned: %+v", res)
	}
}

func TestHandler_checkHashNotFound(t *testing.T) {
	h, err := newHandlerHarness("checkHash")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := getRequest(h.server.URL + "/100")
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusNotFound {
		t.Errorf("checkHash returned: %+v", res)
	}
}

func TestHandler_statsMethodNotAllowed(t *testing.T) {
	h, err := newHandlerHarness("stats")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := postPassword("password", "secret", "POST", h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("stats returned: %+v", res)
	}
}

func TestHandler_checkHashOK(t *testing.T) {
	h, err := newHandlerHarness("stats")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := getRequest(h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusOK {
		t.Errorf("stats returned: %+v", res)
	}
}

func TestHandler_shutdownMethodNotAllowed(t *testing.T) {
	h, err := newHandlerHarness("shutdown")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := getRequest(h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusMethodNotAllowed {
		t.Errorf("shutdown returned: %+v", res)
	}
}

func TestHandler_shutdownOK(t *testing.T) {
	h, err := newHandlerHarness("shutdown")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := sendPost(h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusOK {
		t.Errorf("shutdown returned: %+v", res)
	}
}

func TestHandler_doubleShutdown(t *testing.T) {
	// NOTE: This test might not be deterministic
	// If the machine is low on resoruces.

	h, err := newHandlerHarness("shutdown")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := sendPost(h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	res, err = sendPost(h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusInternalServerError {
		t.Errorf("shutdown returned: %+v", res)
	}
}

func TestHandler_failServerShutdown(t *testing.T) {

	h, err := newHandlerHarness("shutdown")
	h.ts.shutdownErr = fmt.Errorf("Fail shutdown")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	defer h.server.Close()
	res, err := sendPost(h.server.URL)
	if err != nil {
		t.Errorf("Failed to post password: %v", err)
		return
	}

	if res.Code != http.StatusOK {
		t.Errorf("shutdown returned: %+v", res)
	}

}
