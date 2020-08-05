package onepassword

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func stringDiag() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		val, ok := v.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Value is not a string",
				Detail:        fmt.Sprintf("Value is not a string (type = %T)", val),
				AttributePath: path,
			})
		}
		return diags
	}
}

func emailValidateDiag() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		diags := stringDiag()(v, path)
		val, _ := v.(string)
		
		emailRegexp := regexp.MustCompile(
			"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
				"[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
		)
		if len(diags) == 0 && !emailRegexp.MatchString(val) {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Value is not email",
				Detail:   	   fmt.Sprintf("%s is not email", val),
				AttributePath: path,
			})
		}

		return diags
	}
}

func urlValidateDiag() schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		diags := stringDiag()(v, path)
		val, _ := v.(string)
		if len(diags) == 0 {
			_, err := url.ParseRequestURI(val)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Value is not URL",
					Detail:        fmt.Sprintf("%s is not an URL", val),
					AttributePath: path,
				})
			}
		}
		return diags
	}
}

func stringInSliceDiag(ss []string, empty bool) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		diags := stringDiag()(v, path)
		val, _ := v.(string)
		if len(diags) == 0 {
			if (!empty || val != "") && !stringInSlice(val, ss) {
				diags = append(diags, diag.Diagnostic{
					Severity:      diag.Error,
					Summary:       "Value has incorect value",
					Detail:        fmt.Sprintf("%s one from next list (%s)", val, strings.Join(ss,",")),
					AttributePath: path,
				})
			}
		}
		return diags
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func fieldNumber() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(fmt.Sprintf("%x", b)), nil
}
