package schema

// Task represents the database model for tasks.
type Task struct {
	UID            string            `json:"uid"`
	ClaimTimestamp int64             `json:"claim_timestamp"`
	PluginID       string            `json:"plugin_id"`
	ActionUID      string            `json:"action_uid"`
	DataItemUID    string            `json:"data_item_uid"`
	Status         string            `json:"status"`
	Error          string            `json:"error"`
	CreatedAt      int64             `json:"created_at"`
	UpdatedAt      int64             `json:"updated_at"`
	PluginConfig   map[string]string `json:"plugin_config"`
}

// Email represents a generic email entity within dataq.
type Email struct {
	From        string `json:"from"`
	To          string `json:"to"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	Date        string `json:"date"`
	MessageID   string `json:"message_id"`
	ThreadID    string `json:"thread_id"`
	InReplyTo   string `json:"in_reply_to"`
	References  string `json:"references"`
	Cc          string `json:"cc"`
	Bcc         string `json:"bcc"`
	Attachments string `json:"attachments"`
	MimeType    string `json:"mime_type"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	HTML        string `json:"html"`
	Text        string `json:"text"`
}

// FinancialTransaction represents a financial transaction entity within dataq.
type FinancialTransaction struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Category    string `json:"category"`
	Account     string `json:"account"`
	Subcategory string `json:"subcategory"`
	Notes       string `json:"notes"`
	Type        string `json:"type"`
}
