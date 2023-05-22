package errors

import (
	"errors"
	"fmt"
	"net/http"
)

func IsWrongStatusCode(statusCode int) bool {
	return statusCode >= 300 || statusCode < 200
}

func StatusCodeErr(response *http.Response) error {
	return errors.New(fmt.Sprintf("response returned wrong status code: %d, from url path: %s", response.StatusCode, response.Request.URL.Path))
}
