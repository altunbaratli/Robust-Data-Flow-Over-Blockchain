#!/bin/bash

cd volume &&
kubectl create -f createPVandPVC.yaml &&
sleep 10 &&
cd ../jobs &&
kubectl apply -f copyArtifactsJob.yaml &&
sleep 10 &&
pod=$(kubectl get pods --selector=job-name=copyartifacts --output=jsonpath={.items..metadata.name}) &&
kubectl cp ../artifacts $pod:/shared/
kubectl apply -f generateCryptoConfig.yaml &&
kubectl apply -f generateGenesisBlock.yaml &&
kubectl apply -f generateChanneltx.yaml &&
kubectl apply -f generateAnchorPeerMSPs.yaml &&
cd ../network-deployment &&
sh deployAll.sh &&
sleep 10 &&
cd ../jobs &&
kubectl apply -f create_channel.yaml &&
sleep 10 &&
kubectl apply -f join_channel.yaml &&
sleep 10 &&
kubectl apply -f chaincode_install.yaml &&
sleep 10 &&
kubectl apply -f chaincode_instantaite.yaml &&
sleep 10 &&
kubectl apply -f updateAnchorPeers.yaml &&
cd .. &&
kubectl create deployment rest-api --image=altunbarat/rest-api &&
kubectl expose deployment rest-api --port=3000 --target-port=3000 &&
cd node-red &&
kubectl create deployment node-red --image=yigitpolat/hyperledger-iot-nodered &&
kubectl apply -f node-red-svc-nodePort.yaml &&
cd ../network-deployment &&
kubectl apply -f expose-restapi.yaml &&
kubectl get nodes -o wide