#!/usr/bin/env bash

UPSTREAM_REMOTE="upstream"
MASTER_BRANCH="master"
HUB_NAMESPACE="tekton-hub"
HUB_CI_NAMESPACE="tekton-hub-ci"
BINARIES="kubectl git"
CLUSTER=""

set -e

for b in ${BINARIES};do
    type -p ${b} >/dev/null || { echo "'${b}' need to be avail"; exit 1 ;}
done

kubectl version 2>/dev/null >/dev/null || {
    echo "you need to have access to a kubernetes cluster"
    exit 1
}

kubectl get pipelineresource 2>/dev/null >/dev/null || {
    echo "you need to have tekton install onto the cluster"
    exit 1
}

echo; echo 'Tekton Hub Deployment: '

[[ -z ${RELEASE_VERSION} ]] && {
   read -e -p "Enter a target release (i.e: v0.1.2): " RELEASE_VERSION
   [[ -z ${RELEASE_VERSION} ]] && { echo "no target release"; exit 1 ;}
}
[[ ${RELEASE_VERSION} =~ v[0-9]+\.[0-9]*\.[0-9]+ ]] || { echo "invalid version provided, need to match v\d+\.\d+\.\d+"; exit 1 ;}

read -e -p "Are you deploying on OpenShift (Y/n): " YESORNO
if [ "${YESORNO}" == "Y" ] || [ "${YESORNO}" == "y" ]; then
  CLUSTER='openshift'
elif [ "${YESORNO}" == "N" ] || [ "${YESORNO}" == "n" ]; then  
  CLUSTER='kubernetes'
else 
  echo 'invalid input'
  exit 1
fi
echo ${CLUSTER}

cd ${GOPATH}/src/github.com/tektoncd/hub

[[ -n $(git status --porcelain 2>&1) ]] && {
   echo "We have detected some changes in your repo"
   echo "Stash them before executing this script"
   exit 1
}

git checkout ${MASTER_BRANCH}
git reset --hard ${UPSTREAM_REMOTE}/${MASTER_BRANCH}

echo; echo 'Creating tag for new release:  ';
read -e -p "Enter tag message: " TAG_MESSAGE
git tag -a ${RELEASE_VERSION} -m ${TAG_MESSAGE}
git push ${UPSTREAM_REMOTE} --tags

kubectl create namespace ${HUB_NAMESPACE} 2>/dev/null || true

kubectl -n ${HUB_NAMESPACE} get secret db 2>/dev/null >/dev/null || {
    echo; echo "Database Configurations:"
        read -e -p "Enter DB Name: " DB_NAME
        read -e -p "Enter DB Username: " DB_USERNAME
        read -e -p "Enter DB Password: " DB_PASSWORD

        kubectl -n ${HUB_NAMESPACE} create secret generic db \
            --from-literal=POSTGRES_DB=${DB_NAME} \
            --from-literal=POSTGRES_USER=${DB_USERNAME} \
            --from-literal=POSTGRES_PASSWORD=${DB_PASSWORD} \
            --from-literal=POSTGRES_PORT="5432"
        
        kubectl -n ${HUB_NAMESPACE} label secret db app=db
    echo; 
}

kubectl -n ${HUB_NAMESPACE} get secret api 2>/dev/null >/dev/null || {
    echo "API Configurations:"
        read -e -p "Enter GitHub OAuth Client ID: " GH_CLIENT_ID
        read -e -p "Enter GitHub OAuth Client Secret: " GH_CLIENT_SECRET
        read -e -p "Enter JWT Signing key: " JWT_SIGNING_KEY
        read -e -p "Enter the Access JWT expire time: (eg. 1d) " ACCESS_JWT_EXPIRES_IN
        read -e -p "Enter the Refresh JWT expire time: (eg. 1d) " REFRESH_JWT_EXPIRES_IN

        kubectl -n ${HUB_NAMESPACE} create secret generic api \
            --from-literal=GH_CLIENT_ID=${GH_CLIENT_ID} \
            --from-literal=GH_CLIENT_SECRET=${GH_CLIENT_SECRET} \
            --from-literal=JWT_SIGNING_KEY=${JWT_SIGNING_KEY} \
            --from-literal=ACCESS_JWT_EXPIRES_IN=${ACCESS_JWT_EXPIRES_IN} \
            --from-literal=REFRESH_JWT_EXPIRES_IN=${REFRESH_JWT_EXPIRES_IN}
        
        kubectl -n ${HUB_NAMESPACE} label secret api app=api 

        kubectl -n ${HUB_NAMESPACE} create cm ui \
            --from-literal=GH_CLIENT_ID=${GH_CLIENT_ID} \
            --from-literal=API_URL="https://api.hub.tekton.dev" 
        
        kubectl -n ${HUB_NAMESPACE} label cm ui app=ui
    echo;         
}

kubectl -n ${HUB_NAMESPACE} get cm api 2>/dev/null >/dev/null || {
    echo "Hub Config File:"
        read -e -p "Enter Raw URL of the hub config file (Default: https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml): " HUB_CONFIG

        if [ -z "$HUB_CONFIG" ]
        then
            HUB_CONFIG=https://raw.githubusercontent.com/tektoncd/hub/master/config.yaml
        fi

        kubectl -n ${HUB_NAMESPACE} create cm api \
            --from-literal=CONFIG_FILE_URL=${HUB_CONFIG} 
        
        kubectl -n ${HUB_NAMESPACE} label cm api app=api
    echo; 
}

