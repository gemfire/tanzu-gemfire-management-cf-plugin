package impl

import (
	"github.com/gemfire/cloudcache-management-cf-plugin/domain"
	"io"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RequestHelper

// RequestHelper interface provides a way to get request related items
type RequestHelper interface {
	Exchange(url string, method string, bodyReader io.Reader, connectionData *domain.ConnectionData) (urlResponse string, err error)
}
