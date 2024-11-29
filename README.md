# go-url-shortener
URL Shortener API that helps shorten long URLs.

# Prerequisites

`Docker`

# Build docker image

```sh
docker build -t url-shortener .
```

# Run docker container

```sh
docker run -it --rm -v $(pwd):/workspace --workdir /workspace -p 8080:8080 url-shortener bash
```

# Run app in docker

```sh
mongod --fork --logpath /var/log/mongodb/mongod.log
go mod tidy
go run url-shortener
```

# Usage examples
```sh
curl -X POST -d '{"url": "http://someurl"}' localhost:8080/shorten
# {"_id":"674996324dc4add438c190e6","url":"http://someurl","shortCode":"fwVydA","createdAt":"2024-11-29T10:23:46Z","updatedAt":"2024-11-29T10:23:46Z"}
curl localhost:8080/shorten/fwVydA
# {"_id":"674996324dc4add438c190e6","url":"http://someurl","shortCode":"fwVydA","createdAt":"2024-11-29T10:23:46Z","updatedAt":"2024-11-29T10:23:46Z"}
curl localhost:8080/shorten/fwVydA/stats
# {"_id":"674996324dc4add438c190e6","url":"http://someurl","shortCode":"fwVydA","createdAt":"2024-11-29T10:23:46Z","updatedAt":"2024-11-29T10:23:46Z","accessCount":1}
curl -X PUT -d '{"url": "http://someotherurl"}' localhost:8080/shorten/fwVydA
# {"url":"http://someotherurl","shortCode":"fwVydA","createdAt":"2024-11-29T10:23:46Z","updatedAt":"2024-11-29T10:25:27Z"}
```

# Roadmap reference
https://roadmap.sh/projects/url-shortening-service
