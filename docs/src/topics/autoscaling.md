# Auto-scaling

This guide covers auto-scaling for CAPL clusters. The recommended tool for auto-scaling on Cluster API is [Cluster
Autoscaler](https://www.github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler#cluster-autoscaler).

## Flavor

The auto-scaling feature is provided by an add-on as part of the [Cluster Autoscaler
flavor](./flavors/cluster-autoscaler.md).

## Configuration

By default, the Cluster Autoscaler add-on [runs in the management cluster, managing an external workload
cluster](https://cluster-api.sigs.k8s.io/tasks/automated-machine-management/autoscaling#autoscaler-running-in-management-cluster-using-service-account-credentials-with-separate-workload-cluster).

```
+------------+             +----------+
|    mgmt    |             | workload |
| ---------- | kubeconfig  |          |
| autoscaler +------------>|          |
+------------+             +----------+
```

A separate Cluster Autoscaler is deployed for each workload cluster, configured to only monitor node groups for the
specific namespace and cluster name combination.

## Role-based Access Control (RBAC)

### Management Cluster

Due to constraints with the Kubernetes RBAC system (i.e. [roles cannot be subdivided beyond
namespace-granularity](https://www.github.com/kubernetes/kubernetes/issues/56582)), the Cluster Autoscaler add-on is
deployed on the management cluster to prevent leaking Cluster API data between workload clusters.

### Workload Cluster

Currently, the Cluster Autoscaler reuses the `${CLUSTER_NAME}-kubeconfig` Secret generated by the bootstrap provider to
interact with the workload cluster. The kubeconfig contents must be stored in a key named `value`. Due to this, all
Cluster Autoscaler actions in the workload cluster are performed as the `cluster-admin` role.

## Scale Down

> Cluster Autoscaler decreases the size of the cluster when some nodes are consistently unneeded for a significant
> amount of time. A node is unneeded when it has low utilization and all of its important pods can be moved elsewhere.

By default, Cluster Autoscaler scales down a node after it is marked as unneeded for 10 minutes. This can be adjusted
with the [`--scale-down-unneeded-time`
setting](https://www.github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#how-can-i-modify-cluster-autoscaler-reaction-time).

### Kubernetes Cloud Controller Manager for Linode (CCM)

The [Kubernetes Cloud Controller Manager for
Linode](https://www.github.com/linode/linode-cloud-controller-manager?tab=readme-ov-file#the-purpose-of-the-ccm) is
deployed on workload clusters and reconciles Kubernetes Node objects with their backing Linode infrastructure. When
scaling down a node group, the Cluster Autoscaler also deletes the Kubernetes Node object on the workload cluster. This
step preempts the Node-deletion in Kubernetes triggered by the CCM.

## Additional Resources

- [Autoscaling - The Cluster API Book](https://cluster-api.sigs.k8s.io/tasks/automated-machine-management/autoscaling)
- [Cluster Autoscaler
  FAQ](https://www.github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler#cluster-autoscaler)