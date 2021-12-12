# EzBackup

Easy, simple backup solution for [Persistent Volume Claims](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) in your Kubernetes Cluster.

* No components to install into your cluster
* Fully customizable
* Extensive, readable logs for every run
* Fast, secure, efficient backups via [restic](https://github.com/restic/restic)
* All in one small Docker image: [`ghcr.io/lukasknuth/ezbackup`](https://github.com/LukasKnuth/EzBackup/pkgs/container/ezbackup)
* Helpers to statically generate Job/CronJob

## Why it's needed

Taking backups of persistent application state while the application is still running is risky because data might not be fully flushed to disk yet. At best the backup is incomplete, at worst it's corrupted.

This is especially true for applications like Databases which are very picky about when they do disk IO to optimize throughput.

A simple solution to this problem is to stop the application before taking the backup. This is where EzBackup can help.

## How it works

1. The EzBackup CLI finds the Pods mounting a given Persistent Volume Claim with write access
2. For each Pod it walks the ownership chain to find the owning controller (for example a Deployment)
3. All found controllers are scaled down and EzBackup waits for Pod termination
4. restic is used to do the actual backup
5. All controllers are scaled back up to their original replica count

# Usage 

## Setup

Since EzBackup uses the Kubernetes API from inside the cluster, the service account used by it's Pod must be allowed to perform certain actions.

A ClusterRole with all required permissions is available in `install/001-cluster-role.yaml`. The ClusterRole resource is not namespaced, so it only needs to be created once in the cluster.

Now bind the ClusterRole to a ServiceAccount using a RoleBinding. It's advised to create a new ServiceAccount specific to EzBackup rather than changing the namespaces default ServiceAccount. `install/002-account-and-role.yaml` creates both. **Note** that since these resources are namespaced, you need to customize the file to create the resources in the namespace you intend to use.

Lastly, specify the ServiceAccount to use in the Pod template under `spec.serviceAccountName`

## Configuring restic

The actual backup is done by [restic](https://github.com/restic/restic). It supports repositories with [different storage backends](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html).

> **NOTE:** The Docker image does NOT contain rclone, so repository types supported through it are not available!

The executable is started by EzBackup, no parameters are passed on. However, any environment variables set on the container are passed onto restic.

Usually, you'll want to set at least `RESTIC_REPOSITORY` and `RESTIC_PASSWORD`. Depending on the repository type, additional variables are required. See the [restic documentation](https://restic.readthedocs.io/en/stable/040_backup.html#environment-variables) for all available options.

## Manual backup

Simply start the EzBackup container in a [Job](https://kubernetes.io/docs/concepts/workloads/controllers/job/):

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: backup
  namespace: &namespace testing # 1
  annotations:
    ezbackup.codeisland.org/pvc-name: &pvc-name test-claim # 2
spec:
  backoffLimit: 4
  template:
    spec:
      serviceAccountName: ezbackup-automation # 3
      restartPolicy: Never
      containers:
      - name: run-backup
        image: ghcr.io/lukasknuth/ezbackup:latest
        imagePullPolicy: Always
        args: ["backup", *pvc-name, "-n", *namespace]
        env:
        - name: RESTIC_REPOSITORY
          value: "your_repository_here" # 4
        - name: RESTIC_PASSWORD
          value: "YourVerySavePasswordGoesHere"
        volumeMounts:
        - name: target
          mountPath: /mnt/target # 5
          readOnly: true
      volumes:
      - name: target
        persistentVolumeClaim:
          claimName: *pvc-name
```

1. The namespace of the Job must match the namespace of the PVC to backup.
2. Since we the PVC name multiple times, we use a [YAML anchor](https://support.atlassian.com/bitbucket-cloud/docs/yaml-anchors/) here to refer to it.
3. References the Service Account setup in the ["Setup" section](#setup).
4. Configure the restic repository as described in the ["Configuring restic" section](#configuring-restic). **Note**: Since this probably contains sensitive credentials, in production you should use a [Secret resource](https://kubernetes.io/docs/concepts/configuration/secret/) instead.
5. By default, EzBackup expects the PVC to be mounted to `/mnt/target`. Customize with `BACKUP_TARGET_DIR`

## Scheduled backup

Simply adapt the Job from the previous section into a [CronJob](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/) and use the cron notation to specify the schedule.

## Restore

todo: add once "restore" command is implemented...

Meanwhile, since the backups are created by restic, you can simply restore them with any container with a restic binary and the repository information.