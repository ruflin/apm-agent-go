package apmlambda

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/elastic/apm-agent-go"
	"github.com/elastic/apm-agent-go/model"
	"github.com/elastic/apm-agent-go/stacktrace"

	"github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

const (
	// TODO(axw) make this configurable via environment
	payloadLimit = 1024
)

var (
	// nonBlocking is passed to Tracer.Flush so it does not block functions.
	nonBlocking = make(chan struct{})

	// Globals below used during tracing, to avoid reallocating for each
	// invocation. Only one invocation will happen at a time.
	txContext     model.Context
	lambdaContext struct {
		RequestID       string `json:"request_id,omitempty"`
		Region          string `json:"region,omitempty"`
		XAmznTraceID    string `json:"x_amzn_trace_id,omitempty"`
		FunctionVersion string `json:"function_version,omitempty"`
		MemoryLimit     int    `json:"memory_limit,omitempty"`
		Request         string `json:"request,omitempty"`
		Response        string `json:"response,omitempty"`
	}
)

func init() {
	close(nonBlocking)
	txContext.Custom = map[string]interface{}{
		"lambda": &lambdaContext,
	}
	lambdaContext.FunctionVersion = lambdacontext.FunctionVersion
	lambdaContext.MemoryLimit = lambdacontext.MemoryLimitInMB
	lambdaContext.Region = os.Getenv("AWS_REGION")

	if elasticapm.DefaultTracer.Service.Framework == nil {
		executionEnv := os.Getenv("AWS_EXECUTION_ENV")
		version := strings.TrimPrefix(executionEnv, "AWS_Lambda_")
		elasticapm.DefaultTracer.Service.Framework = &model.Framework{
			Name:    "AWS Lambda",
			Version: version,
		}
	}
}

type Function struct {
	client *rpc.Client
	tracer *elasticapm.Tracer
}

func (f *Function) Ping(req *messages.PingRequest, response *messages.PingResponse) error {
	return f.client.Call("Function.Ping", req, response)
}

func (f *Function) Invoke(req *messages.InvokeRequest, response *messages.InvokeResponse) error {
	tx := f.tracer.StartTransaction(lambdacontext.FunctionName, "function")
	defer f.tracer.Flush(nonBlocking)
	defer tx.Done(-1)
	defer f.tracer.Recover(tx)
	tx.Context = &txContext

	lambdaContext.RequestID = req.RequestId
	lambdaContext.XAmznTraceID = req.XAmznTraceId
	lambdaContext.Request = formatPayload(req.Payload)
	lambdaContext.Response = ""

	err := f.client.Call("Function.Invoke", req, response)
	if err != nil {
		e := f.tracer.NewError()
		e.Transaction = tx
		e.SetException(err)
		e.Send()
		return err
	}

	if response.Payload != nil {
		lambdaContext.Response = formatPayload(response.Payload)
	}

	if response.Error != nil {
		e := f.tracer.NewError()
		e.Transaction = tx
		e.Exception = &model.Exception{
			Message: response.Error.Message,
			Type:    response.Error.Type,
		}
		frames := make([]model.StacktraceFrame, len(response.Error.StackTrace))
		for i, f := range response.Error.StackTrace {
			packagePath, functionName := stacktrace.SplitFunctionName(f.Label)
			frames[i].File = filepath.Base(f.Path)
			frames[i].Line = int(f.Line)
			frames[i].Module = packagePath
			frames[i].Function = functionName
			frames[i].LibraryFrame = stacktrace.IsLibraryPackage(packagePath)
		}
	}
	return nil
}

func formatPayload(payload []byte) string {
	if len(payload) > payloadLimit {
		payload = payload[:payloadLimit]
	}
	if !utf8.Valid(payload) {
		return ""
	}
	return string(payload)
}

func init() {
	pipeClient, pipeServer := net.Pipe()
	rpcClient := rpc.NewClient(pipeClient)
	go rpc.DefaultServer.ServeConn(pipeServer)

	origPort := os.Getenv("_LAMBDA_SERVER_PORT")
	lis, err := net.Listen("tcp", "localhost:"+origPort)
	if err != nil {
		log.Fatal(err)
	}
	srv := rpc.NewServer()
	srv.Register(&Function{
		client: rpcClient,
		tracer: elasticapm.DefaultTracer,
	})
	go srv.Accept(lis)

	// Setting _LAMBDA_SERVER_PORT causes lambda.Start
	// to listen on any free port. We don't care which;
	// we don't use it.
	os.Setenv("_LAMBDA_SERVER_PORT", "0")
}

// TODO(axw) Start() function, which wraps a given function
// such that its context is updated with the transaction.
