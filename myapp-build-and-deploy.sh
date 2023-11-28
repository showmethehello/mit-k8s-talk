#!/bin/bash
REGISTRY=wokstok #change this to your hub.docker username
APP=mit-myapp
TIME=$(date +%H%M%S)
TAG=${1:-$TIME}
IMAGE=$REGISTRY/$APP:$TAG

echo "Building $IMAGE"
docker build . --build-arg VERSION=$TAG -f myapp-Dockerfile -t $IMAGE

# push the image from local docker to the cloud
# registry. Optional if using "kind load docker-image" below
# echo "Copying $IMAGE to cloud registry"
# docker push $IMAGE

# Use kind load docker-image to copy the image from docker to
# your kind kubernetes cluster to save time downloading it
echo "Copying $IMAGE to kind cluster"
kind load docker-image $IMAGE --name mit

echo "Applying deployment to kubenetes"
cat "myapp-deployment.yaml" | sed "s|MYAPP_IMAGE|${IMAGE}|g" | kubectl apply -f -
