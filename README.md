# rolf 

Small api with small architecture.

## Installation:
1. Clone repo to src
```
git clone https://github.com/RobinBaeckman/ragnar.git ~/go/src/ragnar
```

2. Import dependencies 
```
cd ~/go/src/ragnar
export GO111MODULE=on
go mod init . 
go mod vendor
```

3. Export env variables
```
export HOST="localhost" &&\
export PORT="3000" &&\
export LOG_PREFIX="rolf-api" &&\
export MYSQL_HOST="127.0.0.1" &&\
export MYSQL_PORT="3306" &&\
export MYSQL_USER="rolf" &&\
export MYSQL_PASS="secret" &&\
export MYSQL_DB="rolf_db" &&\
export REDIS_HOST="127.0.0.1" &&\
export REDIS_PORT="6379" &&\
export COOKIE_NAME="cookie" &&\
export JWT_KEY="secret" &&\
export PROTO="http" &&\
export SMTP_HOST="0.0.0.0" &&\
export SMTP_PORT="1025" &&\
export EMAIL="server@mail.com"
```

4. if you're using docker setup docker mysql
```
cd ~/go/src/ragnar &&\
docker-compose up -d
``` 

else if you're running a mysql server, just make sure you have the right privilege to add tables
```
cd ~/go/src/ragnar &&\
mysql -uroot -psecret -h127.0.0.1 < configs/dump.sql

```

5. Run program
```
go run main.go
```

## Usage:

Register new user
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -d @scripts/tests/curl/create_user.json -X POST http://localhost:3000/v1/register | python -m json.tool
```

Register new admin 
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -i -L -b /tmp/cookie-jar.txt  -d @scripts/tests/curl/create_admin.json -X POST http://localhost:3000/v1/register/admin | python -m json.tool
```

Login
```
curl -v -c /tmp/cookie-jar.txt -d @scripts/tests/curl/login.json http://localhost:3000/v1/login
```

Login Admin
```
curl -v -c /tmp/cookie-jar.txt -d @scripts/tests/curl/login_admin.json http://localhost:3000/v1/login
```

Read user
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -i -L -b /tmp/cookie-jar.txt -X GET http://localhost:3000/v1/users/{id}

```

Read users
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -i -L -b /tmp/cookie-jar.txt -X GET http://localhost:3000/v1/users

```

Update user
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -i -L -b /tmp/cookie-jar.txt -d @scripts/tests/curl/update_user.json -X PUT http://localhost:3000/v1/users/{id}

```

Delete user
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -i -L -b /tmp/cookie-jar.txt -X DELETE http://localhost:3000/v1/users/{id}

```

Run load test
```
cd ragnar-api/scripts/load_test
k6 run perfect-field.js
```

Run benchmark without testing. CreateUser example.
```
go test -bench=BenchmarkCreateUser -run=XXX -cpuprofile=cpu.out -benchtime=100x
go tool pprof rest.test cpu.out
top50
```
