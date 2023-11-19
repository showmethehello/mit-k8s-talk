#!/bin/bash
REGISTRY=registry.mitmullingar.com
APP=myapp
TAG=$(date +%H%M%S)
IMAGE=$REGISTRY/$APP:$TAG

echo "Building $IMAGE"
docker build . --build-arg VERSION=$TAG -f myapp-Dockerfile -t $IMAGE

echo "Copying $IMAGE to kind cluster"
kind load docker-image $IMAGE --name mit

echo "Applying deployment to kubenetes"
cat "myapp-deployment.yaml" | sed "s/MYAPP_TAG/${TAG}/g" | kubectl apply -f -
