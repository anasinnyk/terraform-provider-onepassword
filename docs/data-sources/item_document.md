# onepassword_item_document

This resource can load any document from 1password.

## Example Usage

```hcl
data "onepassword_item_document" "this" {
    name = "some-document-from-vault"
}
```

## Argument Reference

* `name` - (Required) your document title.
* `binary` - (Optional) set to true if your document is binary. Default: false.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - document id.
* `content` - document content. If `binary` is true, this is a base64 encoded string.
