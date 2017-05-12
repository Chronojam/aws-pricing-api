#!/bin/bash
docker run --name some-mysql -e MYSQL_ROOT_PASSWORD=my-secret-pw -d -p 3306:3306 mysql:5.6 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

