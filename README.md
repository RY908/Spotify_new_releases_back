# [NewReleases Server](https://newreleases.tk)

---

## Overview
Server side implementation of [NewReleases](https://newreleases.tk)
This application collects your listening history on spotify and makes a playlist based on the artists you listened to and follow.

## Description
After authenticated by spotify api, your listening history will be collected every 20 minutes and will be inserted in database.
Every weekend, a new playlist named "New Release" is made based on your listening history and artists you follow.

## Usage
1. Set environment variables specified in ```docker-compose.yml```.
2. ```docker-compose build```
3. ```docker-compose run```