# Fix All GitHub Repo Issues

This sample matches the [Auto-fix GitHub issues with TaskSpawner](../../README.md#auto-fix-github-issues-with-taskspawner) example from the README.

It creates a TaskSpawner that polls a GitHub repository for open issues labeled `bug` and automatically creates a Task for each one.

## Prerequisites

- Axon installed on your cluster
- A Claude OAuth token
- A GitHub token with access to the target repository

## Usage

1. Edit `secret-oauth.yaml` and replace `<your-oauth-token>` with your Claude OAuth token.
2. Edit `secret-github-token.yaml` and replace `<your-github-token>` with your GitHub token.
3. Edit `workspace.yaml` and replace the repo URL with your target repository.
4. Apply the resources:

```bash
kubectl apply -f secret-oauth.yaml
kubectl apply -f secret-github-token.yaml
kubectl apply -f workspace.yaml
kubectl apply -f taskspawner.yaml
```

5. Watch the TaskSpawner:

```bash
kubectl get taskspawners -w
```

The TaskSpawner polls for new issues every 5 minutes and creates a Task for each one matching the filter.
