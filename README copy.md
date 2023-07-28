Pre-Requisites
- docker
- docker-compose
- port 8080 should not be used
- ab -> apache benchmark

Run Test
- `docker build -t test_image -f DockerfileTest .`
- `docker run --rm --name testapp test_image`

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

  