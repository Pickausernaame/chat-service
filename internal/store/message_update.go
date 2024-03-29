// Code generated by ent, DO NOT EDIT.

package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"

	"github.com/Pickausernaame/chat-service/internal/store/message"
	"github.com/Pickausernaame/chat-service/internal/store/predicate"
	"github.com/Pickausernaame/chat-service/internal/types"
)

// MessageUpdate is the builder for updating Message entities.
type MessageUpdate struct {
	config
	hooks    []Hook
	mutation *MessageMutation
}

// Where appends a list predicates to the MessageUpdate builder.
func (mu *MessageUpdate) Where(ps ...predicate.Message) *MessageUpdate {
	mu.mutation.Where(ps...)
	return mu
}

// SetInitialRequestID sets the "initial_request_id" field.
func (mu *MessageUpdate) SetInitialRequestID(ti types.RequestID) *MessageUpdate {
	mu.mutation.SetInitialRequestID(ti)
	return mu
}

// SetNillableInitialRequestID sets the "initial_request_id" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableInitialRequestID(ti *types.RequestID) *MessageUpdate {
	if ti != nil {
		mu.SetInitialRequestID(*ti)
	}
	return mu
}

// ClearInitialRequestID clears the value of the "initial_request_id" field.
func (mu *MessageUpdate) ClearInitialRequestID() *MessageUpdate {
	mu.mutation.ClearInitialRequestID()
	return mu
}

// SetIsVisibleForManager sets the "is_visible_for_manager" field.
func (mu *MessageUpdate) SetIsVisibleForManager(b bool) *MessageUpdate {
	mu.mutation.SetIsVisibleForManager(b)
	return mu
}

// SetNillableIsVisibleForManager sets the "is_visible_for_manager" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableIsVisibleForManager(b *bool) *MessageUpdate {
	if b != nil {
		mu.SetIsVisibleForManager(*b)
	}
	return mu
}

// SetBody sets the "body" field.
func (mu *MessageUpdate) SetBody(s string) *MessageUpdate {
	mu.mutation.SetBody(s)
	return mu
}

// SetCheckedAt sets the "checked_at" field.
func (mu *MessageUpdate) SetCheckedAt(t time.Time) *MessageUpdate {
	mu.mutation.SetCheckedAt(t)
	return mu
}

// SetNillableCheckedAt sets the "checked_at" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableCheckedAt(t *time.Time) *MessageUpdate {
	if t != nil {
		mu.SetCheckedAt(*t)
	}
	return mu
}

// ClearCheckedAt clears the value of the "checked_at" field.
func (mu *MessageUpdate) ClearCheckedAt() *MessageUpdate {
	mu.mutation.ClearCheckedAt()
	return mu
}

// SetIsBlocked sets the "is_blocked" field.
func (mu *MessageUpdate) SetIsBlocked(b bool) *MessageUpdate {
	mu.mutation.SetIsBlocked(b)
	return mu
}

// SetNillableIsBlocked sets the "is_blocked" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableIsBlocked(b *bool) *MessageUpdate {
	if b != nil {
		mu.SetIsBlocked(*b)
	}
	return mu
}

// SetIsService sets the "is_service" field.
func (mu *MessageUpdate) SetIsService(b bool) *MessageUpdate {
	mu.mutation.SetIsService(b)
	return mu
}

// SetNillableIsService sets the "is_service" field if the given value is not nil.
func (mu *MessageUpdate) SetNillableIsService(b *bool) *MessageUpdate {
	if b != nil {
		mu.SetIsService(*b)
	}
	return mu
}

