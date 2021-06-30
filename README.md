# project-live-app

![LOGO](https://user-images.githubusercontent.com/1651333/123939585-da0f7780-d9ca-11eb-8985-247d4fa96ed5.png)

This is a project for our GoSchool Project Live Batch 4.

## File structure

```
- api
  |- controllers
    |- routes.go
    |- base.go
    |- users_controller.go
    |- health_controller.go
    |- business_controller.go
    |- category_controller.go
    |- comment_controller.go
  |- models
    |- category.go
    |- user.go
    |- business.go
    |- comment.go
  |- auth (security)
    |- token.go
  |- middlewares
    |- middlewares.go
- tests
- client
  |- controllers
    |- business_controller.go
    |- category_controller.go
    |- home_controller.go
    |- users_controller.go
  |- public
    |- style.css
  |- routes
    |- routes.go
  |- templates
  |- main.go (client start http server)
- go.mod
- go.sum
- .env
- main.go (rest api start http server)
- README.md
```

## Overview
This app is to help HBBs be known to the users of the app.

Registered users are able to create business listings and leave their comments / ratings. Businesses' geographic coordinates, latitude and longitude will be derived from their address via geocoding and shown on an embeded map.

Public users can browse and search for businesses.

## Database Flow
[Miro Board](https://miro.com/app/board/o9J_l-xAAp8=/)

## Team Members

- [Alvin Yeoh](https://github.com/xenodus)
- [Khairul](https://github.com/mofodox)
- [Koh Shao Wei](https://github.com/ksw95)
- [Sherman Lum](https://github.com/Smbsg) (Withdrawn)
