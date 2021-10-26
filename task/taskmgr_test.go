package task

import (
	"fmt"
	"testing"
	"time"

	"github.com/jrpalma/pwdhash/config"
)

func TestManager_ShutdownTwice(t *testing.T) {
	conf := config.Config{}
	mgr := NewManager(conf)

	res := mgr.Shutdown()
	if res.Code != 200 {
		t.Errorf("Shutdown should return 200, not %v", res.Code)
	}

	res = mgr.Shutdown()
	if res.Code != 500 {
		t.Errorf("Shutdown should return 500, not %v", res.Code)
	}
}

func TestManager_Stats(t *testing.T) {
	conf := config.Config{}
	mgr := NewManager(conf)

	_, res := mgr.Stats()
	if res.Code != 200 {
		t.Errorf("Shutdown should return 200, not %v", res.Code)
	}

	res = mgr.Shutdown()
	if res.Code != 200 {
		t.Errorf("Shutdown should return 200, not %v", res.Code)
	}

	_, res = mgr.Stats()
	if res.Code != 500 {
		t.Errorf("Stats should return 500, not %v", res.Code)
	}
}

func TestManager_Check(t *testing.T) {
	conf := config.Config{}
	mgr := NewManager(conf)

	res := mgr.Check("badInteger")
	if res.Code != 400 {
		t.Errorf("Check should return 400, not %v", res.Code)
	}

	res = mgr.Check("0")
	if res.Code != 404 {
		t.Errorf("Stats should return 404, not %v", res.Code)
	}

	res = mgr.Shutdown()
	if res.Code != 200 {
		t.Errorf("Shutdown should return 200, not %v", res.Code)
	}

	res = mgr.Check("0")
	if res.Code != 500 {
		t.Errorf("Stats should return 404, not %v", res.Code)
	}
}

func TestManager_ShutdownNewTask(t *testing.T) {
	conf := config.Config{}

	mgr := NewManager(conf)
	res := mgr.Shutdown()
	if res.Code != 200 {
		t.Errorf("Shutdown should return 200, not %v", res.Code)
	}

	res = mgr.NewTask("pass")
	if res.Code != 500 {
		t.Errorf("Stats should return 500, not %v", res.Code)
	}
}

func TestManager_NewTaskPasswordStrength(t *testing.T) {
	conf := config.Config{}
	conf.CheckPasswordStrength = true
	conf.PasswordStrength.MinLength = 8

	mgr := NewManager(conf)

	res := mgr.NewTask("pass")
	if res.Code != 400 {
		t.Errorf("Stats should return 400, not %v", res.Code)
	}
}

func TestManager_WaitForTasks(t *testing.T) {
	// NOTE: This test might not be determistic
	// on a system that is low on resources.

	conf := config.Config{}

	conf.MaxTaskSeconds = 1
	mgr := NewManager(conf)
	var hashIDs []string

	for i := 0; i < 50; i++ {
		pass := fmt.Sprintf("pass%v", i)
		res := mgr.NewTask(pass)
		if res.Code != 201 {
			t.Errorf("Stats should return 201, not %v", res.Code)
			continue
		}
		hashIDs = append(hashIDs, res.Message)
	}

	// No task should have completed
	for _, hashID := range hashIDs {
		res := mgr.Check(hashID)
		if res.Code != 503 {
			t.Errorf("Stats should return 503, not %v: %v", res.Code, hashID)
		}
	}

	// Two seconds should be enough time so we can check
	// For finished tasks. However, this might not work
	// on a system low on resources.
	time.Sleep(time.Second * 2)

	for _, hashID := range hashIDs {
		res := mgr.Check(hashID)
		if res.Code != 200 {
			t.Errorf("Stats should return 200, not %v: %v", res.Code, hashID)
		}
	}

	mgr.WaitForPendingTasks()
}
