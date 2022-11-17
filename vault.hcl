disable_mlock = true

listener "tcp" {
    tls_disable = 1
    address = "[::]:0"
}

storage "inmem" {}
