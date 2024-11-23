FROM mongo

# Install deps
RUN apt update && apt install -y \
    curl \
    git \
    sudo \
    && rm -rf /var/lib/apt/lists/*

# GO version
ARG GO_VERSION="1.23.2"
# Specify non-root user (may be overridden with '--env USER=<DESIRED USERNAME>')
ENV USER=shortener
# Add non-root user with default UID:GID
RUN useradd ${USER} -m -s /bin/bash

# Install GO
RUN rm -rf /usr/local/go \
    && curl -fsSL https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf /tmp/go${GO_VERSION}.linux-amd64.tar.gz \
    && rm /tmp/go${GO_VERSION}.linux-amd64.tar.gz

# Set Go environment variables
ENV GOPATH=/go
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

WORKDIR /workspace

EXPOSE 27017
EXPOSE 8080

# Run container under non-root user
# USER ${USER}

# Start from a Bash prompt
# CMD [ "mongod" ]


