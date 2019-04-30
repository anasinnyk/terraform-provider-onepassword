package onepassword

import (
	"crypto/rand"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"net/url"
	"regexp"
	"strings"
)

func emailValidate(i interface{}, k string) (s []string, es []error) {
	v, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	emailRegexp := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
			"[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
	if !emailRegexp.MatchString(v) {
		es = append(es, fmt.Errorf("%s is not email", k))
		return
	}
	return
}

func ToSnakeCase(str string) string {
	snake := regexp.MustCompile("(.)([A-Z][a-z]+)").ReplaceAllString(str, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
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

func orEmpty(f schema.SchemaValidateFunc) schema.SchemaValidateFunc {
	return func(i interface{}, k string) ([]string, []error) {
		v, ok := i.(string)
		if ok && v == "" {
			return nil, nil
		}
		return f(i, k)
	}
}

func fieldNumber() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(fmt.Sprintf("%x", b)), nil
}
