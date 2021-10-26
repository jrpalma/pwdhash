package logs

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestNewStreamLogger_Fail(t *testing.T) {
	// Test invalid destination
	_, err := NewStreamLogger(Destination("INVALID"), INFO)
	if err == nil {
		t.Errorf("NewStreamLogger should fail with invalid destination")
	}

	// Test invalid log level
	_, err = NewStreamLogger(STDOUT, LogLevel("INVALID"))
	if err == nil {
		t.Errorf("NewStreamLogger should fail with invalid log level")
	}
}

func TestNewStreamLogger_Success(t *testing.T) {
	_, err := NewStreamLogger(STDOUT, DEBUG)
	if err != nil {
		t.Errorf("NewStreamLogger failed: %v", err)
	}
	_, err = NewStreamLogger(STDERR, DEBUG)
	if err != nil {
		t.Errorf("NewStreamLogger failed: %v", err)
	}
}

func TestNewStreamLogger_Errorf(t *testing.T) {

	log, err := NewStreamLogger(STDOUT, ERROR)
	if err != nil {
		t.Errorf("NewStreamLogger failed: %v", err)
	}

	// Access the internal implementation to test
	buff := &bytes.Buffer{}
	imp := log.(*logger)
	imp.file = buff

	log.Errorf("Error")
	log.Warnf("Warn")
	log.Infof("Info")
	log.Debugf("Debug")

	lines := strings.Split(buff.String(), "\n")
	if len(lines) != 2 {
		t.Errorf("Invalid number of lines: %v", len(lines))
		return
	}
	if !strings.Contains(lines[0], "Error") {
		t.Errorf("Error should be printed to the logs")
	}
}

func TestNewStreamLogger_Warnf(t *testing.T) {

	log, err := NewStreamLogger(STDOUT, WARN)
	if err != nil {
		t.Errorf("NewStreamLogger failed: %v", err)
	}

	// Access the internal implementation to test
	buff := &bytes.Buffer{}
	imp := log.(*logger)
	imp.file = buff

	log.Errorf("Error")
	log.Warnf("Warn")
	log.Infof("Info")
	log.Debugf("Debug")

	lines := strings.Split(buff.String(), "\n")
	if len(lines) != 3 {
		t.Errorf("Invalid number of lines: %v", len(lines))
		return
	}
	if !strings.Contains(lines[0], "Error") {
		t.Errorf("Error should be printed to the logs")
	}
	if !strings.Contains(lines[1], "Warn") {
		t.Errorf("Warning should be printed to the logs")
	}
}

func TestNewStreamLogger_Infof(t *testing.T) {

	log, err := NewStreamLogger(STDOUT, INFO)
	if err != nil {
		t.Errorf("NewStreamLogger failed: %v", err)
	}

	// Access the internal implementation to test
	buff := &bytes.Buffer{}
	imp := log.(*logger)
	imp.file = buff

	log.Errorf("Error")
	log.Warnf("Warn")
	log.Infof("Info")
	log.Debugf("Debug")

	lines := strings.Split(buff.String(), "\n")
	if len(lines) != 4 {
		t.Errorf("Invalid number of lines: %v", len(lines))
		return
	}
	if !strings.Contains(lines[0], "Error") {
		t.Errorf("Error should be printed to the logs")
	}
	if !strings.Contains(lines[1], "Warn") {
		t.Errorf("Warning should be printed to the logs")
	}
	if !strings.Contains(lines[2], "Info") {
		t.Errorf("Info should be printed to the logs")
	}
}

func TestNewStreamLogger_Debugf(t *testing.T) {

	log, err := NewStreamLogger(STDOUT, DEBUG)
	if err != nil {
		t.Errorf("NewStreamLogger failed: %v", err)
	}

	// Access the internal implementation to test
	buff := &bytes.Buffer{}
	imp := log.(*logger)
	imp.file = buff

	log.Errorf("Error")
	log.Warnf("Warn")
	log.Infof("Info")
	log.Debugf("Debug")

	lines := strings.Split(buff.String(), "\n")
	if len(lines) != 5 {
		t.Errorf("Invalid number of lines: %v", len(lines))
		return
	}
	if !strings.Contains(lines[0], "Error") {
		t.Errorf("Error should be printed to the logs")
	}
	if !strings.Contains(lines[1], "Warn") {
		t.Errorf("Warning should be printed to the logs")
	}
	if !strings.Contains(lines[2], "Info") {
		t.Errorf("Info should be printed to the logs")
	}
	if !strings.Contains(lines[3], "Debug") {
		t.Errorf("Debug should be printed to the logs")
	}
}

func TestNewFileLogger_Fail(t *testing.T) {

	_, err := NewFileLogger("./bad/file", ERROR)
	if err == nil {
		t.Errorf("NewFileLogger should have failed with invalid path")
	}

	_, err = NewFileLogger("./bad/file", LogLevel("INVALID"))
	if err == nil {
		t.Errorf("NewFileLogger should have failed with invalid log level")
	}
}

func TestNewFileLogger_Success(t *testing.T) {
	tmpFile := "./tmp"
	defer os.Remove(tmpFile)

	_, err := NewFileLogger(tmpFile, ERROR)
	if err != nil {
		t.Errorf("NewFileLogger failed: %v", err)
	}

}
