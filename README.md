# mendsail-cli

Command line tool to quickly send formatted email.

## Usage

```bash
mendsail send
  --api-key XXXX-XXXX-XXXX-XXXX
  --to [email_address]
  --subject [subject]
  [--add-heading <text>]
  [--add-paragraph <text>]
  [--add-code-block <text>]
  [--add-list [<list-items>]]
  [--add-list-item <text>]
```

### Examples

Send error notifications from scripts:

```bash
mendsail send \
  --api-key XXXX-XXXX-XXXX-XXXX \
  --to you@example.com \
  --subject "Error in cronjob.sh" \
  --add-heading "Data processing failed" \
  --add-paragraph "Log output:" \
  --add-code-block $(tail -n50 log.txt)
```

Send email with various types of content:

```bash
mendsail send \
  --api-key XXXX-XXXX-XXXX-XXXX \
  --to you@example.com \
  --subject "New domains ($(date))" \
  --add-heading "New domains ($(date))" \
  --add-paragraph "Today's new domains:" \
  --add-list "[example1.com](https://example1.com)\\n[example2.com](https://example2.com)"
```
