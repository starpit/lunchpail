For example:

```shell
KUBECONFIG=~/.kube/config docker buildx build --push \
    --platform=linux/arm64/v8,linux/amd64 \
        --tag ghcr.io/project-codeflare/workerpool-worker-ubuntu-python-openroad:0.0.1 .
```