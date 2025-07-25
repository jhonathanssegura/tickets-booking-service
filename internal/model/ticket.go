package model

type Ticket struct {
	ID         string `dynamodbav:"id"`
	UserEmail  string `dynamodbav:"user_email"`
	EventID    string `dynamodbav:"event_id"`
	ReservedAt string `dynamodbav:"reserved_at"`
}
