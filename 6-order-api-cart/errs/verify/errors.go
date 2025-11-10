package errs

type CodeDoesNotMatchError struct{}

func (e *CodeDoesNotMatchError) Error() string {
	return "Неправильный код"
}
