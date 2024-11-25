# Gmail Plugin for DataQ

This plugin extracts emails and metadata from Gmail using the Gmail API.

## Setup

1. Create a Google Cloud Project and enable the Gmail API:
   - Go to [Google Cloud Console](https://console.cloud.google.com)
   - Create a new project or select an existing one
   - Search for and enable the Gmail API

2. Create OAuth credentials:
   - Go to APIs & Services > Credentials
   - Click "Create Credentials" > "OAuth client ID"
   - Choose "Desktop application" as the application type
   - Download the credentials JSON file as `credentials.json`

3. Get an OAuth token using the auth helper:
   ```bash
   cd auth
   go build -o auth.exe
   # Place your credentials.json in this directory
   ./auth.exe
   ```
   This will:
   - Open your browser for authentication
   - Save the token as `token.json`
   - Print both credentials and token in the format needed for config.yaml

## Configuration

The plugin requires two configuration parameters:

- `credentials_json`: The contents of your OAuth client credentials JSON file
- `token_json`: The OAuth token JSON obtained after authentication

Example configuration:
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

## Usage

The plugin will extract all emails from the authenticated Gmail account, including:
- Raw email content (RFC822 format)
- Email metadata (headers)
- Message IDs
- Timestamps

The plugin handles pagination automatically and will retrieve all available emails.
