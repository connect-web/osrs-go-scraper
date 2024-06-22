# Low Latency

## What's the main focus of Low Latency?
- Live OSRS Hiscores data for datasets in Machine Learning.

## What is low latency?
- Go project that scrapes the Runescape Hiscores at high speeds.

## Why do we need it?
- There are multiple bot farms achieving hundreds of millions of experience in one skill.
- Low latency is the backend for an API to easily find these accounts and ban the bots earlier.



## What's the inspiration behind Low Latency?

I have worked with the owners of Bot Detector, and they shared their dataset on the period of November 2023 to Jan 2024. 
This dataset added more insightful features to my time series dataset on the same period which helped achieve high accuracy ML Models while also giving insights to the locations of bot behaviour.




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

## Environmental variables

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

### Username collector

**Navigate to directory:**
src/nameFinder

Add a file named `proxies.txt` in the directory.
`proxies.txt` must be a file containing proxies with the following:
```ip:port:user:pass``` in each line, one line per proxy.

```shell
go run name.go limitfinder.go
```

or you can build it and run it in a new directory without any source files

```shell
go build name.go limitfinder.go
```

then you can simply run with
```shell
./name
```

### Stats finder

**Navigate to directory:**
src/statFinder

Add a file named `proxies.txt` in the directory.
`proxies.txt` must be a file containing proxies with the following:
```ip:port:user:pass``` in each line, one line per proxy.

```shell
go run stats.go
```

or you can build it and run it in a new directory without any source files

```shell
go build stats.go
```

then you can simply run with
```shell
./name
```

