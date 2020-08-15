# Provider

Terraform provider for 1password usage with your infrastructure, for example you can share password from your admin panel via some vault in you 1password company account. This provider is based on 1Password CLI client version 1.4.0, but you can rewrite it by env variable `OP_VERSION`.

## Example Usage

```hcl
provider "onepassword" {
    email      = "john.smith@example.com"
    password   = "super secret master password"
    secret_key = "A3-XXXXXX-XXXXXXX-XXXXX-XXXXX-XXXXX-XXXXX"
    subdomain  = "company"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Optional) your email address in 1password or via env variable `OP_EMAIL`.
* `password` - (Optional) your master password from 1password or via env variable `OP_PASSWORD`.
* `secret_key` - (Optional) secret key which you can download after registration or via env variable `OP_SECRET_KEY`.
* `subdomain` - (Optional) If you use corporate account you must fill subdomain form your 1password site. Defaults to `my` or via env variable `OP_SUBDOMAIN`.

If `email`, `password` and `secret_key` is not set through the arguments or env variables, then the env variable `OP_SESSION_<subdomain>` is checked for existence. If set it will be assumed to be a valid session token and used while executing the `op` commands. Note that any dash `-` character within `subdomain` will be substituted upon `OP_SESSION_<subdomain>` env variable evaluation (e.g, if `subdomain=team-foo`, `OP_SESSION_team_foo` will be looked up).
