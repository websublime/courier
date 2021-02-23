# Courier Publish/Subscribe

Courier is an publish/subscriber socket to broadcast messages to every susbscription. All communication must be authenticated with a valid JWT token. 

## Rest endpoints

Courier provide endpoints to manage the service. With this endpoints you will be able to create audience, obtain a signed url to use socket service and a hook enpoint to send messages to pub/sub.

Create a audience to group content and connections communications. Audience (aud) should be define in your JWT token. All endpoints are secure with a valid JWT. If you need a system to administrate user sign, etc please consider use [gotrue](https://github.com/websublime/gotrue) system.

- POST /v1/audience

Creates and audience from your jwt token

Returns
```json
{
  "id": "11111111-2222-3333-4444-5555555555555",
  "name": "my-audience.com",
  "created_at": "2016-05-15T19:53:12.368652374-07:00",
  "updated_at": "2016-05-15T19:53:12.368652374-07:00"
}
```

- POST /v1/topic

Creates and a topic for a audience (jwt will give audience)

```json {
  {
    "topic": "system/events"
  }
}

Returns
```json
{
  "id": "11111111-2222-3333-4444-5555555555555",
  "topic": "system/events",
  "audience": "...",
  "audienceId": "11111111-2222-3333-4444-5555555555555",
  "created_at": "2016-05-15T19:53:12.368652374-07:00",
  "updated_at": "2016-05-15T19:53:12.368652374-07:00"
}
```

- GET /v1/topic

Get topics for a audience (jwt will give audience)

Returns
```json
[{
  "id": "11111111-2222-3333-4444-5555555555555",
  "topic": "system/events",
  "audience": "...",
  "audienceId": "11111111-2222-3333-4444-5555555555555",
  "created_at": "2016-05-15T19:53:12.368652374-07:00",
  "updated_at": "2016-05-15T19:53:12.368652374-07:00"
}]
```

- POST /v1/hook

Hook endpoint to publish events (message has to be object)

```json
{
  "topic": "system/events",
  "action": "publish",
  "message": {
    "type": "database"
    ...
  }
}
```

Returns same object from body

- GET /v1/sign

Obtain a signed to url to connect to socket. JWT audience will be use to generate your secure connection.

Returns
```json
{
  "key": "8p2uzHww5V6tMdOMJKyxdUDYDzw6DyE2yucduoKtM_HwpXNR0JvHK_KanL9xX1bTlkjJ3lj5eZ2hOr0x-OQZfOFQcd2n4ukUQ2Tde1dXkLOvAMBbpJt14Fe",
  "url": "ws://localhost:8883/ws?token=8p2uzHww5V6tMdOMJKyxdUDYDzw6DyE2yucduoKtM_HwpXNR0JvHK_KanL9xX1bTlkjJ3lj5eZ2hOr0x-OQZfOFQcd2n4ukUQ2Tde1dXkLOvAMBbpJt14Fe"
}
```

## PUB/SUB

Pub/Sub allows you to subscribe, unsubscribe, publish and broadcasts messages. After get a valid signed url the followings signatures for each type are:

- Subscribe
```json
{
  "action": "subscribe",
  "topic": "system/events"
}
```

- Unsubscribe
```json
{
  "action": "unsubscribe",
  "topic": "system/events"
}
```

- Publish
```json
{
  "action": "publish",
  "topic": "system/events",
  "message": {} //always object
}
```

- Broadcast
```json
{
  "action": "broadcast",
  "message": {} //always object
}
```
