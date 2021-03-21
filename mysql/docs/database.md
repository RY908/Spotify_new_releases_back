# Table Definition
---

## Artist
| Field    | Type         | Null | Key | Default | Extra |
| ---- | ---- | ---- | ---- | ---- | ---- |
| artistId | varchar(100) | NO   | PRI | NULL    |       |
| name     | varchar(100) | NO   |     | NULL    |       |
| url      | varchar(200) | NO   |     | NULL    |       |
| iconUrl  | varchar(200) | NO   |     | NULL    |       |

##### Description
ArtistId, name, url, iconUrl are fetched from spotify api.

---
## ListenTo
| Field       | Type         | Null | Key | Default | Extra          |
| ---- | ---- | ---- | ---- | ---- | ---- |
| listenId    | int          | NO   | PRI | NULL    | auto_increment |
| userId      | varchar(100) | NO   | MUL | NULL    |                |
| artistId    | varchar(100) | NO   | MUL | NULL    |                |
| listenCount | int          | NO   |     | NULL    |                |
| timestamp   | datetime     | NO   |     | NULL    |                |
| ifFollowing | tinyint(1)   | NO   |     | NULL    |                |

##### Description
if the artist is the one user follows, ifFollowing will be set to 1 and listenCount will be set to 1000.
timestamp represents the last time user listened to the artist.

---

## User
| Field         | Type         | Null | Key | Default | Extra |
| ---- | ---- | ---- | ---- | ---- | ---- |
| userId        | varchar(100) | NO   | PRI | NULL    |       |
| accessToken   | varchar(300) | NO   |     | NULL    |       |
| tokenType     | varchar(20)  | NO   |     | NULL    |       |
| refreshToken  | varchar(200) | NO   |     | NULL    |       |
| expiry        | datetime     | NO   |     | NULL    |       |
| playlistId    | varchar(100) | NO   |     | NULL    |       |
| ifRemixAdd    | tinyint(1)   | YES  |     | 1       |       |
| ifAcousticAdd | tinyint(1)   | YES  |     | 1       |       |

##### Description
UserId, accessToken, refreshToken, refreshToken, expiry are fetched from spotify api.