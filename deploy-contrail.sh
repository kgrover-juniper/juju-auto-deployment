#!/bin/bash

juju destroy-model contrail <<-EOF
y
EOF

juju add-model contrail && 

if [ $1 = "queens" ]; then
	filename=contrail-docker-bundle-queens.yaml
elif [ $1 = "ussuri" ]; then
	filename=contrail-docker-bundle-ussuri.yaml
else
	filename=contrail-docker-bundle-train.yaml
fi

sed -i "s/image-tag.*/image-tag: \"$2\"/g" $filename &&

juju deploy ./$filename
