package errs

import "fmt"

type SessionNotFoundError struct {
	Key string
}

func (e *SessionNotFoundError) Error() string {
	return fmt.Sprintf("Сессии %s не существует", e.Key)
}
