data "local_file" "pgp_key" {
  for_each = toset(["one", "two", "three", "four", "five"])
  filename = "${path.module}/${each.key}.gpg"
}
data "local_file" "root_token_pgp_key" {
  filename = "${path.module}/root.gpg"
}
resource "vaultoperator_init" "example" {
  secret_shares      = 5
  secret_threshold   = 3
  pgp_keys           = data.local_file.pgp_key.*.content_base64
  root_token_pgp_key = data.local_file.root_token_pgp_key.content_base64
}
