package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/jrpalma/pwdhash/logs"
	"github.com/jrpalma/pwdhash/strength"
)

// Config Represents the service configuration.
type Config struct {

	// LogLevel The service log level.
	LogLevel logs.LogLevel

	// LogDestinatio The destination for the logs.
	LogDestination logs.Destination

	// LogFile The path to the log file where logs will be written.
	LogFile string

	// CheckPasswordStrength Flag used to enable the password strength checks.
	CheckPasswordStrength bool

	// PasswordStrength Rules used to check the password strength.
	PasswordStrength strength.PasswordStrength

	// MaxTaskSeconds The maximum number of seconds the hash task will take.
	MaxTaskSeconds uint

	// ServerAddress The server address to listen on. For example: ":80"
	ServerAddress string
}

// OpenFile Opens or creates a configuration file. If the file exist,
// the file is opened and loaded. If the file does not exis, the file
// is created with the default values and saved. The default values
// used are: Log level WARN, Log Destination STDERR, and checks for
// password strength. The default password strength has a mininum
// length of 8 plus a minimum 1 upper case, 1 lower case, 1 digit,
// and 1 special character. The default runtime is 5 seconds.
func (c *Config) OpenFile(filePath string) error {

	content, err := ioutil.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		if c.LogLevel == "" {
			c.LogLevel = logs.WARN
		}
		if c.LogDestination == "" {
			c.LogDestination = logs.STDERR
		}
		if c.MaxTaskSeconds == 0 {
			c.MaxTaskSeconds = 5
		}
		if c.ServerAddress == "" {
			c.ServerAddress = ":8080"
		}
		if c.PasswordStrength.MaxLength == 0 {
			c.PasswordStrength.MaxLength = 50
		}

		c.CheckPasswordStrength = true
		c.PasswordStrength.MinDigits = 1
		c.PasswordStrength.MinUpperCase = 1
		c.PasswordStrength.MinLowerCase = 1
		c.PasswordStrength.MinDigits = 1
		c.PasswordStrength.MinSpecial = 1
		c.PasswordStrength.MinLength = 8

		return c.SaveFile(filePath)
	}

	err = json.Unmarshal(content, c)
	if err != nil {
		return err
	}

	return nil
}

// SaveFile Saves the JSON configuration to the filePath.
func (c *Config) SaveFile(filePath string) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, bytes, 0644)

	return err
}

// CreateLogger Creates a logger from this configuration.
func (c *Config) CreateLogger() (logs.Logger, error) {
	if c.LogDestination == logs.FILE {
		return logs.NewFileLogger(c.LogFile, c.LogLevel)
	}
	return logs.NewStreamLogger(c.LogDestination, c.LogLevel)
}
