# Architecture

## login
```
\(web)
↓
\login(api)
↓
spotify accounts service
(user authentication)
↓
\callback(api)
(if user's not in database, collect user's following artist and make playlist)
↓
\user/access_token(web)
```
