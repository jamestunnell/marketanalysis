package env

import (
	"fmt"
	"os"
	"strconv"
)

type ErrInvalid struct {
	Name, Val string
	Err       error
}

const (
	NameDBConn = "BACKEND_DBCONN"
	// NameDBUser = "BACKEND_DBUSER"
	// NameDBPass = "BACKEND_DBPASS"
	NameDebug = "BACKEND_DEBUG"
	NamePort  = "BACKEND_PORT"
)

func LoadValues() (*Values, error) {
	vals := &Values{}

	if val := os.Getenv(NameDBConn); val != "" {
		vals.DBConn = val
	}

	// if val := os.Getenv(NameDBUser); val != "" {
	// 	vals.DBUser = val
	// }

	// if val := os.Getenv(NameDBPass); val != "" {
	// 	vals.DBPass = val
	// }

	if val, err := LoadOptionalBool(NameDebug, false); err != nil {
		return nil, err
	} else {
		vals.Debug = val
	}

	if val, err := LoadOptionalInt(NamePort, 0); err != nil {
		return nil, err
	} else {
		vals.Port = val
	}

	return vals, nil
}

func LoadOptionalBool(name string, defaultVal bool) (bool, error) {
	str := os.Getenv(name)
	if str == "" {
		return defaultVal, nil
	}

	val, err := strconv.ParseBool(str)
	if err != nil {
		return false, NewErrInvalid(name, str, err)
	}

	return val, nil
}

func LoadOptionalInt(name string, defaultVal int) (int, error) {
	str := os.Getenv(name)
	if str == "" {
		return defaultVal, nil
	}

	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, NewErrInvalid(name, str, err)
	}

	return val, nil
}

func NewErrInvalid(name, val string, err error) *ErrInvalid {
	return &ErrInvalid{Name: name, Val: val, Err: err}
}

func (err *ErrInvalid) Error() string {
	return fmt.Sprintf("env var %s has invalid value %s: %v", err.Name, err.Val, err.Err)
}
