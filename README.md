# Golang juju-auto-deployment
Required files:
1. juju-golang-deployment.go
2. deploy-contrail.sh
3. Atleast one bundle yaml file: contrail-docker-bundle-queens.yaml or contrail-docker-bundle-train.yaml 

> Deploy from inside tf-charms

Command line execution:
-----------------------
go run juju-golang-deployment.go 'openstack-version' 'contrail-build'


'openstack-version'
 queens
 or
 train

'contrail-build'
 2008.12 
 
 Output:
 -------
 result.txt
 #Comment on/off write_result() in main() of juju-golang-deployment.go

