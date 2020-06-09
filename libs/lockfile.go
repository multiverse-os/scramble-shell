package lockfile

import (
	"os"
	paths "path"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// Lockfile represents an acquired lockfile.
type Lockfile struct {
	file os.File
}

// Acquire creates the given file path if it doesn't exist and
// obtains an exclusive lock on it. An error is returned if the lock
// has been obtained by another process.
func Acquire(path string) (*Lockfile, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "failed opening lock path")
	}

	ft := &unix.Flock_t{
		Pid:  int32(os.Getpid()),
		Type: unix.F_WRLCK,
	}

	if err = unix.FcntlFlock(file.Fd(), unix.F_SETLK, ft); err != nil {
		return nil, errors.Wrap(err, "failed obtaining lock")
	}

	lf := Lockfile{*file}

	return &lf, nil
}

// CreateAndAcquire creates any non-existing directories needed to
// create the lock file, then acquires a lock on it
func CreateAndAcquire(path string, newDirMode os.FileMode) (*Lockfile, error) {
	if err := os.MkdirAll(paths.Dir(path), newDirMode); err != nil {
		return nil, err
	}

	return Acquire(path)
}

// Release releases the lock on the file and removes the file.
func (lf Lockfile) Release() error {
	ft := &unix.Flock_t{
		Pid:  int32(os.Getpid()),
		Type: unix.F_UNLCK,
	}

	if err := unix.FcntlFlock(lf.file.Fd(), unix.F_SETLK, ft); err != nil {
		return errors.Wrap(err, "failed releasing lock")
	}

	if err := lf.file.Close(); err != nil {
		return errors.Wrap(err, "failed closing lock file descriptor")
	}

	if err := os.Remove(lf.file.Name()); err != nil {
		return errors.Wrap(err, "failed removing lock file")
	}

	return nil
}
