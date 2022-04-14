# API (This document in not updated)
---

## GET /login
#### Description
Proceed spotify authentication process.
After accessing, user will be redirected to spotify login page. 
#### Response
| code | Description |
| ---- | ---- |
| 302 | Redirect to spotify login page |

## GET /callback
#### Description
After authentication, Spotify Accounts service redirect to /callback. 
In /callback, user's information (userId, following artists...) is inserted in database.

#### Response
| code | Description |
| ---- | ---- |
| 302 | if successful, redirect to /user/(access token) page (frontend). if unsuccessful, redirect to / page (frontend) |


## Get /user
#### Description
Using cookie given, returns artists whose songs will be added to next playlist update. 

#### Response
```
{
    "artists": [
        {
            "ArtistId": "abcd1234"
            "IconUrl": "https:..."
            "Name": "artist name1"
            "Url": "https:..."
        },
        {
            "ArtistId": "abcd1235"
            "IconUrl": "https:..."
            "Name": "artist name2"
            "Url": "https:..."
        },
        ...
    ]
}
```
Still working on how to change json key's first character from uppercase letter to lowercase letter.

| code | Description |
| ---- | ---- |
| 200 |  |
| 400 | cannot get cookie |
| 401 | cookie is set but never logged in |
| 500 | server error |

## POST /delete
#### Description
user can stop adding songs by artists the user selected.
In this process, first check if the request has set cookie and after that, delete relations between user and artist from database.

#### Request
```
{
    "artistsId": [
        "xxx", "yyy", ...    
    ]
}
```

#### Response
```
{
    "artists": [
        {
            "ArtistId": "abcd1234"
            "IconUrl": "https:..."
            "Name": "artist name1"
            "Url": "https:..."
        },
        {
            "ArtistId": "abcd1235"
            "IconUrl": "https:..."
            "Name": "artist name2"
            "Url": "https:..."
        },
        ...
    ]
}
```

| code | Description |
| ---- | ---- |
| 200 |  |
| 400 | cannot get cookie |
| 401 | cookie is set but never logged in |
| 500 | server error |

## GET /setting
#### Description
Returns user's settings (if remix/acoustic songs will be added to playlist or not).

#### Response
```
{
    "ifRemixAdd": true or false,
    "ifAcousticAdd: true or false
}
```
| code | Description |
| ---- | ---- |
| 200 |  |
| 400 | cannot get cookie |
| 401 | cookie is set but never logged in |
| 500 | server error |

## POST /setting/save
#### Description
Set user's setting.

#### Response
| code | Description |
| ---- | ---- |
| 200 |  |
| 400 | cannot get cookie |
| 401 | cookie is set but never logged in |
| 500 | server error |