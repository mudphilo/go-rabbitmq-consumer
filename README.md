## How to run it

### 1. Setup environment variable
This application retrieves rabbitMQ credentials from environment variables. the following must be set for the application to run
```sh
RABBITMQ_HOST=127.0.0.1
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_PORT=5672
RABBITMQ_VHOST=test
JOB_QUEUES=FILES_QUEUE,SMS_QUEUE
```
