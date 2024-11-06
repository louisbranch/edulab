# EduLab

## Running locally

By default, the project runs on port `8080`: [http://localhost:8080/edulab/](http://localhost:8080/edulab/).

### With Docker Compose
```
docker compose -f server-compose.yaml up -d --build server
```

### With Go

```
go mod tidy
go run cmd/server/main.go
```

Two databases are supported [SQLite](https://github.com/mattn/go-sqlite3) (default) and PostgreSQL.
For PostgreSQL, you have to set the following environmental variables:
```
POSTGRES_USER=%your db user%
POSTGRES_PASSWORD=%your db password%
POSTGRES_HOSTNAME=%your db hostname%
POSTGRES_DB=%your db name%
```

## Adding a new language

Add the language go to [translations.go](translations/translations.go)

For Spanish, for example:
```
//go:generate gotext -srclang=en update -out=catalog.go -lang=en,pt-BR,es github.com/louisbranch/edulab/cmd/server
```

Then run:
```
go generate translations/translations.go
```

A new file for the language will be created at [translations/locales](translations/locales/).

Copy the source file for the one containing the translated messages:
```
cp translations/locales/es/out.gotext.json translations/locales/es/messages.gotext.json
```

After translating the messages in `messages.gotext.json`, run the generator again to update the language catalog:
```
go generate translations/translations.go
```

Finally, for the language to appear on the website, add the new language to the [server](/web/server/i18n.go).
