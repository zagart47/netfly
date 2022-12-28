# Netfly

Little conception of social network


## Description

The network has the ability to register, send messages to the user by name, view messages.

## Getting Started

### Dependencies

* Go 1.16+
* PostgreSQL
* Windows, Linux or Mac OS

### Using
First, you need to run the PostgreSQL server.
```
docker pull postgres
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD={PASSWORD} -e POSTGRES_USER={USERNAME} -d postgres
```

To run client need to create .env file in root project folder.

In .env file should have written strings:

```
DBHOST=postgres://{USERNAME}:{PASSWORD}@{DBHOST}:5432/postgres
TOKENLIFESPAN={How many hours your authorization token will be valid}
```
You must use Postman to work with the application.
To register a user, use:
```
POST: http://127.0.0.1:8080/api/register
```
and JSON raw:
```
{
    "username": "USERNAME",
    "password": "PASSWORD"
}
```
To login you need to use:
```
POST: http://127.0.0.1:8080/api/login
```
and JSON raw:
```
{
    "username": "USERNAME",
    "password": "PASSWORD"
}
```
Postman response shows you token.

To send messages, you need to use method, API and your token.
You need to add your token in "Authorization" tab in "Bearer Token".
```
POST: http://127.0.0.1:8080/api/message
```
and JSON raw:
```
{
    "recipient": "Patrick",
    "text": "Hello"
}
```
To read all your messages need to use:
```
GET: http://127.0.0.1:8080/api/message
```

## Authors

Artur Zagirov  
[@zagart47](https://t.me/zagart47)
