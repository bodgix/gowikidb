# gowikidb

A GO implementation of a MySQL backed wiki engine.

# Deployment

Docker and docker-compose are used for the deployment

# Running

```bash
docker-compose build
docker-compose up
```

Point the browser at (http://localhost:8080/view/test)

# Creating the DATABASE

DB creation is not dockerized yet. Needs to be created manually.

Start the mysql container with

```bash
docker-compose up
```

Start a shell:

```bash
docker exec -ti gowikidb_mysql_1 mysql -u root -proot
```

Copy and paste `gowiki.sql` to the terminal. This will create the database
and the article tables.
