# dreamstats

### Building binary for linux

```bash
set GOARCH=amd64
set GOOS=linux
go build .
```


```bash
set MONGO_DB=dreamstats
set MONGO_URI=mongodb://root:dreamstats@localhost:27017
```
MONGO_DB=dreamstats
MONGO_URI=mongodb://root:dreamstats@localhost:27017

Player name
format 
limit 
all

```bash
sudo service dreamstats status
sudo service dreamstats start
sudo service dreamstats stop
sudo service dreamstats restart
```


https://dev.to/hackmamba/build-a-rest-api-with-golang-and-mongodb-gin-gonic-version-269m

https://www.digitalocean.com/community/tutorials/how-to-secure-mongodb-on-ubuntu-20-04

<!-- Utkarsh Rode8:08 PM
team1_name
team2_name
team1_icon
team2_icon
innings { number, team, runs, wickets, overs }
event_name
match_type
match_number
date
outcome_string

Utkarsh Rode8:15 PM
match{
event_name
match_type
match_number
date
outcome_string
team{
  name
  icon
  innings{
    {runs, overs, wickets},
    {runs, overs, wickets}
  }
},
team{
  name
  icon
  innings{
    {runs, overs, wickets},
    {runs, overs, wickets}
  }
}
} -->