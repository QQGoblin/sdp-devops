package shell

import "bytes"

type OutPut struct {
	Title  string
	StdOut *bytes.Buffer
	StdErr *bytes.Buffer
}
