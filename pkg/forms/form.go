package forms

import (
	"fmt"
	"net/url"
	//"strings"
	"unicode/utf8"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		if f.Get(field) == "" {
			f.Errors.Add(field, "This field is required.")
		}
	}
}

func (f *Form) MaxLength(field string, d int) {
	value := f.Get(field)
	if value == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field cannot be more than %d characters.", d))
	}
}
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)
	if value == "" {
		return
	}
	for _, opts := range opts {
		if value == opts {
			return
		}
	}
	f.Errors.Add(field, "this field is invalid")
}
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
	}
