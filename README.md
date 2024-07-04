superchat-core
===

The first version monolith for superchat, a chat service for fun...because why not.

more coming soon...

# Development

## Create Migration

```shell
migrate create -ext sql -dir database/migrations migration_name
```

## Fix dirty database version after running invalid migration

```shell
go run cmd/fix.go -1
```