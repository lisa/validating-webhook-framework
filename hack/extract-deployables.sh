#!/bin/bash

set -e

## Extract deployables from the template rendered by the `make syncset` target

__PWD=$(dirname $0)

GIT_REV=$(git -C ${__PWD}/.. rev-parse --short=7 HEAD)
REPO_NAME=${1:-quay.io/${USER}/validating-webhook-framework}
IMAGE_TAG=${2:-$GIT_REV}

IMAGE_DIGEST=$(docker inspect ${REPO_NAME}:${IMAGE_TAG} | jq -r '.[0].Id')

make -C ${__PWD}/../ IMG_ORG="${USER}" IMG="${REPO_NAME}" all syncset 1>/dev/null


oc \
  process                                        \
  --local                                        \
  -f ${__PWD}/../build/selectorsyncset.yaml \
  IMAGE_TAG=$IMAGE_TAG REPO_NAME=${REPO_NAME} \
  IMAGE_DIGEST=${IMAGE_DIGEST} \
  REGISTRY_IMG=${REPO_NAME} \
  -o yaml >${__PWD}/selectorsyncset.yaml

ruby <<EOF
require 'yaml'
y = YAML::load_file('${__PWD}/selectorsyncset.yaml')
res = y['items'][0]['spec']['resources']
c = 0
khist = Hash.new(0)
res.each do |r|
  kind = r['kind']
  khist[kind] += 1
  c += 1
  File.open("${__PWD}/#{c.to_s.rjust(3,"0")}_#{kind}_#{khist[kind]}.yaml","w") do |f|
    f.write r.to_yaml
  end
end
EOF
