# go-url-shortener
URL Shortener API that helps shorten long URLs.

# Prerequisites

`Docker`

# Build docker image

```sh
cd go-url-shortener
docker build -t url-shortener .
```

# Run docker container

```sh
docker run -it --rm --user $(id -u):$(id -g) -v $(pwd):/workspace --workdir /workspace url-shortener
```

# Roadmap reference
https://roadmap.sh/projects/url-shortening-service
