package errors

// This package can be used to define custom error types or error handling utilities
// specific to the Data Explorer application, if needed.

// Example custom error type:
/*
type DataSourceError struct {
	SourceType string
	Message    string
	Err        error // Underlying error
}

func (e *DataSourceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("datasource [%s]: %s: %v", e.SourceType, e.Message, e.Err)
	}
	return fmt.Sprintf("datasource [%s]: %s", e.SourceType, e.Message)
}

func (e *DataSourceError) Unwrap() error {
	return e.Err
}

func NewDataSourceError(sourceType, message string, err error) error {
	return &DataSourceError{
		SourceType: sourceType,
		Message:    message,
		Err:        err,
	}
}
*/

// You can also define sentinel errors here:
/*
import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInvalidInput  = errors.New("invalid input provided")
	// ... other common errors
)
*/
