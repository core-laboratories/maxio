# Deploy MaxIO on Kubernetes [![Slack](https://slack.min.io/slack?type=svg)](https://slack.min.io)  [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

MaxIO is a high performance distributed object storage server, designed for large-scale private cloud infrastructure. MaxIO is designed in a cloud-native manner to scale sustainably in multi-tenant environments. Orchestration platforms like Kubernetes provide perfect cloud-native environment to deploy and scale MaxIO.

## MaxIO Deployment on Kubernetes

There are multiple options to deploy MaxIO on Kubernetes:

- MaxIO-Operator: Operator offers seamless way to create and update highly available distributed MaxIO clusters. Refer [MaxIO Operator documentation](https://github.com/minio/minio-operator/blob/master/README.md) for more details.

- Helm Chart: MaxIO Helm Chart offers customizable and easy MaxIO deployment with a single command. Refer [MaxIO Helm Chart documentation](https://github.com/core-laboratories/maxio/tree/master/helm/maxio) for more details.

## Monitoring MaxIO in Kubernetes

MaxIO server exposes un-authenticated liveness endpoints so Kubernetes can natively identify unhealthy MaxIO containers. MaxIO also exposes Prometheus compatible data on a different endpoint to enable Prometheus users to natively monitor their MaxIO deployments.

## Explore Further

- [MaxIO Erasure Code QuickStart Guide](https://docs.min.io/community/minio-object-store/operations/concepts/erasure-coding.html)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Helm package manager for kubernetes](https://helm.sh/)
