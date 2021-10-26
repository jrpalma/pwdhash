package rest

import (
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/jrpalma/pwdhash/config"
	"github.com/jrpalma/pwdhash/logs"
)

type ServerHarness struct {
	conf   config.Config
	log    logs.Logger
	server *Server
}

func newServerHarness(addr string) (*ServerHarness, error) {
	h := &ServerHarness{}
	log, err := logs.NewStreamLogger(logs.STDERR, logs.INFO)
	if err != nil {
		return nil, err
	}
	h.log = log
	h.conf.ServerAddress = addr
	h.server = NewServer(h.conf, h.log)

	return h, nil
}

func TestServer_NewServer(t *testing.T) {
	_, err := newServerHarness(":3700")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}
}

func TestServer_RunShutdown(t *testing.T) {
	sh, err := newServerHarness(":3701")
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		wg.Done()
		err = sh.server.Run()
		if err != http.ErrServerClosed {
			t.Errorf("Faield to run server: %v", err)
		}

	}()

	wg.Wait()
	time.Sleep(time.Second)

	err = sh.server.Shutdown()
	if err != nil {
		t.Errorf("Faield to shutdown server: %v", err)
	}
}

func TestServer_RunRequest(t *testing.T) {
	port := ":3702"
	sh, err := newServerHarness(port)
	if err != nil {
		t.Errorf("Failed to setup test: %v", err)
		return
	}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		wg.Done()
		err = sh.server.Run()
		if err != http.ErrServerClosed {
			t.Errorf("Faield to run server: %v", err)
		}

	}()

	wg.Wait()
	time.Sleep(time.Second)
	defer sh.server.Shutdown()

	URL := "http://127.0.0.1" + port + v1 + "/stats"
	res, err := getRequest(URL)

	if err != nil {
		t.Errorf("Request failed: %+v", err)
	}

	if res.Code != 200 {
		t.Errorf("stats returned: %+v", res)
	}
}
