package version

import "fmt"

const (
	unhandledErrorCode = iota + 2000
	fileAlreadyExistsErrorCode
	fileNotFoundErrorCode
	fileFormatErrorCode
	projectDirectoryNotFoundErrorCode
	inputValueErrorCode
	writeOperationFailedErrorCode
)

type UnhandledError string

func (p UnhandledError) Error() string {
	return fmt.Sprintf("error code: %d - unhandled exception: %q", unhandledErrorCode, string(p))
}

type FileAlreadyExistsError string

func (p FileAlreadyExistsError) Error() string {
	return fmt.Sprintf("error code: %d - Version File %q Already Exists", fileAlreadyExistsErrorCode, string(p))
}

type FileNotFoundError string

func (p FileNotFoundError) Error() string {
	return fmt.Sprintf("error code: %d - Version File %q Not Found", fileNotFoundErrorCode, string(p))
}

type FileFormatError string

func (p FileFormatError) Error() string {
	return fmt.Sprintf("error code: %d - Invalid Version Format in File %q", fileFormatErrorCode, string(p))
}

type ProjectDirectoryNotFoundError string

func (p ProjectDirectoryNotFoundError) Error() string {
	return fmt.Sprintf("error code: %d - Project Directory not Found %q", projectDirectoryNotFoundErrorCode, string(p))
}

type InputValueError string

func (p InputValueError) Error() string {
	return fmt.Sprintf("error code: %d - Invalid Version Format %q", inputValueErrorCode, string(p))
}

type WriteOperationFailedError struct {
	file  string
	error error
}

func (p WriteOperationFailedError) Error() string {
	return fmt.Sprintf("error code: %d - Failed to Write to File %q: %v", writeOperationFailedErrorCode, p.file, p.error)
}
