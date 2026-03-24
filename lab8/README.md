# lab8

This lab starts an ActiveMQ consumer and listens to both a queue and a topic.

## Run

```bash
go run ./cmd/app
```

Use another config file if needed:

```bash
go run ./cmd/app -config custom-config.json
```

## Configuration

The program reads broker and destination settings from `config.json`.

```json
{
  "broker": {
    "host": "127.0.0.1",
    "port": 61613,
    "username": "admin",
    "password": "admin",
    "host_name": "localhost"
  },
  "destinations": {
    "queue": "lab8.queue",
    "topic": "lab8.topic"
  }
}
```

- `broker.host`: ActiveMQ server IP or host name
- `broker.port`: STOMP port, default `61613` when omitted
- `broker.username`: optional login user
- `broker.password`: optional login password
- `broker.host_name`: optional STOMP host header, useful when broker expects a specific virtual host
- `destinations.queue`: queue name, auto-expanded to `/queue/<name>` if needed
- `destinations.topic`: topic name, auto-expanded to `/topic/<name>` if needed

## Behavior

- Connects to ActiveMQ on startup
- Subscribes to configured queue and topic
- Prints every received message to console
- Stops cleanly on `Ctrl+C`

## Else
如果要用docker來測試，以下提供二種：

### ARTEMIS
這是下一代的activeMQ，效能上都有所提升，但是，還沒去研究。
`docker run -d --name artemis -p 8161:8161 -p 61616:61616 -p 61613:61613 -e ARTEMIS_USER=admin -e ARTEMIS_PASSWORD=admin -e ANONYMOUS_LOGIN=true apache/activemq-artemis:latest-alpine`

### ACTIVEMQ
經典款，舊系統可能都使用這個。
`docker run -d --name activemq-classic -p 8161:8161 -p 61616:61616 -p 61613:61613 symptoma/activemq:5.18.3`