# Namespace Cleaner

A simple, yet beautiful tool designed to delete Kubernetes namespaces based on specific strings or text patterns.

## How does it work?

This tool runs as a Cronjob within a Kubernetes Cluster (in the default namespace) and is scheduled to execute every hour. It utilizes the environment variable "NAMESPACE_SELECTOR" to identify and delete desired namespaces along with all their contents.

## How to use

1. Build and publish the tool to your preferred container registry:

docker build -t namespace-cleaner:1.0 .
docker tag namespace-cleaner:1.0 rmnobarra/namespace-cleaner:1.0
docker push rmnobarra/namespace-cleaner:1.0

2. Adjust the job execution schedule parameter inside the "namespace-selector.yaml" file (line number 40).

3. Update the image URL with the location of your own built image inside the "namespace-selector.yaml" file (line number 48).

4. Create the Cronjob resource within your Kubernetes cluster:

```bash
kubectl apply -f namespace-selector.yaml
```

4. Enjoy the benefits of automated namespace cleaning!