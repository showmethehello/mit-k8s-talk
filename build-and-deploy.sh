#!/bin/bash
REGISTRY=registry.mitmullingar.com
APP=myapp
TAG=$(date +%H%M%S)
IMAGE=$REGISTRY/$APP:$TAG

echo "Building $IMAGE"
docker build . --build-arg VERSION=$TAG -f Dockerfile -t $IMAGE

echo "Copying $IMAGE to kind cluster"
kind load docker-image $IMAGE

echo "Applying deployment to kubenetes"
cat "deployment.yaml" | sed "s/MYAPP_TAG/${TAG}/g" | kubectl apply -f -
