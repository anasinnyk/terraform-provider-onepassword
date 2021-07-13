# onepassword_item_document

This resource can create any document for 1password.

## Example Usage

```hcl
resource "onepassword_item_document" "this" {
  name      = "document-name"
  vault     = var.vault_id
  file_path = "${path.module}/test.txt"
}
```

## Argument Reference

* `name` - (Required) your document title.
* `field_path` - (Required) path to your document, which will be upload to 1password.
* `binary` - (Optional) set to true if your document is binary. Default: false.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - document id.
* `content` - document content. If `binary` is true, this is a base64 encoded string.