// Mutation returns the MessageMutation object of the builder.
func (mu *MessageUpdate) Mutation() *MessageMutation {
	return mu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (mu *MessageUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, MessageMutation](ctx, mu.sqlSave, mu.mutation, mu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (mu *MessageUpdate) SaveX(ctx context.Context) int {
	affected, err := mu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (mu *MessageUpdate) Exec(ctx context.Context) error {
	_, err := mu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (mu *MessageUpdate) ExecX(ctx context.Context) {
	if err := mu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (mu *MessageUpdate) check() error {
	if v, ok := mu.mutation.InitialRequestID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "initial_request_id", err: fmt.Errorf(`store: validator failed for field "Message.initial_request_id": %w`, err)}
		}
	}
	if v, ok := mu.mutation.Body(); ok {
		if err := message.BodyValidator(v); err != nil {
			return &ValidationError{Name: "body", err: fmt.Errorf(`store: validator failed for field "Message.body": %w`, err)}
		}
	}
	if _, ok := mu.mutation.ProblemID(); mu.mutation.ProblemCleared() && !ok {
		return errors.New(`store: clearing a required unique edge "Message.problem"`)
	}
	if _, ok := mu.mutation.ChatID(); mu.mutation.ChatCleared() && !ok {
		return errors.New(`store: clearing a required unique edge "Message.chat"`)
	}
	return nil
}

func (mu *MessageUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := mu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(message.Table, message.Columns, sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID))
	if ps := mu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if mu.mutation.AuthorIDCleared() {
		_spec.ClearField(message.FieldAuthorID, field.TypeUUID)
	}
	if value, ok := mu.mutation.InitialRequestID(); ok {
		_spec.SetField(message.FieldInitialRequestID, field.TypeUUID, value)
	}
	if mu.mutation.InitialRequestIDCleared() {
		_spec.ClearField(message.FieldInitialRequestID, field.TypeUUID)
	}
	if value, ok := mu.mutation.IsVisibleForManager(); ok {
		_spec.SetField(message.FieldIsVisibleForManager, field.TypeBool, value)
	}
	if value, ok := mu.mutation.Body(); ok {
		_spec.SetField(message.FieldBody, field.TypeString, value)
	}
	if value, ok := mu.mutation.CheckedAt(); ok {
		_spec.SetField(message.FieldCheckedAt, field.TypeTime, value)
	}
	if mu.mutation.CheckedAtCleared() {
		_spec.ClearField(message.FieldCheckedAt, field.TypeTime)
	}
	if value, ok := mu.mutation.IsBlocked(); ok {
		_spec.SetField(message.FieldIsBlocked, field.TypeBool, value)
	}
	if value, ok := mu.mutation.IsService(); ok {
		_spec.SetField(message.FieldIsService, field.TypeBool, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, mu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{message.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	mu.mutation.done = true
	return n, nil
}

// MessageUpdateOne is the builder for updating a single Message entity.
type MessageUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *MessageMutation
}

// SetInitialRequestID sets the "initial_request_id" field.
func (muo *MessageUpdateOne) SetInitialRequestID(ti types.RequestID) *MessageUpdateOne {
	muo.mutation.SetInitialRequestID(ti)
	return muo
}

// SetNillableInitialRequestID sets the "initial_request_id" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableInitialRequestID(ti *types.RequestID) *MessageUpdateOne {
	if ti != nil {
		muo.SetInitialRequestID(*ti)
	}
	return muo
}

// ClearInitialRequestID clears the value of the "initial_request_id" field.
func (muo *MessageUpdateOne) ClearInitialRequestID() *MessageUpdateOne {
	muo.mutation.ClearInitialRequestID()
	return muo
}

// SetIsVisibleForManager sets the "is_visible_for_manager" field.
func (muo *MessageUpdateOne) SetIsVisibleForManager(b bool) *MessageUpdateOne {
	muo.mutation.SetIsVisibleForManager(b)
	return muo
}

// SetNillableIsVisibleForManager sets the "is_visible_for_manager" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableIsVisibleForManager(b *bool) *MessageUpdateOne {
	if b != nil {
		muo.SetIsVisibleForManager(*b)
	}
	return muo
}

// SetBody sets the "body" field.
func (muo *MessageUpdateOne) SetBody(s string) *MessageUpdateOne {
	muo.mutation.SetBody(s)
	return muo
}

// SetCheckedAt sets the "checked_at" field.
func (muo *MessageUpdateOne) SetCheckedAt(t time.Time) *MessageUpdateOne {
	muo.mutation.SetCheckedAt(t)
	return muo
}

// SetNillableCheckedAt sets the "checked_at" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableCheckedAt(t *time.Time) *MessageUpdateOne {
	if t != nil {
		muo.SetCheckedAt(*t)
	}
	return muo
}

// ClearCheckedAt clears the value of the "checked_at" field.
func (muo *MessageUpdateOne) ClearCheckedAt() *MessageUpdateOne {
	muo.mutation.ClearCheckedAt()
	return muo
}

// SetIsBlocked sets the "is_blocked" field.
func (muo *MessageUpdateOne) SetIsBlocked(b bool) *MessageUpdateOne {
	muo.mutation.SetIsBlocked(b)
	return muo
}

