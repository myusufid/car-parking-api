// NOSONAR
package request_logger

import (
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LogRequestParams encapsulates parameters for logging the request
type LogRequestParams struct {
	C          *fiber.Ctx
	StartTime  time.Time
	Start      time.Time
	Stop       time.Time
	ErrPadding int
	PID        string
	ChainErr   error
}

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Set PID once
	pid := strconv.Itoa(os.Getpid())

	// Set variables
	var (
		once           sync.Once
		errHandler     fiber.ErrorHandler
		errPadding     = 15
		latencyEnabled = contains("latency", cfg.Fields)
		skipURIs       = make(map[string]struct{})
	)

	for _, uri := range cfg.SkipURIs {
		skipURIs[uri] = struct{}{}
	}

	// Return new handler
	return func(c *fiber.Ctx) (err error) {
		startTime := time.Now()

		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		if _, ok := skipURIs[c.Path()]; ok {
			return c.Next()
		}

		once.Do(func() {
			errPadding = calculateErrPadding(c)
			errHandler = c.App().Config().ErrorHandler
		})

		var start, stop time.Time
		if latencyEnabled {
			start = time.Now()
		}

		chainErr := c.Next()

		if chainErr != nil {
			handleErrorHandler(c, errHandler, chainErr)
		}

		if latencyEnabled {
			stop = time.Now()
		}

		logRequest(LogRequestParams{
			C:          c,
			StartTime:  startTime,
			Start:      start,
			Stop:       stop,
			ErrPadding: errPadding,
			PID:        pid,
			ChainErr:   chainErr,
		}, cfg)

		return nil
	}
}

func calculateErrPadding(c *fiber.Ctx) int {
	stack := c.App().Stack()
	maxLength := 15 // Default value
	for m := range stack {
		for r := range stack[m] {
			if len(stack[m][r].Path) > maxLength {
				maxLength = len(stack[m][r].Path)
			}
		}
	}
	return maxLength
}

func handleErrorHandler(c *fiber.Ctx, errHandler fiber.ErrorHandler, chainErr error) {
	if err := errHandler(c, chainErr); err != nil {
		_ = c.SendStatus(fiber.StatusInternalServerError)
	}
}

func logRequest(params LogRequestParams, cfg Config) {
	s := params.C.Response().StatusCode()
	levelIndex := calculateLevelIndex(s)
	messageIndex := calculateMessageIndex(s, cfg.Messages)

	ce := cfg.Logger.Check(cfg.Levels[levelIndex], cfg.Messages[messageIndex])
	if ce == nil {
		return
	}

	fields := buildFields(params, cfg)

	logRAMUsage(fields)
	logRequestDuration(fields, params.StartTime)
	logRequestID(fields, params.C)

	ce.Write(fields...)
}

func calculateLevelIndex(statusCode int) int {
	switch {
	case statusCode >= 500:
		return 0
	case statusCode >= 400:
		return 1
	default:
		return 2
	}
}

func calculateMessageIndex(statusCode int, messages []string) int {
	index := calculateLevelIndex(statusCode)
	if index >= len(messages) {
		return len(messages) - 1
	}
	return index
}

func buildFields(params LogRequestParams, cfg Config) []zap.Field {
	fields := make([]zap.Field, 0, len(cfg.Fields)+1)
	fields = append(fields, zap.Error(params.ChainErr))

	if cfg.FieldsFunc != nil {
		fields = append(fields, cfg.FieldsFunc(params.C)...)
	}

	for _, field := range cfg.Fields {
		switch field {
		case "referer":
			fields = append(fields, zap.String("referer", params.C.Get(fiber.HeaderReferer)))
		case "protocol":
			fields = append(fields, zap.String("protocol", params.C.Protocol()))
		case "pid":
			fields = append(fields, zap.String("pid", params.PID))
		case "port":
			fields = append(fields, zap.String("port", params.C.Port()))
		case "ip":
			fields = append(fields, zap.String("ip", params.C.IP()))
		case "ips":
			fields = append(fields, zap.String("ips", params.C.Get(fiber.HeaderXForwardedFor)))
		case "host":
			fields = append(fields, zap.String("host", params.C.Hostname()))
		case "path":
			fields = append(fields, zap.String("path", params.C.Path()))
		case "url":
			fields = append(fields, zap.String("url", params.C.OriginalURL()))
		case "ua":
			fields = append(fields, zap.String("ua", params.C.Get(fiber.HeaderUserAgent)))
		case "latency":
			fields = append(fields, zap.String("latency", params.Stop.Sub(params.Start).String()))
		case "status":
			fields = append(fields, zap.Int("status", params.C.Response().StatusCode()))
		case "resBody":
			fields = respBody(fields, cfg, params)
		case "queryParams":
			fields = append(fields, zap.String("queryParams", params.C.Request().URI().QueryArgs().String()))
		case "body":
			if cfg.SkipBody == nil || !cfg.SkipBody(params.C) {
				fields = append(fields, zap.ByteString("body", params.C.Body()))
			}
		case "bytesReceived":
			fields = append(fields, zap.Int("bytesReceived", len(params.C.Request().Body())))
		case "bytesSent":
			fields = append(fields, zap.Int("bytesSent", len(params.C.Response().Body())))
		case "route":
			fields = append(fields, zap.String("route", params.C.Route().Path))
		case "method":
			fields = append(fields, zap.String("method", params.C.Method()))
		case "requestId":
			fields = append(fields, zap.String("requestId", params.C.GetRespHeader(fiber.HeaderXRequestID)))
		case "error":
			if params.ChainErr != nil {
				fields = append(fields, zap.String("error", params.ChainErr.Error()))
			}
		case "reqHeaders":
			params.C.Request().Header.VisitAll(func(k, v []byte) {
				fields = append(fields, zap.ByteString(string(k), v))
			})
		}
	}

	fields = logRAMUsage(fields)
	fields = logRequestDuration(fields, params.StartTime)
	fields = logRequestID(fields, params.C)

	// Other fields logic remains the same as before

	return fields
}

func respBody(fields []zap.Field, cfg Config, params LogRequestParams) []zap.Field {
	if cfg.SkipResBody == nil || !cfg.SkipResBody(params.C) {
		if cfg.GetResBody == nil {
			fields = append(fields, zap.ByteString("resBody", params.C.Response().Body()))
		} else {
			fields = append(fields, zap.ByteString("resBody", cfg.GetResBody(params.C)))
		}
	}
	return fields
}

func logRAMUsage(fields []zap.Field) []zap.Field {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return append(fields, zap.Any("ram_usage", m.Alloc/1024/1024))
}

func logRequestDuration(fields []zap.Field, startTime time.Time) []zap.Field {
	duration := time.Since(startTime)
	return append(fields, zap.Any("duration", duration))
}

func logRequestID(fields []zap.Field, c *fiber.Ctx) []zap.Field {
	requestID := c.GetRespHeader(fiber.HeaderXRequestID)
	return append(fields, zap.Any("request_id", requestID))
}

func contains(needle string, slice []string) bool {
	for _, e := range slice {
		if e == needle {
			return true
		}
	}

	return false
}
