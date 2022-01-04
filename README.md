<h1 align="center">Kubernetes distributed voting app</h1>
<p align="center">
  <a href="https://github.com/hbollon/k8s-voting-app-aws/blob/master/LICENSE" target="_blank">
    <img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" />
  </a>

<p align="center"> ☸️ Example of a distributed voting app running on Kubernetes. Written in Golang with Terraform definitions to deploy to AWS</p>

---

## Table of Contents

- [Presentation](#presentation)
  - [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Run with Docker Compose](#run-with-docker-compose)
  - [Run with Kubernetes](#run-with-kubernetes)
    - [Minikube](#minikube)
      - [Using k8s templates](#using-k8s-templates)
      - [Using Helm Chart](#using-helm-chart)
- [Contributing](#-contributing)
- [Author](#author)
- [Show your support](#show-your-support)
- [License](#-license)

## Presentation

This repository provide a complete and modern ready to deploy example of a dockerized and distributed app. Deployable using [Docker-Compose](https://docs.docker.com/compose/), [Kubernetes](https://kubernetes.io/) templates or even [Helm](https://helm.sh/) Chart.

### Architecture

![App's architecture scheme](docs/app-architecture.jpg)

## Getting Started

### Run with Docker Compose

1. Download and install [Docker](https://docs.docker.com/get-docker/) and [Docker-Compose](https://docs.docker.com/compose/install/)
2. Clone this repository: `git clone git@github.com:hbollon/k8s-voting-app-aws.git` (you can alternatively use http)
3. Open a terminal inside the cloned repository folder and build Docker images: `docker-compose build`
4. Start all services: `docker-compose up -d`

The result app should be now accessible through `localhost:9091` and the vote one to `localhost:9090`
To stop all deployed ressources run: `docker-compose down`

### Run with Kubernetes

#### Minikube

Before deploying the app, you must install Minikube and start a cluster:

1. Install [Minikube](https://minikube.sigs.k8s.io/docs/start/)
2. Start a Minikube cluster: `minikube start`
3. Check that Minikube is fully up (`minikube status`) and kubectl successfully linked (`kubectl get pods -A`)
4. Enable Nginx Ingress Controller addon: `minikube addons enable ingress`

##### Using k8s templates

1. Deploy all k8s ressources: `kubectl apply -f k8s-specifications`
2. Get your cluster IP using: `minikube ip`
3. Enable ingress access:
   - **On Linux:** Edit your hosts file located at `/etc/hosts` by adding `<minikube ip> result.votingapp.com vote.votingapp.com` to the end of it, of course replace `<minikube ip>` by the real cluster ip.
   - **On Windows:** Edit your hosts file located at `c:\Windows\System32\Drivers\etc\hosts` by adding `127.0.0.1 result.votingapp.com vote.votingapp.com` to the end of it.
   After that, start a Minikube tunnel: `minikube tunnel`

The result app should be now accessible through `result.votingapp.com` and the vote one to `vote.votingapp.com`
To stop and destroy all k8s deployed ressources run: `kubectl delete -f k8s-specifications` and stop minikube using `minikube stop`

##### Using Helm Chart

1. Update Helm repositories and download dependencies: `helm dependency update ./helm/voting-app`
2. Deploy the Helm Chart: `helm template voting-app ./helm/voting-app --namespace=voting-app-stack | kubectl apply -f -`
3. Get your cluster IP using: `minikube ip`
4. Enable ingress access:
   - **On Linux:** Edit your hosts file located at `/etc/hosts` by adding `<minikube ip> result.votingapp.com vote.votingapp.com` to the end of it, of course replace `<minikube ip>` by the real cluster ip.
   - **On Windows:** Edit your hosts file located at `c:\Windows\System32\Drivers\etc\hosts` by adding `127.0.0.1 result.votingapp.com vote.votingapp.com` to the end of it.
   After that, start a Minikube tunnel: `minikube tunnel`

The result app should be now accessible through `result.votingapp.com` and the vote one to `vote.votingapp.com`
To stop and destroy all k8s deployed ressources run: `helm template voting-app ./helm/voting-app --namespace=voting-app-stack | kubectl delete -f -` and stop minikube using `minikube stop`

## 🤝 Contributing

Contributions are greatly appreciated!

1. Fork the project
2. Create your feature branch (```git checkout -b feature/AmazingFeature```)
3. Commit your changes (```git commit -m 'Add some amazing stuff'```)
4. Push to the branch (```git push origin feature/AmazingFeature```)
5. Create a new Pull Request

Issues and feature requests are welcome!
Feel free to check [issues page](https://github.com/hbollon/k8s-voting-app-aws/issues).

## Author

👤 **Hugo Bollon**

* Github: [@hbollon](https://github.com/hbollon)
* LinkedIn: [@Hugo Bollon](https://www.linkedin.com/in/hugobollon/)
* Portfolio: [hugobollon.me](https://www.hugobollon.me)

## Show your support

Give a ⭐️ if this project helped you!
You can also consider so sponsor me [here](https://github.com/sponsors/hbollon) ❤️

## 📝 License

This project is under [MIT](https://github.com/hbollon/k8s-voting-app-aws/blob/master/LICENSE.md) license.
