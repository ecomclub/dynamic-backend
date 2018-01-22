# dynamic-backend
Map E-Com Plus Storefront dynamic backend with Go and Redis

# Technology stack
+ [Go](https://golang.org/) 1.9.x
+ [Redis](https://redis.io/) 3
+ Redis client for Golang https://github.com/go-redis/redis

# Setting up
For security, we recommend to download and install the app as root,
and let the files owned by `root:root` as default.

```bash
sudo git clone https://github.com/ecomclub/dynamic-backend.git
cd dynamic-backend
sudo go build main.go
```

Start application with CLI arguments:
+ Root directory to static files
+ HTTP/TCP port
+ Optional log file path

Example:

```bash
./main /var/www :3000 /var/log/app.log
```
