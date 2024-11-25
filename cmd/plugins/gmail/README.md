# Gmail Plugin for DataQ

This plugin extracts emails and metadata from Gmail using the Gmail API.

## Setup

1. Create a Google Cloud Project and enable the Gmail API
2. Create OAuth 2.0 credentials (OAuth client ID) for a desktop application
3. Download the credentials JSON file
4. Run the authentication flow to get a token (you can use Google's OAuth 2.0 Playground)

## Configuration

The plugin requires two configuration parameters:

- `credentials_json`: The contents of your OAuth client credentials JSON file
- `token_json`: The OAuth token JSON obtained after authentication

## Usage

The plugin will extract all emails from the authenticated Gmail account, including:
- Raw email content (RFC822 format)
- Email metadata (headers)
- Message IDs
- Timestamps

The plugin handles pagination automatically and will retrieve all available emails.
