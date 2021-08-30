package cmn

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FormValToInt64 Tries to convert a form val to Int64, if str = "", then returns nil. Used mainly converting html form values
// to BigInt column values
func FormValToInt64(aformVal string) *int64 {
	if aformVal == "" {
		return nil
	}

	res, err := strconv.ParseInt(aformVal, 10, 64)

	if err != nil {
		return nil
	}

	return &res
}

// Int64ToStr is shorthand for strconv.FormatInt with base 10
func Int64ToStr(aval int64) string {
	res := strconv.FormatInt(aval, 10)
	return res
}

// StrToInt64 is shorthand for strconv.ParseInt with base 10, bitSize 64, returns 0 if parsing error occurs.
func StrToInt64(aval string) int64 {
	aval = strings.Trim(strings.TrimSpace(aval), "\n")
	i, err := strconv.ParseInt(aval, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// StrToInt is shorhant for strconv.Atoi, returns 0 if parsing error occurs.
func StrToInt(aval string) int {
	i, err := strconv.Atoi(aval)
	if err != nil {
		return 0
	}
	return i
}

// FloatToStr, formats float number for text representation.
// Todo: add formatting options as "#,##0.00"
func FloatToStr(aval float64) string {
	return fmt.Sprintf("%f", aval)
}

// StrToFloat is shorhand for strconv.ParseFşoat with bitSize 64, returns 0 if parsing error occurs.
func StrToFloat(aval string) float64 {
	i, err := strconv.ParseFloat(aval, 64)
	if err != nil {
		return 0
	}
	return i
}

// todo: add other possible formats for date conversion
func ParseTime(val string) (time.Time, error) {
	var err error

	if res, err := time.Parse("15:04", val); err == nil {
		return res, nil
	}

	if res, err := time.Parse("02.01.2006", val); err == nil {
		return res, nil
	}

	if res, err := time.Parse("01-02-2006", val); err == nil {
		return res, nil
	}

	if res, err := time.Parse("02.01.2006 15:04", val); err == nil {
		return res, nil
	}

	return time.Time{}, err
}


func StrToTime(aval string) (time.Time, error) {
	dt := time.Time{}

	dt, err := time.Parse("15:04", aval)
	if err != nil {
		return dt, err
	}

	return dt, nil
}


func JoinInt64Array(lns []int64, sep string) string {
	lnsStr := make([]string, len(lns))
	for ndx, ln := range lns {
		lnsStr[ndx] = Int64ToStr(ln)
	}
	return strings.Join(lnsStr, sep)
}

func PtrToString(str *string) string {
	if str != nil {
		return *str
	} else {
		return ""
	}
}

func InterfaceArrayToStringArray(t []interface{}) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = fmt.Sprint(v)
	}

	return s
}

func BodyToStringReq(r *http.Request) string {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body)
}

func BodyToJsonReq(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	defer r.Body.Close()

	return nil
}

func BodyToJsonAndStringReq(r *http.Request, data interface{}) (string, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	return string(body), nil
}

func BodyToJsonResp(r *http.Response, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	return nil
}

func BodyToStringResp(r *http.Response) string {
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	return string(body)
}

func RequestIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return IPAddress
}
