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

- `credentials_path`: Path to your OAuth client credentials JSON file
- `token_path`: Path to your OAuth token JSON file

Example configuration:
```yaml
plugins:
  - id: gmail
    enabled: true
    config:
      credentials_path: "./gmail/credentials.json"
      token_path: "./gmail/token.json"
```

## Usage

The plugin will extract all emails from the authenticated Gmail account, including:
- Raw email content (RFC822 format)
- Email metadata (headers)
- Message IDs
- Timestamps

The plugin handles pagination automatically and will retrieve all available emails.
