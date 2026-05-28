# tally

```docker
docker run --rm --name tally-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=yourpassword \
  -e POSTGRES_DB=tally \
  -p 5432:5432 \
  postgres:latest
```
