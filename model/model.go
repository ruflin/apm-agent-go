package model

import (
	"net/http"
	"net/url"
	"time"
)

// Service represents the service handling transactions being traced.
type Service struct {
	// Name is the immutable name of the service.
	Name string `json:"name"`

	// Version is the version of the service, if it has one.
	Version string `json:"version,omitempty"`

	// Environment is the name of the service's environment, if it has
	// one, e.g. "production" or "staging".
	Environment string `json:"environment,omitempty"`

	// Agent holds information about the Elastic APM agent tracing this
	// service's transactions.
	Agent Agent `json:"agent"`

	// Framework holds information about the service's framework, if any.
	Framework *Framework `json:"framework,omitempty"`

	// Language holds information about the programming language in which
	// the service is written.
	Language *Language `json:"language,omitempty"`

	// Runtime holds information about the programming language runtime
	// running this service.
	Runtime *Runtime `json:"runtime,omitempty"`
}

// Agent holds information about the Elastic APM agent.
type Agent struct {
	// Name is the name of the Elastic APM agent, e.g. "Go".
	Name string `json:"name"`

	// Version is the version of the Elastic APM agent, e.g. "1.0.0".
	Version string `json:"version"`
}

// Framework holds information about the framework (typically web)
// used by the service.
type Framework struct {
	// Name is the name of the framework.
	Name string `json:"name"`

	// Version is the version of the framework.
	Version string `json:"version"`
}

// Language holds information about the programming language used.
type Language struct {
	// Name is the name of the programming language.
	Name string `json:"name"`

	// Version is the version of the programming language.
	Version string `json:"version,omitempty"`
}

// Runtime holds information about the programming language runtime.
type Runtime struct {
	// Name is the name of the programming language runtime.
	Name string `json:"name"`

	// Version is the version of the programming language runtime.
	Version string `json:"version"`
}

// System represents the system (operating system and machine) running the
// service.
type System struct {
	// Architecture is the system's hardware architecture.
	Architecture string `json:"architecture,omitempty"`

	// Hostname is the system's hostname.
	Hostname string `json:"hostname,omitempty"`

	// Platform is the system's platform, or operating system name.
	Platform string `json:"platform,omitempty"`
}

// Process represents an operating system process.
type Process struct {
	// Pid is the process ID.
	Pid int `json:"pid"`

	// Ppid is the parent process ID, if known.
	Ppid *int `json:"ppid,omitempty"`

	// Title is the title of the process.
	Title string `json:"title,omitempty"`

	// Argv holds the command line arguments used to start the process.
	Argv []string `json:"argv,omitempty"`
}

// Transaction represents a transaction handled by the service.
type Transaction struct {
	// ID holds the hex-formatted UUID of the transaction.
	ID string `json:"id"`

	// Name holds the name of the transaction.
	Name string `json:"name"`

	// Type identifies the service-domain specific type of the request,
	// e.g. "request" or "backgroundjob".
	Type string `json:"type"`

	// Timestamp holds the time at which the transaction started.
	Timestamp time.Time `json:"-"`

	// Duration records how long the transaction took to complete.
	Duration time.Duration `json:"-"`

	// Result holds the result of the transaction, e.g. the status code
	// for HTTP requests.
	Result string `json:"result,omitempty"`

	// Context holds contextual information relating to the transaction.
	Context *Context `json:"context,omitempty"`

	// Sampled indicates that the transaction was sampled, and
	// includes all available information. Non-sampled transactions
	// omit Context and Spans.
	//
	// If Sampled is unspecified (nil), it is equivalent to setting
	// it to true.
	Sampled *bool `json:"sampled,omitempty"`

	// SpanCount holds statistics on spans within a transaction.
	SpanCount *SpanCount `json:"span_count,omitempty"`

	// Spans holds the transaction's spans.
	Spans []*Span `json:"spans,omitempty"`
}

// SpanCount holds statistics on spans within a transaction.
type SpanCount struct {
	// Dropped holds statistics on dropped spans within a transaction.
	Dropped *SpanCountDropped `json:"dropped,omitempty"`
}

// SpanCountDropped holds statistics on dropped spans.
type SpanCountDropped struct {
	// Total holds the total number of spans dropped by the
	// agent within a transaction.
	Total int `json:"total"`
}

