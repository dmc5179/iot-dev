#!/bin/bash

BASE_DOMAIN="example.io"
CLUSTER_NAME="caas"
CEPH_USER='danclark'

go run main.go knative setup

go run main.go kafka setup

go run main.go kafka bridge

curl -v GET my-bridge-route-kafka.apps.${CLUSTER_NAME}.${BASE_DOMAIN}/healthy

go run main.go ceph setup

go run main.go ceph user ${CEPH_USER}

go run main.go ceph secrets ${CEPH_USER}

export CEPH_ENDPOINT='ceph-route-rook-ceph.apps.${CLUSTER_NAME}.${BASE_DOMAIN}'
export CEPH_SECRET=''
export CEPH_KEY=''

go run main.go tensorflowServing setup -n kafka

go run main.go knative service video-analytics -n kafka --cephEndpoint "${CEPH_ENDPOINT}" --cephAccessKey "${CEPH_KEY}" --cephSecretKey "${CEPH_SECRET}"

go run main.go knative source kafka video-analytics -n kafka

#
#
# Start the IoT Video Simulator here
# cd ~/workspace/2020Summit-IoT-Streaming-Demo/iotDeviceSimulator-kafka
# export STREAMURL=https://www.youtube.com/watch?v=wjsHwDFlPKc
# export ENDPOINT=my-bridge-route-kafka.apps.${CLUSTER_NAME}.${BASE_DOMAIN}/topics/my-topic
# go run ./cmd


go run main.go knative service video-serving -n kafka --cephEndpoint "${CEPH_ENDPOINT}" --cephAccessKey "${CEPH_KEY}" --cephSecretKey "${CEPH_SECRET}"

oc get ksvc

