# Usage:

go run main.go [command] [flags]

## Available Commands:
  **repo**:        Generates a terraform script for setting up a github repo for the customer. \
  **cluster**:     Generates a terraform script for deploying a new k8s cluster. (subcommands: gke, eks, aks) \
  **manifest**:    Generates k8s templates to deploy the managed fusion services. \
  **help**:        Help about any command

### Example:

To generate a terraform script for setting up customer's github repo:
```
go run main.go repo --customerName=<customer_name>
```

To generate a terraform script for deploying a new gke cluster:
```
go run main.go cluster gke --customerName=<customer_name> --projectId=<gcp_project_id> --region=<gke_region>  --clusterType=<non-prod or prod> --clusterSize=large/small

### Create Large Prod cluster
go run main.go cluster gke --customerName=<customer_name> --projectId=<gcp_project_id> --region=<gke_region>  --clusterType=prod --clusterSize=large

### Create Small Prod Cluster
go run main.go cluster gke --customerName=<customer_name> --projectId=<gcp_project_id> --region=<gke_region>  --clusterType=prod --clusterSize=small

```

To generate k8s templates for deploying the managed fusion services:
```
go run main.go manifest --customerName=<customer_name> --clusterType=<non-prod or prod> --env=<dev,stg, or prd> --namespace=<namespace> --cloud=<gke, eks, or aks> --version=<fusion_version> --clusterSize=small/large
```