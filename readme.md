Clone project and run server:

```
go run ./main.go
```

Then use curl or postman.

Get zipcode from attached out.csv, service returns prefecture+city

Request:

```
{host}:8050/?z=9140011
```

Response:

```
{
"address": "福井県 敦賀市"
}
```
