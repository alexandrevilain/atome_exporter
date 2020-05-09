.EXPORT_ALL_VARIABLES:

OBJC_DISABLE_INITIALIZE_FORK_SAFETY = YES

deployer-deps:
	pip install ansible --user
	pip install jmespath --user
deploy:
	GOOS=linux GOARCH=amd64 go build -o build/atome_exporter ./cmd/atome_exporter
	ansible-galaxy install -r ./deployment/requirements.yaml
	ansible-playbook -i ./deployment/environment/production/hosts --vault-password-file .ansible_vault ./deployment/playbook.yml