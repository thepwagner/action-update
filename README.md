# action-update

Golang toolkit for building GitHub Actions that analyze and perform dependency updates.

To build an update action:
* Implement the [Updater](https://github.com/thepwagner/action-update/blob/86bcb48b1d7395e207073cc60b789a6da677bde0/updater/updater.go#L41-L45) interface.
* Extend `updateaction.Environment` into a struct that implements `updater.Factory` [sample](https://github.com/thepwagner/action-update-go/blob/68c5ba279d625e2fd526e9dfa4919612960f2158/gomodules/env.go#L8-L15)
* Write a `main()` that passes the environment to `handlers.ParseAndHandle()` [sample](https://github.com/thepwagner/action-update-go/blob/68c5ba279d625e2fd526e9dfa4919612960f2158/main.go#L16-L22)

### Implementations

* https://github.com/thepwagner/action-update-docker
* https://github.com/thepwagner/action-update-dockerurl
* https://github.com/thepwagner/action-update-go
* WIP: https://github.com/thepwagner/action-update-twirp
* Abandoned: https://github.com/thepwagner/action-update-brewformula