kubectl create namespace ${HUB_CI_NAMESPACE} 2>/dev/null || true

kubectl -n ${HUB_CI_NAMESPACE} delete secret registry-sec --ignore-not-found
kubectl -n ${HUB_CI_NAMESPACE} get secret registry-sec 2>/dev/null >/dev/null || {
    echo; echo "Enter Quay registry credentials to push the images: (quay.io/tekton-hub) "
        read -e -p "Enter Username: " USERNAME
        read -e -sp "Enter Password: " PASSWORD

        kubectl -n ${HUB_CI_NAMESPACE} create secret generic registry-sec \
            --type="kubernetes.io/basic-auth"  \
            --from-literal=username=${USERNAME} \
            --from-literal=password=${PASSWORD}

        kubectl -n ${HUB_CI_NAMESPACE} annotate secret registry-sec tekton.dev/docker-0=quay.io
}

echo; echo 'Creates service account and necessary role to create resources: '

kubectl -n ${HUB_CI_NAMESPACE} delete serviceaccount registry-login --ignore-not-found
cat <<EOF | kubectl -n ${HUB_CI_NAMESPACE} create -f-
apiVersion: v1
kind: ServiceAccount
metadata:
  name: registry-login
secrets:
  - name: registry-sec
EOF

kubectl -n ${HUB_NAMESPACE} delete role hub-pipeline --ignore-not-found
kubectl -n ${HUB_NAMESPACE} delete rolebinding hub-pipeline --ignore-not-found
kubectl -n ${HUB_NAMESPACE} create role hub-pipeline \
  --resource=deployment,services,pvc,job \
  --verb=create,get,list,delete
kubectl -n ${HUB_NAMESPACE} create rolebinding hub-pipeline \
  --serviceaccount=${HUB_CI_NAMESPACE}:registry-login \
  --role=hub-pipeline

if [ "${CLUSTER}" == "openshift" ]; then
  oc adm policy add-scc-to-user privileged system:serviceaccount:${HUB_CI_NAMESPACE}:registry-login
fi

echo; echo 'Install Tasks: '

kubectl -n ${HUB_CI_NAMESPACE} apply -f https://raw.githubusercontent.com/tektoncd/catalog/master/task/git-clone/0.2/git-clone.yaml
kubectl -n ${HUB_CI_NAMESPACE} apply -f https://raw.githubusercontent.com/tektoncd/catalog/master/task/buildah/0.2/buildah.yaml
kubectl -n ${HUB_CI_NAMESPACE} apply -f https://raw.githubusercontent.com/tektoncd/catalog/master/task/golangci-lint/0.1/golangci-lint.yaml
kubectl -n ${HUB_CI_NAMESPACE} apply -f https://raw.githubusercontent.com/tektoncd/catalog/master/task/kubernetes-actions/0.2/kubernetes-actions.yaml
kubectl -n ${HUB_CI_NAMESPACE} apply -f https://raw.githubusercontent.com/tektoncd/catalog/master/task/npm/0.1/npm.yaml
kubectl -n ${HUB_CI_NAMESPACE} apply -f ./tekton/api/golang-db-test.yaml

echo; echo 'Install Pipelines: '

kubectl -n ${HUB_CI_NAMESPACE} apply -f ./tekton/api/pipeline.yaml
kubectl -n ${HUB_CI_NAMESPACE} apply -f ./tekton/ui/pipeline.yaml

echo; echo 'Start Pipelines: '

cat <<EOF | kubectl -n ${HUB_CI_NAMESPACE} create -f-
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  generateName: hub-api-
spec:
  serviceAccountName: registry-login
  pipelineRef:
    name: api-deploy
  params:
    - name: HUB_REPO
      value: https://github.com/tektoncd/hub
    - name: REVISION
      value: ${MASTER_BRANCH}
    - name: API_IMAGE
      value: quay.io/tekton-hub/hub-api
    - name: DB_MIGRATION_IMAGE
      value: quay.io/tekton-hub/hub-db
    - name: TAG
      value: ${RELEASE_VERSION}
    - name: HUB_NAMESPACE
      value: ${HUB_NAMESPACE}
    - name: K8S_VARIANT #it will accept either openshift or kubernetes
      value: ${CLUSTER}
  workspaces:
    - name: shared-workspace
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 500Mi
EOF

cat <<EOF | kubectl -n ${HUB_CI_NAMESPACE} create -f-
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  generateName: hub-ui-
spec:
  serviceAccountName: registry-login
  pipelineRef:
    name: ui-pipeline
  params:
    - name: HUB_REPO
      value: https://github.com/tektoncd/hub
    - name: REVISION
      value: ${MASTER_BRANCH}
    - name: IMAGE
      value: quay.io/tekton-hub/hub-ui
    - name: TAG
      value: ${RELEASE_VERSION}
    - name: HUB_NAMESPACE
      value: ${HUB_NAMESPACE}
    - name: K8S_VARIANT #it will accept either openshift or kubernetes
      value: ${CLUSTER}
  workspaces:
    - name: shared-workspace
      volumeClaimTemplate:
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
EOF