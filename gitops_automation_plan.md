# ⚔️ Operation: Full-Auto GitOps

This plan implements a professional CI/CD bridge between your Go backend and your Kubernetes staging environment.

## 🎯 Objectives
*   **Eliminate `:latest`**: Stop using insecure and caching-prone tags.
*   **Full Automation**: Pushing code to `devix-backend` must update `staging-ops` automatically.
*   **Audit Trail**: Every deployment in ArgoCD will link directly to a specific GitHub commit.

## 🛠️ The Tactical Setup

### 1. GitHub Secrets (Backend Repo)
You need to add these to your `devix-backend` repository settings:
*   `DOCKER_USERNAME`: `riteshsingh193`
*   `DOCKER_PASSWORD`: Your Docker Hub token.
*   `OPS_REPO_TOKEN`: A Personal Access Token (PAT) with "Contents: Write" access to `staging-ops`.

### 2. The GitHub Action Workflow
Create `.github/workflows/deploy-staging.yml` in your **Backend** repo:

```yaml
name: Staging Deployment

on:
  push:
    branches: [ main ]

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and Push
        uses: docker/build-push-action@v5
        with:
          push: true
          tags: riteshsingh193/devix-staging:${{ github.sha }}

      - name: Update Staging Ops
        run: |
          git clone https://x-access-token:${{ secrets.OPS_REPO_TOKEN }}@github.com/neutron420/staging-ops.git
          cd staging-ops
          sed -i "s|image: riteshsingh193/devix-staging:.*|image: riteshsingh193/devix-staging:${{ github.sha }}|g" staging/backend/deployment.yml
          git config user.name "Devix CI/CD"
          git config user.email "ci@devix.app"
          git add .
          git commit -m "deploy: update staging image to ${{ github.sha }}"
          git push origin main
```

## 📈 The Result
*   **Push Code** to Backend.
*   **Wait 2 mins**.
*   **ArgoCD** detects the change in `staging-ops`.
*   **New Pods** roll out with the exact code you just wrote.

---
> [!IMPORTANT]
> This setup ensures that your `staging-ops` repo is ALWAYS the source of truth for what is currently running in your cluster.
