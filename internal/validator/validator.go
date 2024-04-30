package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FiledErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FiledErrors) == 0
}

func (v *Validator) AddFiledError(key, message string) {
	if v.FiledErrors == nil {
		v.FiledErrors = make(map[string]string)
	}

	if _, exists := v.FiledErrors[key]; !exists {
		v.FiledErrors[key] = message
	}
}

func (v *Validator) CheckFiled(ok bool, key, message string) {
	if !ok {
		v.AddFiledError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
