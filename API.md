# REST API

## Ping

### Healthcheck - `/v1/ping`

#### Request

```text
GET /v1/ping HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:5000
User-Agent: HTTPie/3.2.1
```

#### Response

```text
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Length: 30
Content-Type: application/json; charset=utf-8
Date: Sat, 08 Apr 2023 11:58:22 GMT
Vary: Accept-Encoding

"pong"

```

## Scan

### Scan Site - `/v1/scan`

#### Request

```text
POST /v1/scan HTTP/1.1
Accept: application/json, */*;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 31
Content-Type: application/json
Host: localhost:5000
User-Agent: HTTPie/3.2.1

{
    "base_url": "https://example.com"
}

```

#### Response

```text
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Length: 966
Content-Type: application/json; charset=utf-8
Date: Sat, 08 Apr 2023 12:13:25 GMT
Vary: Accept-Encoding

...omitted...

```
