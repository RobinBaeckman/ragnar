# ragnar

Small api with small architecture.

Installation:
1. Clone repo to src
'''
git clone https://github.com/RobinBaeckman/ragnar.git ~/go/src/ragnar
'''

2. Import dependencies 
```
cd ~/go/src/ragnar
export GO111MODULE=on
go mod init . 
go mod vendor
```

2. Export env variables
```
export HOST="localhost" &&\
export PORT="3000" &&\
export LOG_PREFIX="something-api: " &&\
export MYSQL_HOST="127.0.0.1" &&\
export MYSQL_USER="ruser" &&\
export MYSQL_PASS="secret" &&\
export MYSQL_DB="ragnar_db" &&\
export REDIS_HOST="127.0.0.1" &&\
export REDIS_PORT="6379" &&\
export COOKIE_NAME="cookie"
```

3. if you're using docker setup docker mysql
```
cd ~/go/src/ragnar &&\
docker-compose up -d
``` 

else if you're running a mysql server, just make sure you have the right privilege to add tables
```
cd ~/go/src/ragnar &&\
mysql -uruser -psecret -h127.0.0.1 < configs/dump.sql

```

4. Run program
```
go run main.go
```

Usage:

Create new user
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -d @tests/endpoints/create_user.json -X POST http://localhost:3000/v1/users | python -m json.tool
```

Login
```
curl -v -c /tmp/cookie-jar.txt -d @tests/endpoints/login.json http://localhost:3000/v1/login
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
curl -v -i -L -b /tmp/cookie-jar.txt -d @tests/endpoints/update_user.json -X PUT http://localhost:3000/v1/users/{id}

```

Delete user
```
// TODO: use something else than python for showing json because the status code is not shown
curl -v -i -L -b /tmp/cookie-jar.txt -X DELETE http://localhost:3000/v1/users/{id}

```
