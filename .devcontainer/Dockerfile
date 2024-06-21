FROM mcr.microsoft.com/vscode/devcontainers/go:1.22

# Install Dapr CLI
RUN wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | /bin/bash

WORKDIR /home/vscode
USER vscode
