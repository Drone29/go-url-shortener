FROM mongo

# Install deps
RUN apt update && apt install -y \
    curl \
    git \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# GO version
ARG GO_VERSION="1.23.2"

# Install GO
RUN rm -rf /usr/local/go \
    && curl -fsSL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go${GO_VERSION}.linux-amd64.tar.gz \
    && rm /tmp/go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH=$PATH:/usr/local/go/bin

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Set up MongoDB log directory
RUN mkdir -p /var/log/mongodb

EXPOSE 27017
EXPOSE 8080