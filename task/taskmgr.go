package task

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/jrpalma/pwdhash/config"
	"github.com/jrpalma/pwdhash/strength"
)

// NewManager Creates a new task manager with the given configuration.
func NewManager(config config.Config) *Manager {
	tm := &Manager{config: config}
	tm.tasks = make(map[uint64]*task)
	return tm
}

// Result The task manager operation result.
type Result struct {
	// Code The HTTP status code to return
	Code int
	// Message The HTTP message to use in a response.
	Message string
}

// Stats The task manager statistics
type Stats struct {
	// Total The number of completed tasks.
	Total uint64 `json:"total"`
	// Average The average microsecons it has taken to process all completed tasks.
	Average uint64 `json:"average"`
}

// Manager A task manager is charage of all the password hash operations.
type Manager struct {
	taskID         uint64
	taskRuntime    time.Duration
	completedTasks uint64

	done   bool
	config config.Config
	tasks  map[uint64]*task
	mutex  sync.Mutex
	wg     sync.WaitGroup
}

// Shutdown Shutdown the task manager. A call to WaitForPendingTasks is expected
// in order to wait for pending tasks after this call. This call might fail if a
// shutdown is pending.
func (tm *Manager) Shutdown() Result {
	result := Result{}

	if tm.done {
		result.Message = shutdownMsg
		result.Code = 500
		return result
	}

	tm.done = true
	result.Code = 200

	return result
}

// WaitForPendingTasks Waits for all pending tasks.
func (tm *Manager) WaitForPendingTasks() {
	tm.wg.Wait()
}

// Stats Returns the statisc object. This call might fail
// it the shutdown is pending.
func (tm *Manager) Stats() (Stats, Result) {
	stats := Stats{}
	result := Result{}

	if tm.done {
		result.Message = shutdownMsg
		result.Code = 500
		return stats, result
	}

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	stats.Total = tm.completedTasks
	avg := float64(tm.taskRuntime*time.Nanosecond) / float64(tm.completedTasks)
	stats.Average = uint64(avg)
	result.Code = 200

	return stats, result
}

// NewTask Creates a new password hash task. This call
// might fail if shutdown is pending.
func (tm *Manager) NewTask(pwd string) Result {
	result := Result{}

	if tm.done {
		result.Message = shutdownMsg
		result.Code = 500
		return result
	}

	if tm.config.CheckPasswordStrength {
		strongPassword := strength.Check(tm.config.PasswordStrength, pwd)
		if !strongPassword {
			result.Message = fmt.Sprintf("Password is too weak")
			result.Code = 400
			return result
		}
	}

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.taskID++
	task := &task{ID: tm.taskID, Password: pwd}
	tm.tasks[task.ID] = task

	tm.wg.Add(1)
	go tm.runTask(task)

	result.Code = 201
	result.Message = strconv.FormatUint(tm.taskID, 10)
	return result
}

// Check Checks if a new password hash task has completed. This call
// might fail if shutdown is pending.
func (tm *Manager) Check(hashID string) Result {
	result := Result{}

	if tm.done {
		result.Message = shutdownMsg
		result.Code = 500
		return result
	}

	id, err := strconv.ParseUint(hashID, 10, 64)
	if err != nil {
		result.Message = fmt.Sprintf("Invalid hash ID %v", hashID)
		result.Code = 400
		return result
	}

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	task, exists := tm.tasks[id]
	if !exists {
		result.Message = fmt.Sprintf("No such hash ID %v", hashID)
		result.Code = 404
		return result
	}

	if !task.Done {
		result.Code = 503
		return result
	}

	result.Code = 200
	result.Message = task.Hash

	return result
}

func (tm *Manager) runTask(task *task) {
	defer tm.wg.Done()

	start := time.Now()
	hash := sha512.Sum512([]byte(task.Password))
	data := base64.StdEncoding.EncodeToString(hash[:])
	time.Sleep(time.Second * time.Duration(tm.config.MaxTaskSeconds))
	taskDuration := time.Since(start)

	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.taskRuntime += taskDuration
	tm.completedTasks++

	task.Hash = data
	task.Done = true
}

const shutdownMsg = "Service is shutting down"

type task struct {
	ID       uint64
	Done     bool
	Hash     string
	Password string
}
