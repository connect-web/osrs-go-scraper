# High performance OSRS Hiscore scraper

## Purpose
- This project is used for downloading player data from the Old-School Runescape hiscores and supports the [Low Latency website](https://low-latency.co.uk) by providing the backend with data.

[Low Latency Source](https://github.com/connect-web/Low-Latency-API)



# Setup
## Creating docker instance

```shell
docker run -d --name runescape-database \
  -p 127.0.0.1:5432:5432  \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=database \
  -e POSTGRES_PASSWORD=mypassword  \
  timescale/timescaledb-ha:pg16
```
### Required:
replace ```mypassword``` with a secure password

### Optional
You can change the external port which is the first ```5432``` if you are already running a PostgreSQL database on your host.

```postgres``` can be replaced with a desired username
```database``` can be replaced with a desired database
```127.0.0.1``` can be replaced with a host to expose the container outside of localhost.

- After creating your docker instance you must create the extensions & tables in the ```./sql``` folder in the respective files.
- Coming soon: Docker file with setup for database

## Environment variables

Then set your environmental variables for the Go scraper to connect with your docker database.

edit ```~/.bashrc``` and set these to the top of your file.

```shell
export lowLatencyUser="postgres"
export lowlatencyPassword="mypassword" # plz change this
export lowLatencyHost="localhost"
export lowLatencyPort="5432"
export lowLatencyDatabase="database"
```

# Running

### Username finder

**Navigate to directory:**
```src/cmd/name_finder```

Add a file named `proxies.txt` in the directory.
`proxies.txt` must be a file containing proxies with the following:
```ip:port:user:pass``` in each line, one line per proxy.

```shell
go run namefinder.go
```

or you can build it and run it in a new directory without any source files

```shell
go build namefinder.go
```

then you can run with
```shell
./namefinder
```

### Stats finder

**Navigate to directory:**
```src/cmd/player_live```

Add a file named `proxies.txt` in the directory.
`proxies.txt` must be a file containing proxies with the following:
```ip:port:user:pass``` in each line, one line per proxy.

```shell
go run playerlive.go
```

or you can build it and run it in a new directory without any source files

```shell
go build playerlive.go
```

then you can simply run with
```shell
./playerlive
```

### Gains finder

**Navigate to directory:**
```src/cmd/player_gains```

Add a file named `proxies.txt` in the directory.
`proxies.txt` must be a file containing proxies with the following:
```ip:port:user:pass``` in each line, one line per proxy.

```shell
go run playergains.go
```

or you can build it and run it in a new directory without any source files

```shell
go build playergains.go
```

then you can simply run with
```shell
./playergains
```


