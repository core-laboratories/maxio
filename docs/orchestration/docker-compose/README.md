# MaxIO with Docker Compose

The example runs four MaxIO servers behind Nginx on one Docker host. It is useful for local development and testing.

From the repository root, build the image:

```sh
docker build -t maxio:local .
```

Then start the services:

```sh
cd docs/orchestration/docker-compose
docker compose up -d
```

The S3 API is exposed on port 9000 and the console on port 9001. Edit [docker-compose.yaml](docker-compose.yaml) to change credentials, storage mounts, or the number of nodes. The MaxIO configuration environment variables still use their established `MINIO_` names.
