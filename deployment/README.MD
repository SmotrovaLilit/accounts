## K8s locally

### MAC 

#### Install minikube

https://minikube.sigs.k8s.io/docs/start/?arch=%2Fmacos%2Fx86-64%2Fstable%2Fhomebrew

```bash
brew install minikube
```

#### Turn on/off minikube

```bash
brew unlink minikube
brew link minikube
```

#### Start minikube

```bash
minikube start -n 3
```
s
```bash
minikube kubectl -- get po -A
kubectl get po -A
```

#### Install kubectl
https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/#install-with-homebrew-on-macos

```bash
brew install kubectl
```

#### Delete minikube

```bash
minikube delete
```

## Install terraform
https://developer.hashicorp.com/terraform/tutorials/docker-get-started/install-cli
