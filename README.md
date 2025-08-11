# go-scrape

Little web scraper project to learn some Go

## Example Usage

Run the server with

```
go run main.go
```

Then send a request to `http://localhost:8080/scrape`

```
curl -X POST http://localhost:8080/scrape \
-H "Content-Type: application/json" \
-d '{"url": "https://job-boards.greenhouse.io/figma/jobs/5552522004?gh_jid=5552522004"}'
```
