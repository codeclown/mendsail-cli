# Development

## Tests

```bash
make test
```

## Build and run

```bash
make build && ./bin/mendsail send \
  --api-key XXXX-XXXX-XXXX-XXXX \
  --to you@example.com \
  --subject "Error in cronjob.sh" \
  --heading "Data processing failed" \
  --paragraph "Log output:" \
  --code-block "foobar"
```
