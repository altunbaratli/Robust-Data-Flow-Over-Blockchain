#!/bin/bash

kubectl delete --all deployments &&
kubectl delete --all jobs &&
kubectl delete --all services &&
kubectl delete --all pods &&
kubectl delete --all pvc &&
kubectl delete --all pv
