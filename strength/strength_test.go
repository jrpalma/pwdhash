package strength

import (
	"testing"
)

func defaultStrength() PasswordStrength {
	rules := PasswordStrength{}
	rules.MinDigits = 1
	rules.MinLowerCase = 2
	rules.MinUpperCase = 2
	rules.MinDigits = 1
	rules.MinSpecial = 1
	rules.MinLength = 8
	rules.MaxLength = 10
	return rules
}

func TestCheck_Success(t *testing.T) {
	pwd := "OneTwo3!"
	rules := defaultStrength()
	pass := Check(rules, pwd)
	if !pass {
		t.Errorf("Password: %v, Failed with: %+v", pwd, rules)
	}
}

func TestCheck_Failure(t *testing.T) {
	rules := defaultStrength()

	pwd := "OneTwo3X"
	pass := Check(rules, pwd)
	if pass {
		t.Errorf("Password: %v, Should fail with: %+v", pwd, rules)
	}

	pwd = "OneTw3!"
	pass = Check(rules, pwd)
	if pass {
		t.Errorf("Password: %v, Should fail with: %+v", pwd, rules)
	}

	pwd = "OneTwoX!"
	pass = Check(rules, pwd)
	if pass {
		t.Errorf("Password: %v, Should fail with: %+v", pwd, rules)
	}

	pwd = "Onetwo3!"
	pass = Check(rules, pwd)
	if pass {
		t.Errorf("Password: %v, Should fail with: %+v", pwd, rules)
	}

	pwd = "ONETWo3!"
	pass = Check(rules, pwd)
	if pass {
		t.Errorf("Password: %v, Should fail with: %+v", pwd, rules)
	}

	pwd = "ONETWo3!XXXX"
	pass = Check(rules, pwd)
	if pass {
		t.Errorf("Password: %v, Should fail with: %+v", pwd, rules)
	}
}
