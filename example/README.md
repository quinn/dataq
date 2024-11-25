# DataQ Example

This example demonstrates how to use DataQ to collect and organize data using plugins.

## Directory Structure

```
example/
├── config.yaml     # Plugin configuration file
├── data/          # Sample data directory
│   ├── document1.txt
│   └── notes.docx
```

## Configuration

The `config.yaml` file specifies which plugins to use and their configuration:

- `enabled`: Whether the plugin should be active
- `config`: Plugin-specific configuration parameters

### File Scanner Plugin Configuration

```yaml
plugins:
  - id: filescan
    enabled: true
    config:
      root_path: "./data"  # Directory to scan
      include_patterns:    # File patterns to include
        - "*.txt"
        - "*.pdf"
        - "*.doc*"
      exclude_patterns:    # File patterns to exclude
        - "*.tmp"
        - "~$*"
```

### Gmail Plugin Configuration

The Gmail plugin requires OAuth2 credentials from Google Cloud Console. Set up your credentials using the auth helper in `cmd/plugins/gmail/auth/`:

```yaml
plugins:
  - id: gmail
    enabled: true
    config:
      credentials_json: |
        {
          "installed": {
            "client_id": "YOUR_CLIENT_ID",
            "client_secret": "YOUR_CLIENT_SECRET",
            "redirect_uris": ["http://localhost:8085"]
          }
        }
      token_json: |
        {
          "access_token": "YOUR_ACCESS_TOKEN",
          "token_type": "Bearer",
          "refresh_token": "YOUR_REFRESH_TOKEN",
          "expiry": "2024-01-01T00:00:00Z"
        }
```

## Running the Example

1. Make sure you have the DataQ binary in your PATH
2. Configure your plugins in `config.yaml`
3. Run DataQ:
   ```bash
   dataq extract
   ```

The plugins will:
- File Scanner: Scan the `data` directory for matching files
- Gmail: Extract emails using your OAuth credentials
