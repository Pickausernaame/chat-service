// Code generated by ent, DO NOT EDIT.

package message

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"

	"github.com/Pickausernaame/chat-service/internal/types"
)

const (
	// Label holds the string label denoting the message type in the database.
	Label = "message"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldChatID holds the string denoting the chat_id field in the database.
	FieldChatID = "chat_id"
	// FieldProblemID holds the string denoting the problem_id field in the database.
	FieldProblemID = "problem_id"
	// FieldAuthorID holds the string denoting the author_id field in the database.
	FieldAuthorID = "author_id"
	// FieldIsVisibleForClient holds the string denoting the is_visible_for_client field in the database.
	FieldIsVisibleForClient = "is_visible_for_client"
	// FieldIsVisibleForManager holds the string denoting the is_visible_for_manager field in the database.
	FieldIsVisibleForManager = "is_visible_for_manager"
	// FieldBody holds the string denoting the body field in the database.
	FieldBody = "body"
	// FieldCheckedAt holds the string denoting the checked_at field in the database.
	FieldCheckedAt = "checked_at"
	// FieldIsBlocked holds the string denoting the is_blocked field in the database.
	FieldIsBlocked = "is_blocked"
	// FieldIsService holds the string denoting the is_service field in the database.
	FieldIsService = "is_service"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeProblem holds the string denoting the problem edge name in mutations.
	EdgeProblem = "problem"
	// EdgeChat holds the string denoting the chat edge name in mutations.
	EdgeChat = "chat"
	// Table holds the table name of the message in the database.
	Table = "messages"
	// ProblemTable is the table that holds the problem relation/edge.
	ProblemTable = "messages"
	// ProblemInverseTable is the table name for the Problem entity.
	// It exists in this package in order to avoid circular dependency with the "problem" package.
	ProblemInverseTable = "problems"
	// ProblemColumn is the table column denoting the problem relation/edge.
	ProblemColumn = "problem_id"
	// ChatTable is the table that holds the chat relation/edge.
	ChatTable = "messages"
	// ChatInverseTable is the table name for the Chat entity.
	// It exists in this package in order to avoid circular dependency with the "chat" package.
	ChatInverseTable = "chats"
	// ChatColumn is the table column denoting the chat relation/edge.
	ChatColumn = "chat_id"
)

// Columns holds all SQL columns for message fields.
var Columns = []string{
	FieldID,
	FieldChatID,
	FieldProblemID,
	FieldAuthorID,
	FieldIsVisibleForClient,
	FieldIsVisibleForManager,
	FieldBody,
	FieldCheckedAt,
	FieldIsBlocked,
	FieldIsService,
	FieldCreatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultIsVisibleForClient holds the default value on creation for the "is_visible_for_client" field.
	DefaultIsVisibleForClient bool
	// DefaultIsVisibleForManager holds the default value on creation for the "is_visible_for_manager" field.
	DefaultIsVisibleForManager bool
	// BodyValidator is a validator for the "body" field. It is called by the builders before save.
	BodyValidator func(string) error
	// DefaultIsBlocked holds the default value on creation for the "is_blocked" field.
	DefaultIsBlocked bool
	// DefaultIsService holds the default value on creation for the "is_service" field.
	DefaultIsService bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() types.MessageID
)

// OrderOption defines the ordering options for the Message queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByChatID orders the results by the chat_id field.
func ByChatID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldChatID, opts...).ToFunc()
}

// ByProblemID orders the results by the problem_id field.
func ByProblemID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProblemID, opts...).ToFunc()
}

// ByAuthorID orders the results by the author_id field.
func ByAuthorID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAuthorID, opts...).ToFunc()
}

// ByIsVisibleForClient orders the results by the is_visible_for_client field.
func ByIsVisibleForClient(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsVisibleForClient, opts...).ToFunc()
}

// ByIsVisibleForManager orders the results by the is_visible_for_manager field.
func ByIsVisibleForManager(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsVisibleForManager, opts...).ToFunc()
}

// ByBody orders the results by the body field.
func ByBody(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBody, opts...).ToFunc()
}

// ByCheckedAt orders the results by the checked_at field.
func ByCheckedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCheckedAt, opts...).ToFunc()
}

// ByIsBlocked orders the results by the is_blocked field.
func ByIsBlocked(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsBlocked, opts...).ToFunc()
}

// ByIsService orders the results by the is_service field.
func ByIsService(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsService, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByProblemField orders the results by problem field.
func ByProblemField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProblemStep(), sql.OrderByField(field, opts...))
	}
}

// ByChatField orders the results by chat field.
func ByChatField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newChatStep(), sql.OrderByField(field, opts...))
	}
}

func newProblemStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProblemInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ProblemTable, ProblemColumn),
	)
}

func newChatStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ChatInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ChatTable, ChatColumn),
	)
}
