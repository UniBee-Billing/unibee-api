package exception

import "fmt"

type AlipayLibraryError struct {
	Message string
}

func (p *AlipayLibraryError) Error() string {
	return fmt.Sprintf("AlipaySDKErrore=%s", p.Message)
}
