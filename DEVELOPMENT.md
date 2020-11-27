# Development

## Build and run

```bash
make build && ./bin/mendsail send \
  --to you@example.com \
  --subject "Error in cronjob.sh" \
  --add-heading "Data processing failed" \
  --add-paragraph "Log output:" \
  --add-code-block "foobar"
```
