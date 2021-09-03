# EzBackup

An easy, simple way for backing up Persistent Volume Claims inside your Kubernetes Cluster without installing anything into your cluster.

## How does it work

This repository contains a CLI written in Go, a Docker image (for both x86 and ARMv7) shipping the CLI and [restic](https://github.com/restic/restic) for the actual backups and multiple example YAML files to just `kubectl apply` to your cluster.

This will allow you to generate backup Jobs or CronJobs using standard Kubernetes resources to scale down any Pods accessing your PVC, creating an incremental backup to a remote location (anything supported by restic) and scaling everything back as expected.