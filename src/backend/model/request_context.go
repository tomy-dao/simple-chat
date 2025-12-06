package model

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// RequestContext wraps context.Context with type-safe accessors for request-specific values
type RequestContext struct {
	ctx      context.Context
	Token    string
	UserID   uint
	SessionID string
	span     trace.Span
}

// NewRequestContext creates a new RequestContext from a context.Context
// Extracts token, user_id, session_id, and tracer from the context
func NewRequestContext(ctx context.Context) *RequestContext {
	reqCtx := &RequestContext{
		ctx: ctx,
	}

	// Extract token
	if token, ok := ctx.Value("token").(string); ok {
		reqCtx.Token = token
	}

	// Extract user_id
	if userID, ok := ctx.Value("user_id").(float64); ok {
		reqCtx.UserID = uint(userID)
	} else if userID, ok := ctx.Value("user_id").(uint); ok {
		reqCtx.UserID = userID
	}

	// Extract session_id
	if sessionID, ok := ctx.Value("session_id").(string); ok {
		reqCtx.SessionID = sessionID
	}

	// Extract span from context (OpenTelemetry)
	reqCtx.span = trace.SpanFromContext(ctx)

	return reqCtx
}

// Context returns the underlying context.Context
func (r *RequestContext) Context() context.Context {
	return r.ctx
}

// WithToken sets the token in the context and returns a new RequestContext
func (r *RequestContext) WithToken(token string) *RequestContext {
	ctx := context.WithValue(r.ctx, "token", token)
	return &RequestContext{
		ctx:       ctx,
		Token:     token,
		UserID:    r.UserID,
		SessionID: r.SessionID,
		span:      r.span,
	}
}

// WithUserID sets the user_id in the context and returns a new RequestContext
func (r *RequestContext) WithUserID(userID uint) *RequestContext {
	ctx := context.WithValue(r.ctx, "user_id", userID)
	return &RequestContext{
		ctx:       ctx,
		Token:     r.Token,
		UserID:    userID,
		SessionID: r.SessionID,
		span:      r.span,
	}
}

// WithSessionID sets the session_id in the context and returns a new RequestContext
func (r *RequestContext) WithSessionID(sessionID string) *RequestContext {
	ctx := context.WithValue(r.ctx, "session_id", sessionID)
	return &RequestContext{
		ctx:       ctx,
		Token:     r.Token,
		UserID:    r.UserID,
		SessionID: sessionID,
		span:      r.span,
	}
}

// WithClaims sets token, user_id, and session_id from JWT claims
func (r *RequestContext) WithClaims(token string, userID uint, sessionID string) *RequestContext {
	ctx := context.WithValue(r.ctx, "token", token)
	ctx = context.WithValue(ctx, "user_id", userID)
	ctx = context.WithValue(ctx, "session_id", sessionID)
	return &RequestContext{
		ctx:       ctx,
		Token:     token,
		UserID:    userID,
		SessionID: sessionID,
		span:      r.span,
	}
}

// Span returns the OpenTelemetry span from context
func (r *RequestContext) Span() trace.Span {
	return r.span
}

// TraceID returns the trace ID from the span
func (r *RequestContext) TraceID() string {
	if r.span == nil || !r.span.SpanContext().IsValid() {
		return ""
	}
	return r.span.SpanContext().TraceID().String()
}

// SpanID returns the span ID from the span
func (r *RequestContext) SpanID() string {
	if r.span == nil || !r.span.SpanContext().IsValid() {
		return ""
	}
	return r.span.SpanContext().SpanID().String()
}

