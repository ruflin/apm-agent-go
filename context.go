package elasticapm

import "context"

// ContextWithSpan returns a copy of parent in which the given span
// is stored, associated with the key ContextSpanKey.
func ContextWithSpan(parent context.Context, s *Span) context.Context {
	return context.WithValue(parent, contextSpanKey{}, s)
}

// ContextWithTransaction returns a copy of parent in which the given
// transaction is stored, associated with the key ContextTransactionKey.
func ContextWithTransaction(parent context.Context, t *Transaction) context.Context {
	return context.WithValue(parent, contextTransactionKey{}, t)
}

// SpanFromContext returns the current Span in context, if any. The span must
// have been added to the context previously using either ContextWithSpan
// or SetSpanInContext.
func SpanFromContext(ctx context.Context) *Span {
	span, _ := ctx.Value(contextSpanKey{}).(*Span)
	return span
}

// TransactionFromContext returns the current Transaction in context, if any.
// The transaction must have been added to the context previously using either
// ContextWithTransaction or SetTransactionInContext.
func TransactionFromContext(ctx context.Context) *Transaction {
	tx, _ := ctx.Value(contextTransactionKey{}).(*Transaction)
	return tx
}

// StartSpan starts and returns a new Span within the sampled transaction
// and parent span in the context, if any, and returns the span along with
// a new context containing the span.
//
// If there is no transaction in the context, or it is not being sampled,
// StartSpan returns nil.
func StartSpan(ctx context.Context, name, spanType string) (*Span, context.Context) {
	tx := TransactionFromContext(ctx)
	if tx == nil || !tx.Sampled() {
		return nil, ctx
	}
	span := tx.StartSpan(name, spanType, SpanFromContext(ctx))
	return span, context.WithValue(ctx, contextSpanKey{}, span)
}

// CaptureError returns a new Error related to the sampled transaction
// present in the context, if any, and calls its SetException method
// with the given error. The Exception.Handled field will be set to true.
//
// If there is no transaction in the context, or it is not being sampled,
// CaptureError returns nil. As a convenience, if the provided error is
// nil, then CaptureError will also return nil.
func CaptureError(ctx context.Context, err error) *Error {
	if err == nil {
		return nil
	}
	tx := TransactionFromContext(ctx)
	if tx == nil || !tx.Sampled() {
		return nil
	}
	e := tx.tracer.NewError()
	e.SetException(err)
	e.Exception.Handled = true
	e.Transaction = tx
	return e
}

type contextSpanKey struct{}
type contextTransactionKey struct{}
