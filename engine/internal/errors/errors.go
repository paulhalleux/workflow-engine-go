package errors

type UnsupportedProtocolError struct{}

func (e *UnsupportedProtocolError) Error() string {
	return "unsupported agent protocol"
}

var ErrUnsupportedProtocol = &UnsupportedProtocolError{}
