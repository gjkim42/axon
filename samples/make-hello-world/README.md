# Make Hello World

This sample matches the [Quick Start kubectl/YAML](../../README.md#3-run-your-first-task) example from the README.

It runs a single Task that creates a hello world program in Python, using a Workspace to clone a git repo.

## Prerequisites

- Axon installed on your cluster
- A Claude OAuth token

## Usage

1. Edit `secret.yaml` and replace `<your-oauth-token>` with your actual token.
2. Edit `workspace.yaml` and replace the repo URL with your target repository.
3. Apply the resources:

```bash
kubectl apply -f secret.yaml
kubectl apply -f workspace.yaml
kubectl apply -f task.yaml
```

4. Watch the Task progress:

```bash
kubectl get tasks -w
```
