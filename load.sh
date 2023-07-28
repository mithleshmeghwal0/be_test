#!/bin/bash

# The URL of your Go HTTP server
SERVER_URL="http://localhost:8080/api/user/v1/users"
ACCESS_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJEYXRhIjoiN2Q0MTI2MjgtMGJmYi00NTdjLWFiYjgtMTFmZjQ5MGZlMDM0IiwiaXNzIjoidGVzdCIsInN1YiI6InNvbWVib2R5IiwiZXhwIjoxNjkzMTAyNTA5LCJuYmYiOjE2OTA1MTA1MTYsImlhdCI6MTY5MDUxMDUxNiwianRpIjoiZDk0ZDYxNDAtMzc0My00MmJmLThkODQtMGQ0ZDcwNTQ0OTVlIn0.pd4Z1S8wsQQCDf2b-uIgLw3azb4thY0RCH7-FCdzp_U"

# Perform load testing with ab
ab -k -n 1000 -c 100 -p testdata.txt -T application/json -H "Accept-Encoding: gzip, deflate" -H "Authorization: Bearer ${ACCESS_TOKEN}" "${SERVER_URL}"

