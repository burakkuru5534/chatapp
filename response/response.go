package response

import (
	"encoding/json"
	"example.com/m/conf"
	"fmt"
	"github.com/rs/zerolog"
	"log"
	"log/syslog"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"net/http"
	"runtime"
)

// swagger:model Response
type Response struct {
	ID      string
	Success bool
	Message string
	Data    interface{}
}

func NewResponse(success bool, message string, data interface{}) (*Response, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	m := &Response{
		ID:      u.String(),
		Success: success,
		Message: message,
		Data:    data,
	}

	return m, nil
}

func (m *Response) SendWithStatus(w http.ResponseWriter, statusCode int) error {
	encjson, err := json.Marshal(m)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(encjson)
	return err
}

func (m *Response) Send(w http.ResponseWriter) error {
	return m.SendWithStatus(w, http.StatusOK)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	ar, err := NewResponse(true, "", data)
	if err != nil {
		SendServerError(w, err, http.StatusInternalServerError)
		return
	}
	err = ar.Send(w)
	if err != nil {
		SendServerError(w, err, http.StatusInternalServerError)
	}
}

func SendServerError(w http.ResponseWriter, err error, statusCode int) {
	l := GetLogger(conf.LogLevel)

	_, fn, line, _ := runtime.Caller(1)

	logMessage := fmt.Sprintf("%v -- %s:%d", err, fn, line)
	l.Logger.Error().Msg(logMessage)

	http.Error(w, logMessage, statusCode)
}
type Logger struct {
	zerolog.Logger
}
var loggerinstance *Logger
type logOutputType int
const (
	Console = iota
	Syslog
	GrayLog
)

func New(text string) error {
	return &errorString{text}
}
type errorString struct {
	s string
}
func (e *errorString) Error() string {
	return e.s
}
func NewLogger(outputType logOutputType, level zerolog.Level) (*Logger, error) {
	var zlogger zerolog.Logger

	switch outputType {
	case Console:
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("Msg: %s |", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("| %s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return fmt.Sprintf("%s |", i)
		}

		zlogger = zerolog.New(output).With().Timestamp().Logger()
	case Syslog:
		syslogger, err := syslog.New(syslog.LOG_DEBUG, conf.Product)
		if err != nil {
			log.Fatal(err)
		}
		zlogger = zerolog.New(syslogger).With().Timestamp().Logger()
	case GrayLog:
		//todo: send to graylog
		return nil, New("graylog logging not implemented")
	default:
		return nil, New("unknown logging output type")
	}

	loggerinstance = &Logger{zlogger.Level(level)}
	return loggerinstance, nil
}

func GetLogger(level zerolog.Level) *Logger {
	if loggerinstance == nil {
		_, _ = NewLogger(Console, level)
	}

	return loggerinstance
}