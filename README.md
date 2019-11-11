# signup

This component allows the signup of a new organization on the **Nalej™** platform.
It creates a new organization domain, and creates the two first users on the organization:

* Nalej Admin: Ops user, allowed to manage the cluster and the organization, create users, cordon clusters. Has super-cow powers.
* Owner: First organization user. Able to run instances, upload descriptors, etc.

## Getting Started

### Prerequisites

* [system-model](https://github.com/nalej/system-model)
* [user-manager](https://github.com/nalej/user-manager)

### Build and compile

To build and compile this repository use the provided Makefile:

```shell script
make all
```

This operation generates the binaries for this repo, download dependencies,
run existing tests and generate ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```shell script
make test
```

## Manually testing client/server certificate

* Make everything: `make`
* Generate a client and secret certificate pairs. (Requires CFSSL to work: `brew install cfssl`)
* Create docker image and deploy it in Kubernetes (with its dependent services)
* Because the way certificates work, it's necessary to override the domain in `/etc/hosts` to make it work.
You have to add your minikube ip (`minikube ip`) as an entry in your `/etc/hosts` pointing to `signup.nalej`.
* Retrieve the port of the signup service running on minikube:
  ```shell script
  service_url=$(minikube service --namespace=nalej --url signup); echo ${service_url#http://*:}
  ```

* Launch the signup-cli with the required parameters replacing the placeholders with its correct values (paths are absolute): 
  
  ```shell script
  ./bin/signup-cli signup --signupAddress=signup.nalej:SERVICE_PORT --orgName=test --ownerEmail=test -- ownerName=test --ownerPassword=test --caPath=CLUSTER_CA_PATH --clientCertPath=CLIENT_CERT_PATH --clientKeyPath=CLIENT_KEY_PATH
  ```


### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```shell script
make dep
```

To have all dependencies up-to-date run:

```shell script
dep ensure -update -v
```

## User client interface

After compiling the project, the `signup-cli` will be available. With it, you will be able to create and manage organizations
in a cluster.

Before launching the service you need to generate the server and client certificates.
To do it go to `scripts` and execute `setup_service_certificates.sh`.
When it is done, you'll have the client certificate and its key in that directory.

Example:

```shell script
./bin/signup-cli signup --signupAddress=signup.nalej:SERVICE_PORT --orgName=test --ownerEmail=test --ownerName=test --ownerPassword=test --caPath=CLUSTER_CA_PATH --clientCertPath=CLIENT_CERT_PATH --clientKeyPath=CLIENT_KEY_PATH
```

## Known Issues

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting
pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the
[tags on this repository](https://github.com/nalej/signup/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/signup/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
