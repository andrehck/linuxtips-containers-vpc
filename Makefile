# Definindo variáveis
TF_INIT_FLAGS = -backend-config=env/dev/backend.tfvars
TF_APPLY_FLAGS = --auto-approve -var-file=env/dev/terraform.tfvars
TF_DESTROY_FLAGS = --auto-approve -var-file=env/dev/terraform.tfvars

# Inicializa o Terraform
init:
	terraform init $(TF_INIT_FLAGS)

# Aplica a configuração do Terraform
apply: init
	terraform apply $(TF_APPLY_FLAGS)

# Formata os arquivos Terraform
fmt:
	terraform fmt

# Destrói a infraestrutura gerenciada pelo Terraform
destroy: init
	terraform destroy $(TF_DESTROY_FLAGS)
