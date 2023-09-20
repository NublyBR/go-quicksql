package quicksql

import (
	"fmt"
	"io"
)

type spacer struct {
	wri io.Writer

	top string
	mid string
	bot string

	split, n int

	header     bool
	headerinfo []any
	buffered   []any
}

func (s *spacer) push(elem ...any) error {
	var err error

	if s.n <= 0 {
		s.n = s.split
	}

	if !s.header {
		s.header = true

		_, err = fmt.Fprintf(s.wri, s.top, s.headerinfo...)
		if err != nil {
			return err
		}
	}

	if s.buffered != nil {
		s.n--

		if s.n == 0 {
			_, err = fmt.Fprintf(s.wri, s.bot, s.buffered...)
			s.header = false
			s.n = s.split
		} else {
			_, err = fmt.Fprintf(s.wri, s.mid, s.buffered...)
		}

		if err != nil {
			return err
		}
	}

	s.buffered = elem
	return nil
}

func (s *spacer) flush() error {
	var err error

	// This will only be true if no call to push was made
	if !s.header && s.buffered == nil {
		return nil
	}

	if !s.header {
		s.header = true

		_, err = fmt.Fprintf(s.wri, s.top, s.headerinfo...)
		if err != nil {
			return err
		}
	}

	if s.buffered != nil {
		s.n--
		_, err = fmt.Fprintf(s.wri, s.bot, s.buffered...)
		s.header = false
		s.n = s.split
		return err
	}

	return nil
}
