package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jrpalma/pwdhash/config"
	"github.com/jrpalma/pwdhash/logs"
	"github.com/jrpalma/pwdhash/task"
)

type serverShutdownFunc func() error

func newHandler(conf config.Config, log logs.Logger, serverShutdown serverShutdownFunc) *handler {
	return &handler{
		taskMgr:        task.NewManager(conf),
		serverShutdown: serverShutdown,
		config:         conf,
		log:            log,
	}
}

type handler struct {
	config         config.Config
	log            logs.Logger
	taskMgr        *task.Manager
	serverShutdown serverShutdownFunc
}

func (h *handler) newHash(w http.ResponseWriter, r *http.Request) {
	callInfo := h.getCallInfo(r)

	if r.Method != "POST" {
		h.sendStatus(w, callInfo, http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.sendStatus(w, callInfo, http.StatusInternalServerError)
		h.log.Errorf("%v Failed to parse form: %v", callInfo, err)
		return
	}

	if !r.PostForm.Has(formFieldName) {
		h.sendStatus(w, callInfo, http.StatusBadRequest)
		h.log.Errorf("%v Field password is missing", callInfo)
		return
	}

	pwd := r.PostForm.Get(formFieldName)
	res := h.taskMgr.NewTask(pwd)

	h.sendTaskResult(w, callInfo, res)
}
func (h *handler) checkHash(w http.ResponseWriter, r *http.Request) {
	callInfo := h.getCallInfo(r)

	if r.Method != "GET" {
		h.sendStatus(w, callInfo, http.StatusMethodNotAllowed)
		return
	}

	tokens := strings.Split(r.URL.Path, "/")
	hashID := tokens[len(tokens)-1]

	res := h.taskMgr.Check(hashID)

	// TODO: We need to get a better approximation
	// for the retry instead of using this value.
	if res.Code == http.StatusServiceUnavailable {
		seconds := fmt.Sprintf("%v", h.config.MaxTaskSeconds)
		w.Header().Add("Retry-After", seconds)
	}

	h.sendTaskResult(w, callInfo, res)
}
func (h *handler) stats(w http.ResponseWriter, r *http.Request) {
	callInfo := h.getCallInfo(r)

	if r.Method != "GET" {
		h.sendStatus(w, callInfo, http.StatusMethodNotAllowed)
		return
	}

	stats, res := h.taskMgr.Stats()
	if res.Code != http.StatusOK {
		h.sendTaskResult(w, callInfo, res)
		return
	}

	h.sendJSON(w, callInfo, stats)
}
func (h *handler) shutdown(w http.ResponseWriter, r *http.Request) {
	callInfo := h.getCallInfo(r)

	if r.Method != "POST" {
		h.sendStatus(w, callInfo, http.StatusMethodNotAllowed)
		return
	}

	res := h.taskMgr.Shutdown()
	if res.Code != http.StatusOK {
		h.sendTaskResult(w, callInfo, res)
		return
	}

	h.sendTaskResult(w, callInfo, res)
	go h.startShutdown(callInfo)
}

func (h *handler) startShutdown(rid string) {
	h.log.Infof("%v Shutdown: Waiting for pending tasks...", rid)

	h.taskMgr.WaitForPendingTasks()

	h.log.Infof("%v Shutdown: All tasks finished. Shutting down server...", rid)
	err := h.serverShutdown()
	if err != nil {
		h.log.Errorf("%v Server failed to shutdown: %v", rid, err)
	}
}

func (h *handler) sendStatus(w http.ResponseWriter, callInfo string, code int) {
	http.Error(w, http.StatusText(code), code)
	h.logCall(w, callInfo, code)
}

func (h *handler) sendTaskResult(w http.ResponseWriter, callInfo string, result task.Result) {
	http.Error(w, result.Message, result.Code)
	h.logCall(w, callInfo, result.Code)
}

func (h *handler) sendJSON(w http.ResponseWriter, callInfo string, val interface{}) {
	data, err := json.Marshal(val)
	if err != nil {
		res := task.Result{Message: err.Error(), Code: http.StatusInternalServerError}
		h.sendTaskResult(w, callInfo, res)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

	h.logCall(w, callInfo, http.StatusOK)
}

func (h *handler) logCall(w http.ResponseWriter, callInfo string, code int) {
	success := code >= 200 && code <= 299
	status := http.StatusText(code)
	if success {
		h.log.Infof("%v %v %v", callInfo, code, status)
		return
	}

	h.log.Errorf("%v %v %v", callInfo, code, status)
}

func (h *handler) getCallInfo(r *http.Request) string {
	info := " "
	val := r.Header.Get("X-Request-ID")
	if val != "" {
		info += fmt.Sprintf("%v(%v)", ridKey, val)
	}
	info += " " + r.Method
	info += " " + r.URL.Path
	return info
}

const (
	formFieldName = "password"
)
