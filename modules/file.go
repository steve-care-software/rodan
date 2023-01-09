package modules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/fslock"
	"github.com/steve-care-software/interpreter/domain/programs/modules"
)

type file struct {
	absBasePath string
	chunkSize   uint
}

func createFile(
	absBasePath string,
	chunkSize uint,
) *file {
	out := file{
		absBasePath: absBasePath,
		chunkSize:   chunkSize,
	}

	return &out
}

// Execute executes the application
func (app *file) Execute() map[uint]modules.ExecuteFn {
	fileOpen := app.fileOpen()
	fileClose := app.fileClose()
	fileLock := app.fileLock()
	fileUnLock := app.fileUnLock()
	fileInfo := app.fileInfo()
	fileRead := app.fileRead()
	fileWrite := app.fileWrite()
	return map[uint]modules.ExecuteFn{
		ModuleFileOpen:   fileOpen,
		ModuleFileClose:  fileClose,
		ModuleFileLock:   fileLock,
		ModuleFileUnLock: fileUnLock,
		ModuleFileInfo:   fileInfo,
		ModuleFileRead:   fileRead,
		ModuleFileWrite:  fileWrite,
	}
}

func (app *file) formPath(relativePath string, inputIndex uint) (string, error) {
	path := filepath.Join(app.absBasePath, relativePath)
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(absPath, app.absBasePath) {
		return path, nil
	}

	str := fmt.Sprintf("the input at index (%d) was expected to contain a relative path (%s) that was expected to not seek before the base directory (%s)", inputIndex, relativePath, app.absBasePath)
	return "", errors.New(str)
}

func (app *file) fileOpen() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if relativePath, ok := input[0].([]byte); ok {
			path, err := app.formPath(strings.TrimSpace(string(relativePath)), 0)
			if err != nil {
				return nil, err
			}

			return os.Open(path)
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a string", 0)
		return nil, errors.New(str)
	}
}

func (app *file) fileClose() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if pConn, ok := input[0].(*os.File); ok {
			err := pConn.Close()
			if err != nil {
				return nil, err
			}

			return nil, nil
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a file connection", 0)
		return nil, errors.New(str)
	}
}

func (app *file) fileLock() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if relativePath, ok := input[0].(string); ok {
			path, err := app.formPath(relativePath, 0)
			if err != nil {
				return nil, err
			}

			pLock := fslock.New(path)
			err = pLock.TryLock()
			if err != nil {
				return nil, err
			}

			return pLock, nil
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a string", 0)
		return nil, errors.New(str)
	}
}

func (app *file) fileUnLock() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if pLock, ok := input[0].(*fslock.Lock); ok {
			err := pLock.Unlock()
			if err != nil {
				return nil, err
			}

			return nil, nil
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a string", 0)
		return nil, errors.New(str)
	}
}

func (app *file) fileInfo() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if pConn, ok := input[0].(*os.File); ok {
			return pConn.Stat()
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a file connection", 0)
		return nil, errors.New(str)
	}
}

func (app *file) fileRead() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		index := uint(0)
		if idx, ok := input[2].(uint); ok {
			index = idx
		}

		sizeInBytes := int64(-1)
		if amount, ok := input[1].(uint); ok {
			sizeInBytes = int64(amount)
		}

		if pConn, ok := input[0].(*os.File); ok {
			if sizeInBytes == -1 {
				pInfo, err := pConn.Stat()
				if err != nil {
					return nil, err
				}

				sizeInBytes = pInfo.Size()
			}

			data := make([]byte, sizeInBytes)
			readAmount, err := pConn.ReadAt(data, int64(index))
			if err != nil {
				return nil, err
			}

			if int64(readAmount) != sizeInBytes {
				str := fmt.Sprintf("%d bytes were expected to be read, %d actually read", sizeInBytes, readAmount)
				return nil, errors.New(str)
			}

			return data, nil
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a file connection", 0)
		return nil, errors.New(str)
	}
}

func (app *file) fileWrite() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		index := uint(0)
		if idx, ok := input[2].(uint); ok {
			index = idx
		}

		if pConn, ok := input[0].(*os.File); ok {
			if data, ok := input[1].([]byte); ok {
				return pConn.WriteAt(data, int64(index))
			}

			str := fmt.Sprintf("the input at index (%d) was expected to contain []byte", 1)
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a file connection", 0)
		return nil, errors.New(str)
	}
}
