Pre-Requisites
- docker
- docker-compose
- port 8080 should not be used
- ab -> apache benchmark

Run Test
- `docker build -t test_image -f DockerfileTest .`
- `docker run --rm --name testapp test_image`

Postman Collection
- `https://documenter.getpostman.com/view/18502288/2s9XxsUwBW``
- JWT token for requests . it should be placed inside `Authorization` Header
  Authorization: Bearer <token>
- Token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjoiN2Q0MTI2MjgtMGJmYi00NTdjLWFiYjgtMTFmZjQ5MGZlMDM0IiwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiZXhwIjoxNjkzMTAyNTA5LCJuYmYiOjE2OTA1MTA1MTYsImlhdCI6MTY5MDUxMDUxNiwianRpIjoiZDk0ZDYxNDAtMzc0My00MmJmLThkODQtMGQ0ZDcwNTQ0OTVlIn0.pd4Z1S8wsQQCDf2b-uIgLw3azb4thY0RCH7-FCdzp_U 


Run Service
- `docker-compose up --build`

Run load test
-  chmod +x load.sh;./load.sh

Structure of repo
- So i am choosing golang as programming language and postgres as database, there are many patterns to structure repo
  The pattern that i used , 
  `internal` -> this folder holds all the application logic.
  `internal\userservice` -> stores all api 
  `internal]\persist` -> interface for databases. Apis only interact with this and have no idea about database.
  `internal\sql` -> it holds the all db logic
  `internal\sql\migrations` -> it holds migration files, which will create tables, other required stuff for 

  `pkg` -> it holds all external package that we need for this project, such as sql, logger, jwt, environment.

Trade-offs
- If i had more time and had to build this service for production, 
  
  *First i will also build a `history service` which will store history of change to user record.
  For history, i will use `rabbitmq`, and push each change of user database in the queue, and history will listen to each event and push the user record in history.
  
  *Second, for the list api with pattern matching, i will use `rabbitmq`, and push data to Elastic search and build a `view` service which will serve list of users.

  