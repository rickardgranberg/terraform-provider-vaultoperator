disable_mlock = true

listener "tcp" {
    address = "[::]:0"
    tls_disable = "{{ .DisableTLS }}"
    tls_cert_file = "{{ .CertFile }}"
    tls_key_file = "{{ .KeyFile }}"
}

storage "inmem" {}