// Span represents a span within a transaction.
type Span struct {
	// Name holds the name of the span.
	Name string `json:"name"`

	// Start is the start time of the span, as a duration relative to the
	// containing transaction's timestamp.
	Start time.Duration `json:"-"`

	// Duration holds the duration of the span.
	Duration time.Duration `json:"-"`

	// Type identifies the service-domain specific type of the span,
	// e.g. "db.postgresql.query".
	Type string `json:"type"`

	// ID holds an identifier for the span, unique within its
	// containing transaction.
	ID *int64 `json:"id,omitempty"`

	// Parent holds the identifier of the parent span, if any.
	Parent *int64 `json:"parent,omitempty"`

	// Context holds contextual information relating to the span.
	Context *SpanContext `json:"context,omitempty"`

	// Stacktrace holds stack frames corresponding to the span.
	Stacktrace []StacktraceFrame `json:"stacktrace,omitempty"`
}

// SpanContext holds contextual information relating to the span.
type SpanContext struct {
	// Database holds contextual information for database
	// operation spans.
	Database *DatabaseSpanContext `json:"db,omitempty"`
}

// DatabaseSpanContext holds contextual information for database
// operation spans.
type DatabaseSpanContext struct {
	// Instance holds the database instance name.
	Instance string `json:"instance,omitempty"`

	// Statement holds the database statement (e.g. query).
	Statement string `json:"statement,omitempty"`

	// Type holds the database type. For any SQL database,
	// this should be "sql"; for others, the lower-cased
	// database category, e.g. "cassandra", "hbase", "redis".
	Type string `json:"type,omitempty"`

	// User holds the username used for database access.
	User string `json:"user,omitempty"`
}

