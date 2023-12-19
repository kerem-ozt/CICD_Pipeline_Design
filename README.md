# CI/CD Pipeline Task

CI/CD pipeline design using Google Kubernetes Engine provider created using Gitlab CI, ArgoCD and helm chart.

## Table of Contents

- [Overview](#overview)
- [File Hierarchy](#file-hierarchy)
- [Helm Chart](#helm-chart)
- [Gitlab CI](#gitlab-ci)
- [ArgoCD](#argocd)
- [About Google Kubernetes Engine](#about-google-kubernetes-engine)

## Overview

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/1.jpeg)

The designed CI/CD pipeline follows a general flow illustrated in the diagram. GitLab manages the CI process, while ArgoCD manages the CD process. The sequence begins with a developer making changes in the repository and committing them. GitLab, detecting the new commit, triggers the pipeline. There are predefined steps in the pipeline, these steps include operations such as static code analysis, unit tests, image creation and pushing the image to the register. Subsequent sections of this article will delve into these steps and the GitLab YML file that orchestrates the process.

Furthermore, the Kubernetes cluster is hosted on Google Kubernetes Engine (GKE), a cloud provider. ArgoCD plays a pivotal role by checking for changes in the manifest files used in Kubernetes pods. If there are changes, ArgoCD is triggered based on the configurations specified in the Argo application. ArgoCD manages resources in the GKE clusters, detecting alterations in manifest files and applying them to the clusters. This includes tasks like creating new pods, updating services, or synchronizing other Kubernetes resources. Through this streamlined process, a basic yet effective end-to-end CI/CD pipeline is established.

Even a rudimentary pipeline like this can significantly enhance operational efficiency. Regardless of how meticulous one is, the manual building and pushing of images entail a risk of introducing faulty versions because of humans. This automated process not only minimizes human errors but also ensures the consistency of software versions.

## File Hierarchy

The file tree of the main branch of the application is shown below.

```
├── Dockerfile
├── LICENSE
├── README.md
├── controllers
│   ├── auth.go
│   ├── leaderBoard.go
│   ├── ping.go
│   ├── tournament.go
│   └── user.go
├── docker-compose.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── main.go
├── middlewares
│   ├── auth.go
│   ├── cors.go
│   ├── logger.go
│   ├── recovery.go
│   └── validators
│       ├── auth.validator.go
│       ├── base.validator.go
│       ├── tournament.validator.go
│       └── user.validator.go
├── models
│   ├── config.go
│   ├── db
│   │   ├── country.go
│   │   ├── leaderBoard.go
│   │   ├── token.go
│   │   ├── tournament.go
│   │   └── user.go
│   ├── request.go
│   └── response.go
├── routes
│   ├── auth.go
│   ├── leaderBoard.go
│   ├── ping.go
│   ├── router.go
│   ├── tournament.go
│   └── user.go
├── services
│   ├── config.service.go
│   ├── leaderBoard.service.go
│   ├── storage.service.go
│   ├── token.service.go
│   ├── tournament.service.go
│   └── user.service.go
├── tests
│   └── ping_test.go
└── tournament-cicd
    ├── Chart.yaml
    ├── templates
    │   ├── configmap.yaml
    │   ├── deployment.yaml
    │   └── service.yaml
    ├── values-dev.yaml
    └── values.yaml
```
The main directory contains the source codes of the application, Dockerfile, docker-compose.yaml and tournament-cicd, which contains the Helm chart. Apart from the cd pipeline of the application, the repository where its structure is explained can also be accessed at https://github.com/kerem-ozt/Tournament-Management-API.

## Helm Chart

The directory containing Helm chart min and related files is shown in the tree below.

```
└── tournament-cicd
    ├── Chart.yaml
    ├── templates
    │   ├── configmap.yaml
    │   ├── deployment.yaml
    │   └── service.yaml
    ├── values-dev.yaml
    └── values.yaml
```
This structure was created based on the template created using the helm create tool on Helm's official website.

Basic information about the application is available in the Chart.yaml file, including the name of the Helm chart, a brief description, type (indicating whether it's an application or library chart), version number, and the application's deployed version number.

The values.yaml file holds the application's default configurable values. Many properties are defined as defaults for Helm chart creation, such as replicaCount, image details, Kubernetes service type and port, ingress features, resource limitations, autoscaling settings, nodeSelector, tolerances, affinity configurations specific to Kubernetes pods, app name, port, namespace, configmap details, and Docker image name and tag.

The values-dev.yaml file contains development (dev) environment-specific values that will override those in values.yaml, specifically pertaining to namespace and configmap.

The .helmignore file contains files and directories that are excluded from the Helm build.

The configmap.yaml represents the ConfigMap resource populated with dynamic values using the .Values syntax within Helm templates.

The deployment.yaml file represents the Kubernetes Deployment resource, including parameters for the cluster that ArgoCD will use. It specifies the number of replicas, selector labels, pod template, container details (image, ports, envFrom), and resource requests for the container.

The service.yaml file represents the Kubernetes Service resource, indicating the port listened to by the service and the tags used to select pods managed by the service.

This Helm chart is designed for deploying and managing your application on Kubernetes, offering benefits such as configuration management, version control, and reusability. The chart's versatility allows you to change values for different environments, making it suitable for use in development, test, and production environments, facilitating the distribution of your application.

## Gitlab CI

The .gitlab-ci.yml file located in the main directory of the project is used to configure the pipeline of the project. It is located in the root folder of the repo and contains the definition of your Pipelines, Jobs and Environments. In this pipeline I have defined 5 stages. Lint, test, update_image_tag, build and release and 5 jobs respectively; lint, test, update_image_tag, build and build_image. 

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/2.png)

Linting Phase:
In the linting phase, static code analysis is performed to identify formal errors, simple bugs, and suspicious structures.

Testing Phase:
During the testing phase, the application undergoes a series of tests to ensure its functionality and integrity.

Update Tag Name Stage:
After a commit is made at the update tag name stage, the header variable in the values.yaml file is modified to trigger ArgoCD. ArgoCD is then enabled to deploy the new changes to the Kubernetes cluster in its current state. In the current version of the file, this section is in a comment line and can be optionally opened. This flexibility is valuable in development environments where updating pods in the cluster after each commit may not always be advisable.

Build Phase:
During the build phase, a folder is created using a name that can be modified using GitLab variables. Subsequently, the Go project is compiled using the go build command.

Build Image Phase:
In the build image phase, a Docker image is generated with the current source code and pushed to the Docker registry with the specified tag. The project utilizes the mkeremozt/game-api-go directory available on DockerHub as the registry.

## ArgoCD

ArgoCD, along with Helm, was installed on the cluster. Notably, during the installation process, the default namespace was not utilized; instead, a new namespace named 'argocd' was created.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/3.png)

As evidenced by the Google Cloud GUI, when ArgoCD is functioning correctly, it initiates the execution of ArgoCD services within the 'argocd' namespace on the nodes. However, due to the implementation on Google Kubernetes Engine, direct access to the ArgoCD GUI was not feasible.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/4.png)

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/5.png)

First, I had to create firewall rules for the relevant port via gcloud sdk shell, and then export the port from my local computer.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/6.png)

After these steps, I followed the document on ArgoCD's official website, got my user password and logged in as admin. Then, I introduced my github repo, which I was monitoring with gitlab, to ArgoCD and created an ArgoCD application as shown in the figure.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/7.png)

After all these steps, I obtained an ArgoCD app that tracks the changes in the manifest files in my existing repo in gitlab and shows the status, monitoring and health status of my automatically synchronizing pods.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/8.png)

## About Google Kubernetes Engine

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/9.png)

I wanted to use a cloud provider for a more elegant solution in the project and chose Google's Kubernetes Engine service. Some monitoring interfaces of the cluster, which can be accessed through the Google Kubernetes Engine interface, are shown in the photographs.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/10.png)

The image above shows the general status of the cluster, the container's logs and exposing services.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/11.png)

Additionally, as seen in the last interface image, the pods and even the ports held by the pods can also be viewed via the interface. Although accessing these ports may be a bit difficult according to local clusters, it is possible to determine the necessary security rules for this port and open the port to the outside.

Below, the running pods, the status of the pod we deployed and the service logs are displayed via the online Cloud Shell Editor.

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/12.png)

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/13.png)

![](https://gitlab.com/mkeremozt/Tournament-Management-API/-/raw/main/images/14.png)
