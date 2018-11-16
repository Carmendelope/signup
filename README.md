# signup
SignUp component

## Generating server and client certificates for service authentication
Before launching the service you need to generate the server and client certificates.

To do it go to `scripts` and execute `setup_service_certificates.sh`.
When its done, you'll have the client certificate and its key in that directory.

## Manually testing client/server certificate
* Make everything: `make`
* Generate a client and secret certificate pairs. (Requires CFSSL to work: `brew install cfssl`)
* Create docker image and deploy it in Kubernetes (with its dependent services)
* Because the way certificates work, it's necessary to override the domain in `/etc/hosts` to make it work. You have to add you minikube ip (`minikube ip`) as an entry in your `/etc/hosts` pointing to `signup.nalej`. Once we have the ingress up and running we'll have a better solution for this.
* Retrieve the port of the signup service running on minikube:
  ```
  service_url=$(minikube service --namespace=nalej --url signup); echo ${service_url#http://*:}
  ```

* Launch the signup-cli with the required parameters replacing the placeholders with its correct values (paths are absolute): 
  ```
  ./bin/signup-cli signup --signupAddress=signup.nalej:SERVICE_PORT --orgName=test --ownerEmail=test -- ownerName=test --ownerPassword=test --caPath=CLUSTER_CA_PATH --clientCertPath=CLIENT_CERT_PATH --clientKeyPath=CLIENT_KEY_PATH
  ```
