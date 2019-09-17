package libhssh

import "os"

type NoBellStdOut struct{}

func (o *NoBellStdOut) Write(b []byte) (int, error) {
	if len(b) == 1 && b[0] == 7 {
		return 0, nil
	}

	return os.Stderr.Write(b)
}

func (s *NoBellStdOut) Close() error {
	return os.Stderr.Close()
}