// Context holds contextual information relating to a transaction or error.
type Context struct {
	// Request holds details of the HTTP request relating to the
	// transaction or error, if relevant.
	Request *Request `json:"request,omitempty"`

	// Response holds details of the HTTP response relating to the
	// transaction or error, if relevant.
	Response *Response `json:"response,omitempty"`

	// User holds details of the authenticated user relating to the
	// transaction or error, if relevant.
	User *User `json:"user,omitempty"`

	// Custom holds arbitrary additional metadata.
	Custom map[string]interface{} `json:"custom,omitempty"`

	// Tags holds user-defined key/value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// User holds information about an authenticated user.
type User struct {
	// Username holds the username of the user.
	Username string `json:"username,omitempty"`

	// ID identifies the user, e.g. a primary key. This may be
	// a string or number.
	ID interface{} `json:"id,omitempty"`

	// Email holds the email address of the user.
	Email string `json:"email,omitempty"`
}

// Error represents an error occurring in the service.
type Error struct {
	// Timestamp holds the time at which the error occurred.
	Timestamp time.Time `json:"timestamp"`

	// ID holds a hex-formatted UUID for the error.
	ID string `json:"id,omitempty"`

	// TransactionID holds the UUID of the transaction to which
	// this error relates, if any.
	TransactionID string `json:"-"`

	// Culprit holds the name of the function which
	// produced the error.
	Culprit string `json:"culprit,omitempty"`

	// Context holds contextual information relating to the error.
	Context *Context `json:"context,omitempty"`

	// Exception holds details of the exception (error or panic)
	// to which this error relates.
	Exception *Exception `json:"exception,omitempty"`

	// Log holds additional information added when logging the error.
	Log *Log `json:"log,omitempty"`
}

// Exception represents an exception: an error or panic.
type Exception struct {
	// Message holds the error message.
	Message string `json:"message"`

	// Code holds the error code. This may be a number or a string.
	Code interface{} `json:"code,omitempty"`

	// Type holds the type of the exception.
	Type string `json:"type,omitempty"`

	// Module holds the exception type's module namespace.
	Module string `json:"module,omitempty"`

	// Attributes holds arbitrary exception-type specific attributes.
	Attributes map[string]interface{} `json:"attributes,omitempty"`

	// Stacktrace holds stack frames corresponding to the exception.
	Stacktrace []StacktraceFrame `json:"stacktrace,omitempty"`

	// Handled indicates whether or not the error was caught and handled.
	Handled bool `json:"handled"`
}

// StacktraceFrame describes a stack frame.
type StacktraceFrame struct {
	// AbsolutePath holds the absolute path of the source file for the
	// stack frame.
	AbsolutePath string `json:"abs_path,omitemmpty"`

	// File holds the base filename of the source file for the stack frame.
	File string `json:"filename"`

	// Line holds the line number of the source for the stack frame.
	Line int `json:"lineno"`

	// Column holds the column number of the source for the stack frame.
	Column *int `json:"colno,omitempty"`

	// Module holds the module to which the frame belongs. For Go, we
	// use the package path (e.g. "net/http").
	Module string `json:"module,omitempty"`

	// Function holds the name of the function to which the frame belongs.
	Function string `json:"function,omitempty"`

	// LibraryFrame indicates whether or not the frame corresponds to
	// library or user code.
	LibraryFrame bool `json:"library_frame,omitempty"`

	// ContextLine holds the line of source code to which the frame
	// corresponds.
	ContextLine string `json:"context_line,omitempty"`

	// PreContext holds zero or more lines of source code preceding the
	// line corresponding to the frame.
	PreContext []string `json:"pre_context,omitempty"`

	// PostContext holds zero or more lines of source code proceeding the
	// line corresponding to the frame.
	PostContext []string `json:"post_context,omitempty"`

	// Vars holds local variables for this stack frame.
	Vars map[string]interface{} `json:"vars,omitempty"`
}

// Log holds additional information added when logging an error.
type Log struct {
	// Message holds the logged error message.
	Message string `json:"message"`

	// Level holds the severity of the log record.
	Level string `json:"level,omitempty"`

	// LoggerName holds the name of the logger used.
	LoggerName string `json:"logger_name,omitempty"`

	// ParamMessage holds a parameterized message,  e.g.
	// "Could not connect to %s". The string is not interpreted,
	// but may be used for grouping errors.
	ParamMessage string `json:"param_message,omitempty"`

	// Stacktrace holds stack frames corresponding to the error.
	Stacktrace []StacktraceFrame `json:"stacktrace,omitempty"`
}

// Request represents an HTTP request.
type Request struct {
	// URL is the request URL.
	URL URL `json:"url"`

	// Method holds the HTTP request method.
	Method string `json:"method"`

	// Headers holds the request headers.
	Headers *RequestHeaders `json:"headers,omitempty"`

	// Body holds the request body, if body capture is enabled.
	Body *RequestBody `json:"body,omitempty"`

	// HTTPVersion holds the HTTP version of the request.
	HTTPVersion string `json:"http_version,omitempty"`

	// Cookies holds the parsed cookies.
	Cookies []*http.Cookie `json:"-"`

	// Env holds environment information passed from the
	// web framework to the request handler.
	Env map[string]interface{} `json:"env,omitempty"`

	// Socket holds transport-level information.
	Socket *RequestSocket `json:"socket,omitempty"`
}

// RequestBody holds a request body.
//
// Exactly one of Raw or Form must be set.
type RequestBody struct {
	// Raw holds the raw body content.
	Raw string

	// Form holds the form data from POST, PATCH, or PUT body parameters.
	Form url.Values
}

// RequestHeaders holds a limited subset of HTTP request headers.
type RequestHeaders struct {
	// ContentType holds the content-type header.
	ContentType string `json:"content-type,omitempty"`

	// Cookie holds the cookies sent with the request,
	// delimited by semi-colons.
	Cookie string `json:"cookie,omitempty"`

	// UserAgent holds the user-agent header.
	UserAgent string `json:"user-agent,omitempty"`
}

// RequestSocket holds transport-level information relating to an HTTP request.
type RequestSocket struct {
	// Encrypted indicates whether or not the request was sent
	// as an SSL/HTTPS request.
	Encrypted bool `json:"encrypted,omitempty"`

	// RemoteAddress holds the remote address for the request.
	RemoteAddress string `json:"remote_address,omitempty"`
}

// URL represents a request URL.
type URL struct {
	// Full is the full URL, e.g.
	// "https://example.com:443/search/?q=elasticsearch#top".
	Full string `json:"full,omitempty"`

	// Protocol is the scheme of the URL, e.g. "https".
	Protocol string `json:"protocol,omitempty"`

	// Hostname is the hostname for the URL, e.g. "example.com".
	Hostname string `json:"hostname,omitempty"`

	// Port is the port number in the URL, e.g. "443".
	Port string `json:"port,omitempty"`

	// Path is the path of the URL, e.g. "/search".
	Path string `json:"pathname,omitempty"`

	// Search is the query string of the URL, e.g. "q=elasticsearch".
	Search string `json:"search,omitempty"`

	// Hash is the fragment for references, e.g. "top" in the
	// URL example provided for Full.
	Hash string `json:"hash,omitempty"`
}

// Response represents an HTTP response.
type Response struct {
	// StatusCode holds the HTTP response status code.
	StatusCode int `json:"status_code,omitempty"`

	// Headers holds the response headers.
	Headers *ResponseHeaders `json:"headers,omitempty"`

	// HeadersSent indicates whether or not headers were sent
	// to the client.
	HeadersSent *bool `json:"headers_sent,omitempty"`

	// Finished indicates whether or not the response was finished.
	Finished *bool `json:"finished,omitempty"`
}

// ResponseHeaders holds a limited subset of HTTP respponse headers.
type ResponseHeaders struct {
	// ContentType holds the content-type header.
	ContentType string `json:"content-type,omitempty"`
}
