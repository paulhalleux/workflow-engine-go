### Create migrate script
```
 migrate create -ext sql -dir migrations -seq wd_def_add_enabledC
```

### Restart engine
```
docker-compose up -d --build engine
```