// SetNillableIsBlocked sets the "is_blocked" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableIsBlocked(b *bool) *MessageUpdateOne {
	if b != nil {
		muo.SetIsBlocked(*b)
	}
	return muo
}

// SetIsService sets the "is_service" field.
func (muo *MessageUpdateOne) SetIsService(b bool) *MessageUpdateOne {
	muo.mutation.SetIsService(b)
	return muo
}

// SetNillableIsService sets the "is_service" field if the given value is not nil.
func (muo *MessageUpdateOne) SetNillableIsService(b *bool) *MessageUpdateOne {
	if b != nil {
		muo.SetIsService(*b)
	}
	return muo
}

// Mutation returns the MessageMutation object of the builder.
func (muo *MessageUpdateOne) Mutation() *MessageMutation {
	return muo.mutation
}

// Where appends a list predicates to the MessageUpdate builder.
func (muo *MessageUpdateOne) Where(ps ...predicate.Message) *MessageUpdateOne {
	muo.mutation.Where(ps...)
	return muo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (muo *MessageUpdateOne) Select(field string, fields ...string) *MessageUpdateOne {
	muo.fields = append([]string{field}, fields...)
	return muo
}

// Save executes the query and returns the updated Message entity.
func (muo *MessageUpdateOne) Save(ctx context.Context) (*Message, error) {
	return withHooks[*Message, MessageMutation](ctx, muo.sqlSave, muo.mutation, muo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (muo *MessageUpdateOne) SaveX(ctx context.Context) *Message {
	node, err := muo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (muo *MessageUpdateOne) Exec(ctx context.Context) error {
	_, err := muo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (muo *MessageUpdateOne) ExecX(ctx context.Context) {
	if err := muo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (muo *MessageUpdateOne) check() error {
	if v, ok := muo.mutation.InitialRequestID(); ok {
		if err := v.Validate(); err != nil {
			return &ValidationError{Name: "initial_request_id", err: fmt.Errorf(`store: validator failed for field "Message.initial_request_id": %w`, err)}
		}
	}
	if v, ok := muo.mutation.Body(); ok {
		if err := message.BodyValidator(v); err != nil {
			return &ValidationError{Name: "body", err: fmt.Errorf(`store: validator failed for field "Message.body": %w`, err)}
		}
	}
	if _, ok := muo.mutation.ProblemID(); muo.mutation.ProblemCleared() && !ok {
		return errors.New(`store: clearing a required unique edge "Message.problem"`)
	}
	if _, ok := muo.mutation.ChatID(); muo.mutation.ChatCleared() && !ok {
		return errors.New(`store: clearing a required unique edge "Message.chat"`)
	}
	return nil
}

func (muo *MessageUpdateOne) sqlSave(ctx context.Context) (_node *Message, err error) {
	if err := muo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(message.Table, message.Columns, sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID))
	id, ok := muo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`store: missing "Message.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := muo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, message.FieldID)
		for _, f := range fields {
			if !message.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("store: invalid field %q for query", f)}
			}
			if f != message.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := muo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if muo.mutation.AuthorIDCleared() {
		_spec.ClearField(message.FieldAuthorID, field.TypeUUID)
	}
	if value, ok := muo.mutation.InitialRequestID(); ok {
		_spec.SetField(message.FieldInitialRequestID, field.TypeUUID, value)
	}
	if muo.mutation.InitialRequestIDCleared() {
		_spec.ClearField(message.FieldInitialRequestID, field.TypeUUID)
	}
	if value, ok := muo.mutation.IsVisibleForManager(); ok {
		_spec.SetField(message.FieldIsVisibleForManager, field.TypeBool, value)
	}
	if value, ok := muo.mutation.Body(); ok {
		_spec.SetField(message.FieldBody, field.TypeString, value)
	}
	if value, ok := muo.mutation.CheckedAt(); ok {
		_spec.SetField(message.FieldCheckedAt, field.TypeTime, value)
	}
	if muo.mutation.CheckedAtCleared() {
		_spec.ClearField(message.FieldCheckedAt, field.TypeTime)
	}
	if value, ok := muo.mutation.IsBlocked(); ok {
		_spec.SetField(message.FieldIsBlocked, field.TypeBool, value)
	}
	if value, ok := muo.mutation.IsService(); ok {
		_spec.SetField(message.FieldIsService, field.TypeBool, value)
	}
	_node = &Message{config: muo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, muo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{message.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	muo.mutation.done = true
	return _node, nil
}
