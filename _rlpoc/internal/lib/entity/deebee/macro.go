package deebee

import (
	"encoding/binary"
	"fmt"
	"os"
)

type macroStore struct {
	pages       map[macroId]*macroPage
	currentPage *macroPage
	micro       microId
}

func (ms *macroStore) More(minimum uint64) (err error) {
	panic("not implemented")
}

func (ms *macroStore) Checkpoint() error {
	panic("not implemented")
}

func (ms *macroStore) WriteFrame(frameClass int, b []byte) (id ID, err error) {
	if frameClass < 0 {
		return ID{}, fmt.Errorf("invalid frameClass: %d", frameClass)
	}
	return ms.writeFrame(frameClass, b)
}

func (ms *macroStore) writeFrame(frameClass int, b []byte) (id ID, err error) {
	id = ID{
		macro: ms.currentPage.macro,
		micro: ms.micro,
	}
	uvarBuf := make([]byte, binary.MaxVarintLen64*2)
	uvarSize := binary.PutVarint(uvarBuf, int64(len(b)))
	uvarSize += binary.PutVarint(uvarBuf[:uvarSize], int64(frameClass))
	if _, err = ms.currentPage.Write(uvarBuf[:uvarSize]); err != nil {
		return ID{}, fmt.Errorf("writing frame envelope to current page: %w", err)
	}

	if _, err = ms.currentPage.Write(b); err != nil {
		return ID{}, fmt.Errorf("writing data to current page: %w", err)
	}

	ms.micro += microId(uvarSize + len(b))

	return id, nil
}

type macroPage struct {
	file  *os.File
	macro macroId
}

// Write writes the given bytes the the page on disk.
// If it returns an error, the page will have to be reopened.
func (page *macroPage) Write(b []byte) (n int, err error) {
	n, err = page.file.Write(b)
	if err != nil {
		firstError := fmt.Errorf("writing to page file: %w", err)
		//if n > 0 {
		//	_, err := page.file.Seek(-int64(n), 1)
		//	if err != nil {
		//		return 0, fmt.Errorf("seeking back error (%v) AFTER ERROR: %w", err, firstError)
		//	}
		//	if page.currentPosition > math.MaxInt64 {
		//		return 0, fmt.Errorf("current file position (%d) is greater than allowable Truncation range (%d) AFTER ERROR: %w", page.currentPosition, math.MaxInt64, firstError)
		//	}
		//	if err := page.file.Truncate(int64(page.currentPosition)); err != nil {
		//		return 0, fmt.Errorf("truncating file (%v) AFTER ERROR: %w", err, firstError)
		//	}
		//}
		//return 0, newRetryAllowedError(firstError)
		return n, firstError
	}
	if n < 0 {
		return 0, fmt.Errorf("file write returned less than 0: %d", n)
	}
	return n, nil
}
