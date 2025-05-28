- migration user status from int to UserStatus [done]
- remove premium fields [done]
- change user role from string to UserPermission [done]
- implement user wallet [done]

-----

- implement voucher feature [done]
- implement monitor feature:
    + `/monitor/users-activities` [priority 1] [pending]
    + `/monitor/media-analysis` [priority 1] [done]
    + `/monitor/progress-analysis` [priority 1] [done]
    + `/monitor/traffic-analysis` [priority 2] [done]
- return role when login [done]
- combine 3 features: audios, transcriptions, videos into medias [done]
- refactor mlvt_handler
- add traffic log to monitor service
- replace gin.H in admin monitor with response package and add swagger
- change monitor data type -> monitor media
- only allow admin to create/update voucher
- refactor update fields in admin package
- refactor source code (combine feature to reduce number of folders in handler, service, repo)

- premium user can get both free daily token and premium token [processing]
- premium user will receive token even not login in this day [processing]