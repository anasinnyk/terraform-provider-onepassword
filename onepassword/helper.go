package onepassword

import (
	"fmt"
	"net/url"
	"regexp"
)

func emailValidate(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(v) {
		es = append(es, fmt.Errorf("%s is not email", k))
		return
	}
	return
}

func urlValidate(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	_, err := url.ParseRequestURI(v)
	if err != nil {
		es = append(es, fmt.Errorf("%s is not an URL", v))
		return
	}
	return
}
