# Liquibase Lock Release CronJob

## Requirements
- [Go](https://golang.org/dl/)
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/get-started)
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
## Usage
### Build Container Image
```bash
make container
```
### Deploy Cron Job to the current K8s cluster
```bash
make deploy
```
&nbsp;
## Environment variables
In order to make the CronJob work, you need to have a K8s Secret deployed with the name `liquibase-lock-release-secret` and the following variables defined:
```yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: liquibase-lock-release-secret
  namespace: example
type: Opaque
data:
  DB_USER: <base64-encoded-db-user>
  DB_PASS: <base64-encoded-db-pass>
  DB_URL: <base64-encoded-db-url>
  DB_NAME: <base64-encoded-schema-name>
```
&nbsp;
## Configuration
- Environment variable `MAX_LOCK_TIME` determines what is the allowed time for the lock in minutes. The variable needs to be a string (example: 3m, 5m, etc...)
- Some assumptions are made in the code such as DB port being `3306` and DB engine `mysql`