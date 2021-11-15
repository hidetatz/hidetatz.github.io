type: input
timestamp: 2021-11-15 12:08:10
url: https://ably.com/blog/how-we-load-tested-control-api
lang: en
---

* They devide load test scenarios into three;
  - Typical users
  - Power users
    - Too many traffic, but don't have malcious intent
    - To simulate this properly without facing rate limiting, they use a reverse proxy (squid) to distribute the client IP address
  - Bad users/bots
    - Send too many traffic with bad intentions
