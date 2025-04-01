- migration user status from int to UserStatus [done]
- remove premium fields [done]
- change user role from string to UserPermission [done]
- implement user wallet [done]

-----

- implement voucher feature [done]
- implement monitor feature:
    + `/monitor/users-activities` [priority 1]
    + `/monitor/media-analysis` [priority 1] [done]
    + `/monitor/progress-analysis` [priority 1]
    + `/monitor/traffic-analysis` [priority 2]
- combine 3 features: audios, transcriptions, videos into medias
- only allow admin to create/update voucher
- refactor update fields in admin package
- refactor source code (combine feature to reduce number of folders in handler, service, repo)