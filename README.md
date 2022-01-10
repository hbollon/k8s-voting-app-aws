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
  - [Project structure](#project-structure)
- [Getting Started](#getting-started)
  - [Run with Docker Compose](#run-with-docker-compose)
  - [Run with Kubernetes](#run-with-kubernetes)
    - [Minikube](#minikube)
      - [Using k8s templates](#using-k8s-templates)
      - [Using Helm Chart](#using-helm-chart)
    - [AWS](#aws)
- [Roadmap](#roadmap)
- [Contributing](#-contributing)
- [Author](#author)
- [Show your support](#show-your-support)
- [License](#-license)

## Presentation

This repository provide a complete and modern ready to deploy example of a dockerized and distributed app. Deployable using [Docker-Compose](https://docs.docker.com/compose/), [Kubernetes](https://kubernetes.io/) templates or even [Helm](https://helm.sh/) Chart.

### Architecture

![App's architecture scheme](docs/app-architecture.jpg)

### Project structure

```bash
k8s-voting-app-aws/
├─ .github/ # Github workflows
├─ docs/
│  ├─ app-architecture.jpg # App's architcture scheme
│  ├─ README-FR.md # French translation of the readme
├─ helm/ # Helm Chart definitions
├─ k8s-specifications/ # K8s Templates files
├─ voting-app/ # Result, Vote and Worker source code 
├─ *.tf # terraform specs files
├─ *.tfvars # terraform values files
├─ *.yml # docker-compose files
```

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

1. Deploy all k8s ressources: `kubectl apply -f k8s-specifications --namespace=voting-app-stack`
2. Get your cluster IP using: `minikube ip`
3. Enable ingress access:
   - **On Linux:** Edit your hosts file located at `/etc/hosts` by adding `<minikube ip> result.votingapp.com vote.votingapp.com` to the end of it, of course replace `<minikube ip>` by the real cluster ip.
   - **On Windows:** Edit your hosts file located at `c:\Windows\System32\Drivers\etc\hosts` by adding `127.0.0.1 result.votingapp.com vote.votingapp.com` to the end of it.
   After that, start a Minikube tunnel: `minikube tunnel`

The result app should be now accessible through `result.votingapp.com` and the vote one to `vote.votingapp.com`
To stop and destroy all k8s deployed ressources run: `kubectl delete -f k8s-specifications --namespace=voting-app-stack` and stop minikube using `minikube stop`

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

#### AWS

To deploy the app to AWS you must create an infrastructure based on EKS (Elastic Kubernetes Service) first. You have all the needed Terraform definitions to do that easily with an IAC interface (Infrastructure As Code).
You must have an AWS account to follow this guide and be careful, although AWS has free tier for new accounts, this infrastructure can generate some costs although very limited. Especially in case of bad configuration or usage where the costs can be multiplied.

**I will not be responsible for the invoice generated in any way.**

1. Install [Terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli), [AWS-CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) and [KubeCTL](https://kubernetes.io/docs/tasks/tools/)

2. Clone this repo

3. Create a new IAM User on your AWS Account:
   - Go to the IAM section and create a new user named "TerraformUser"
   - Add this user to a new group named "TerraformFullAccessGroup" with **AdministratorAccess** and **AmazonEKSClusterPolicy** rights
   - Once done, keep the _Secret Access Key_ and _Access Key ID_, this will be the only time AWS gives it to you

4. Go to the VPC panel of the [AWS Console](https://console.aws.amazon.com/console/home) and get two differents subnet ids from the default VPC. Add these two ids in the `values.tfvars` file at the root of this project (replace `<subnet_id_1>` and `<subnet_id_2>`).

5. For the following steps you will need to use a credentials management method to use them with Terraform and AWS-CLI. The easier way is to set AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environement variables. But you can also use tools like [Summon](https://github.com/cyberark/summon) or AWS config files.

6. Open a console at the root of this project directory and execute:
  - `terraform init`
  - `terraform plan -var-file=values.tfvars`: check that the output is generated without any errors. 
  - `terraform apply -var-file=values.tfvars` (this operation can take a while don't worry)

7. If previous commands runs well you should now have a working EKS cluster, in order to link your kubectl installation to it you must run: `aws eks update-kubeconfig --region eu-west-3 --name eks_cluster_voting_app` **change the region flag if you have deployed the EKS on another one**.
Once done, run `kubectl get pods -A`, if it working you've done with your fresh EKS cluster.

8. Finally, deploy all the k8s ressources:
   - With k8s templates: `kubectl apply -f k8s-specifications --namespace=voting-app-stack`
   - With Helm Chart: `helm dependency update ./helm/voting-app` and after: `helm template voting-app ./helm/voting-app --namespace=voting-app-stack | kubectl apply -f -`

You have the possibility to get the aws endpoint linked to your EKS cluster by running: `terraform output -json`. However, the ingress ressource is not compatible with it atm.

You can destroy everything just by running this command: `terraform destroy -var-file=values.tfvars`

**Never delete the generated .tfstate file when the infrastucture is deployed!** Without it you will be unable to delete all the AWS ressources with Terraform and you will be forced to do it manually with the Web AWS Console or the AWS-CLI.

## Roadmap

Many additional features are coming, including:

- Ingress compatibility with AWS and domain customization
- Monitoring/Alerting/Dashboarding using kube-prometheus-stack
- Style webapps with CSS
- And many more to come !

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
