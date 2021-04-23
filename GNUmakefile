default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	./scripts/runvault.sh &
	VAULT_ADDR=http://localhost:8200 TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
	./scripts/stopvault.sh
