package logger

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"time"

	"dario.cat/mergo"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/tuantran1810/go-di-template/libs/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type LogFunc func(template string, args ...interface{})

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

type CustomFunc func(ctx context.Context, info *grpc.UnaryServerInfo) (keyAndValues []interface{})

var noSkip = func(_ context.Context, _ *grpc.UnaryServerInfo) bool {
	return false
}

var enable = func(_ context.Context, _ *grpc.UnaryServerInfo) bool {
	return true
}

var defaultCustomFunc = func(_ context.Context, _ *grpc.UnaryServerInfo) []interface{} {
	return nil
}

var DefaultOptions = Options{
	Skipper:         noSkip,
	IncludeRequest:  enable,
	IncludeResponse: enable,
	IncludeCustom:   defaultCustomFunc,
}

type Options struct {
	Skipper         func(ctx context.Context, info *grpc.UnaryServerInfo) bool
	IncludeRequest  func(ctx context.Context, info *grpc.UnaryServerInfo) bool
	IncludeResponse func(ctx context.Context, info *grpc.UnaryServerInfo) bool
	IncludeCustom   CustomFunc
}

// getLogFunc is the default implementation of gRPC return codes and interceptor log level for server side.
func getLogFunc(log Logger, code codes.Code) LogFunc {
	switch code {
	case codes.OK:
		return log.Infow
	case codes.Canceled:
		return log.Warnw
	case codes.Unknown:
		return log.Errorw
	case codes.InvalidArgument:
		return log.Warnw
	case codes.DeadlineExceeded:
		return log.Warnw
	case codes.NotFound:
		return log.Warnw
	case codes.AlreadyExists:
		return log.Warnw
	case codes.PermissionDenied:
		return log.Warnw
	case codes.Unauthenticated:
		return log.Warnw // unauthenticated requests can happen
	case codes.ResourceExhausted:
		return log.Warnw
	case codes.FailedPrecondition:
		return log.Warnw
	case codes.Aborted:
		return log.Warnw
	case codes.OutOfRange:
		return log.Warnw
	case codes.Unimplemented:
		return log.Errorw
	case codes.Internal:
		return log.Errorw
	case codes.Unavailable:
		return log.Errorw
	case codes.DataLoss:
		return log.Errorw
	default:
		return log.Errorw
	}
}

func detectAndInjectCorrelationID(ctx context.Context) (context.Context, string) {
	correlationID := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		correlationIDs, ok := md[string(utils.XCorrelationID)]
		if ok && len(correlationIDs) > 0 {
			correlationID = correlationIDs[0]
		} else { // if cannot detect from header
			correlationID = utils.GetCorrelationID(ctx)
		}
	}

	return utils.InjectCorrelationIDToContext(ctx, correlationID), correlationID
}

func UnaryServerInterceptor(log Logger, args ...Options) grpc.UnaryServerInterceptor {
	opts := DefaultOptions
	for _, arg := range args {
		if err := mergo.Merge(&opts, arg, mergo.WithOverride); err != nil {
			panic(err)
		}
	}

	runtimeMarshaler := new(runtime.JSONPb)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if opts.Skipper(ctx, info) {
			return handler(ctx, info)
		}

		ctx, correlationID := detectAndInjectCorrelationID(ctx)
		logReq := opts.IncludeRequest(ctx, info)
		logResp := opts.IncludeResponse(ctx, info)

		startTime := time.Now()
		resp, err := handler(ctx, req)

		duration := time.Since(startTime)
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)

		code := status.Code(err)
		// try to get the origin error code
		if code == codes.Unknown {
			code = status.Code(errors.Unwrap(err))
		}
		if code == codes.Unknown && errors.Is(err, context.Canceled) && ctx.Err() == context.Canceled {
			code = codes.Canceled
		}

		logFn := getLogFunc(log, code)
		if code == codes.OK && method == "Check" && service == "grpc.health.v1.Health" {
			logFn = log.Debugw
		}

		args := make([]interface{}, 0, 20)
		args = append(args,
			"code", code,
			"latency_ms", duration.Milliseconds(),
			"service", service,
			"method", method,
		)

		if extras := opts.IncludeCustom(ctx, info); len(extras) > 0 {
			args = append(args, extras...)
		}
		if logReq {
			reqBody, _ := json.Marshal(req)
			args = append(args, "request", json.RawMessage(reqBody))
		}
		if logResp {
			respBody, _ := json.Marshal(resp)
			args = append(args, "response", json.RawMessage(respBody))
		}
		args = append(args, "correlation_id", correlationID)

		if err != nil {
			args = append(args, "error", err.Error())
			st, ok := status.FromError(err)
			if ok {
				stDetails := st.Proto().Details
				if len(stDetails) > 0 {
					// details has been encoded, can not use builtin json.
					details, _ := runtimeMarshaler.Marshal(stDetails)
					args = append(args, "details", json.RawMessage(details))
				}
			}
		}
		logFn("", args...)

		return resp, err
	}
}

func StreamServerInterceptor(log Logger) grpc.StreamServerInterceptor {
	runtimeMarshaler := new(runtime.JSONPb)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		startTime := time.Now()

		err := handler(srv, stream)

		duration := time.Since(startTime)
		service := path.Dir(info.FullMethod)[1:]
		method := path.Base(info.FullMethod)

		code := status.Code(err)
		if code == codes.Unknown {
			// try to get the origin error code
			code = status.Code(errors.Unwrap(err))
		}
		logFn := getLogFunc(log, code)

		args := make([]interface{}, 0, 20)
		args = append(args,
			"code", code,
			"latency_ms", duration.Milliseconds(),
			"service", service,
			"method", method,
		)

		if err != nil {
			args = append(args, "error", err.Error())
			st, ok := status.FromError(err)
			if ok {
				stDetails := st.Proto().Details
				if len(stDetails) > 0 {
					// details has been encoded, can not use builtin json.
					details, _ := runtimeMarshaler.Marshal(stDetails)
					args = append(args, "details", json.RawMessage(details))
				}
			}
		}
		logFn("", args...)

		return err
	}
}
