# ezclone
> `ezclone` is a tool which retrieves all your private and public github repositories plus repos you are contributing to and clones them to a given location
`ezclone` is written in go and has zero runtime dependencies (not even git itself).

## Usage
In order to talk to the API, `ezclone` needs an access token. One can be obtained [here](https://github.com/settings/tokens)
Put it inside the `GITHUB_ACCESSTOKEN` environment variable so the tool can access it.
Then, just run:
```
go run ezsync.go /path/to/where/you/want/your/backup
```