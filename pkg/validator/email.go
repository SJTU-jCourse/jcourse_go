package validator

import (
	"fmt"
	"regexp"
)

type EmailValidator interface {
	Validate(email string) bool
}

type CommonEmailValidator struct{}

func (v *CommonEmailValidator) Validate(email string) bool {
	re := regexp.MustCompile(`\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`)
	return re.MatchString(email)
}

type SuffixEmailValidator struct {
	suffixDomain string
}

func (v *SuffixEmailValidator) Validate(email string) bool {
	re := regexp.MustCompile(fmt.Sprintf(`\w+([-+.]\w+)*@%s`, v.suffixDomain))
	return re.MatchString(email)
}

func NewEmailValidator(suffixDomain string) EmailValidator {
	if suffixDomain == "" {
		return &CommonEmailValidator{}
	}
	return &SuffixEmailValidator{suffixDomain: suffixDomain}
}
