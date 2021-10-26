package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/jrpalma/pwdhash/logs"
)

func defaultConfig() *Config {
	conf := &Config{}
	conf.LogLevel = logs.WARN
	conf.LogDestination = logs.STDERR
	conf.CheckPasswordStrength = true
	conf.PasswordStrength.MinDigits = 1
	conf.PasswordStrength.MinUpperCase = 1
	conf.PasswordStrength.MinLowerCase = 1
	conf.PasswordStrength.MinDigits = 1
	conf.PasswordStrength.MinSpecial = 1
	conf.PasswordStrength.MinLength = 8
	conf.PasswordStrength.MaxLength = 50
	conf.MaxTaskSeconds = 5
	conf.ServerAddress = ":8080"
	return conf
}

func TestOpenFile_Defaults(t *testing.T) {
	file := "./OpenFile_Defaults"
	conf := &Config{}
	defaults := defaultConfig()

	defer os.Remove(file)

	err := conf.OpenFile(file)
	if err != nil {
		t.Errorf("OpenFile failed: %v", err)
		return
	}

	osFile, err := os.Open(file)
	if err != nil {
		t.Errorf("OpenFile failed to create file: %v", err)
		return
	}

	defer osFile.Close()
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("Failed to open file: %v", err)
		return
	}

	tmp := &Config{}
	err = json.Unmarshal(bytes, tmp)
	if err != nil {
		t.Errorf("Failed to unmarshal file: %v", err)
		return
	}

	if !reflect.DeepEqual(conf, tmp) {
		t.Errorf("Expected config: %+v. Got %+v", conf, tmp)
	}
	if !reflect.DeepEqual(defaults, tmp) {
		t.Errorf("Expected config: %+v. Got %+v", defaults, tmp)
	}
}

func TestOpenFile_ExistingFile(t *testing.T) {
	file := "./OpenFile_ExistingFile"
	defaults := defaultConfig()
	// Change something to compare
	defaults.LogLevel = logs.INFO

	defer os.Remove(file)

	err := defaults.SaveFile(file)
	if err != nil {
		t.Errorf("SaveFile failed: %v", err)
		return
	}

	tmp := &Config{}
	err = tmp.OpenFile(file)
	if err != nil {
		t.Errorf("OpenFile failed to create file: %v", err)
		return
	}

	if !reflect.DeepEqual(defaults, tmp) {
		t.Errorf("Expected config: %+v. Got %+v", defaults, tmp)
	}
}

func TestSaveFile_Fail(t *testing.T) {
	file := "./bad/SaveFile_Fail"
	defaults := defaultConfig()
	defer os.Remove(file)

	err := defaults.SaveFile(file)
	if err == nil {
		t.Errorf("SaveFile should fail")
	}
}

func TestOpenFile_UnmarshalFail(t *testing.T) {
	file := "./OpenFile_UnmarshalFail"
	defaults := defaultConfig()
	defer os.Remove(file)

	err := ioutil.WriteFile(file, []byte("{badJSON}"), 0644)
	if err != nil {
		t.Errorf("WriteFile failed: %v", err)
		return
	}

	err = defaults.OpenFile(file)
	if err == nil {
		t.Errorf("OpenFile should fail")
	}
}

func TestCreateLogger_Success(t *testing.T) {
	file := "./CreateLogger_Success"
	log := "./LOG_FILE"
	conf := defaultConfig()

	defer os.Remove(file)
	defer os.Remove(log)

	err := conf.OpenFile(file)
	if err != nil {
		t.Errorf("OpenFile failed: %v", err)
		return
	}

	// Use STDERR
	_, err = conf.CreateLogger()
	if err != nil {
		t.Errorf("CreateLogger failed: %v", err)
	}

	// Use FILE
	conf.LogDestination = logs.FILE
	conf.LogFile = log
	_, err = conf.CreateLogger()
	if err != nil {
		t.Errorf("CreateLogger failed: %v", err)
	}
}
