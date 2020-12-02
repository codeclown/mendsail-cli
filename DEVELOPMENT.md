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

Testing against local instance:

```bash
export MENDSAIL_BASE_URL=http://localhost:3000/v1
make build && ./bin/mendsail send \
  --api-key XXXX-XXXX-XXXX-XXXX \
  --to you@example.com \
  --subject "Error in cronjob.sh" \
  --list "item 1" "item 2"
```
