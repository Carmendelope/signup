#!/bin/bash

function abort() {
	echo "Aborting!!"
	exit 1
}

function check_cfssl {
	command -v cfssl >/dev/null 2>&1 || { echo "cfssl is not installed. Please install it with \"brew install cfssl\"." >&2; abort; }
}

function check_yaml_exists {
	if [ ! -f "./bin/yaml/mngtcluster/signup.service.yaml" ]; then
		echo "signup.service.yaml not found. Please generate them using \"make\" or \"make yaml\""
		abort
	fi
}

function get_service_cluster_ip {
	check_yaml_exists
	service_cluster_ip=$(kubectl get -f ./bin/yaml/mngtcluster/signup.service.yaml -o jsonpath={.spec.clusterIP})
}

function get_k8s_config_ca_certificate {
	kubectl config view -o jsonpath='{.clusters[?(@.name=="minikube")].cluster.certificate-authority-data}' --raw | base64 --decode > cluster_ca.crt
} 

function get_client_secret {
	read -s -p "Please enter the client secret: " input_secret
	echo ""
}

function clear_old_files {
	while true; do
		echo "WARNING!! This will remove any other certificate created previously."
		read -p "Do you want to continue? [y/N]: " yn
		case $yn in
			[Yy]* ) rm -rf {server.csr,server-key.pem,server.crt,client.csr,client-key.pem,client.crt,cluster_ca.crt}; break;;
			[Nn]* ) abort;;
			* ) abort;;
		esac
	done
}

function get_server_csr {
	local cluster_ip=$1

	cat <<-EOF | cfssl genkey - | cfssljson -bare server
	{
		"hosts": [
			"signup.nalej.svc.cluster.local",
			"signup.nalej",
                        "signup.nalej.cluster.local",
                        "nalej.cluster.local",
                        "login.nalej.cluster.local",
                        "api.nalej.cluster.local",
			"$cluster_ip"
		],
		"CN": "signup.nalej.cluster.local",
		"key": {
			"algo": "ecdsa",
			"size": 256
		}
	}
	EOF
}

function get_client_csr {
	local client_secret=$1

	cat <<-EOF | cfssl genkey - | cfssljson -bare client
	{
		"hosts": [
			""
		],
		"CN": "${client_secret}",
		"key": {
			"algo": "ecdsa",
			"size": 256
		}
	}
	EOF
}

function request_k8s_signing {
	local csr_type=$1

	kubectl delete csr signup.nalej-${csr_type}
	cat <<-EOF | kubectl create -f -
	apiVersion: certificates.k8s.io/v1beta1
	kind: CertificateSigningRequest
	metadata:
	  name: signup.nalej-${csr_type}
	  namespace: nalej
	spec:
	  groups:
	  - system:authenticated
	  request: $(cat ${csr_type}.csr | base64 | tr -d '\n')
	  usages:
	  - digital signature
	  - key encipherment
	  - ${csr_type} auth
	EOF
	kubectl certificate approve signup.nalej-${csr_type}
	kubectl get csr signup.nalej-${csr_type} -o jsonpath='{.status.certificate}' | base64 --decode > ${csr_type}.crt
}

function create_server_k8s_secret {
	kubectl --namespace=nalej delete secret signup-server-tls
        kubectl --namespace=nalej create secret tls signup-server-tls --key server-key.pem --cert server.crt
}

function create_client_k8s_secret {
	local client_secret=$1

	kubectl --namespace=nalej delete secret signup-client-secret
	cat <<-EOF | kubectl create -f -
	apiVersion: v1
	kind: Secret
	metadata:
	  name: signup-client-secret
	  namespace: nalej
	data:
	  secret: $(echo -n "${client_secret}" | base64)
	EOF
}

# Prepare
clear
check_cfssl
clear_old_files
get_k8s_config_ca_certificate

# Server certificate
echo "######################"
echo "Generating server part"
echo "######################"
get_service_cluster_ip  # Stored in service_cluster_ip
get_server_csr $service_cluster_ip
request_k8s_signing "server"
create_server_k8s_secret
rm -rf server{.csr,-key.pem,.crt}

# Client certificate
echo ""
echo ""
echo "######################"
echo "Generating client part"
echo "######################"
get_client_secret  # Stored in input_secret
get_client_csr $input_secret
request_k8s_signing "client"
create_client_k8s_secret $input_secret
rm -rf client.csr

echo ""
echo "SETUP FINISHED!!"
echo "Now you can use client.crt and client-key.pem to connect to signup service from the CLI."
echo "The CA is cluster_ca.crt."

exit 0
