package feed

import (
	"bytes"
	"errors"
	"os"
	"strconv"

	atom "github.com/denisbrodbeck/atomfeed"
)

func WriteFeed(feed *atom.Feed, outpath string, perm string) error {
	if feed == nil {
		return errors.New("Atom Feed struct is nil")
	}
	perm32, err := strconv.ParseUint(perm, 8, 32)
	if err != nil {
		return err
	}
	if perm32 > 777 || perm32 < 0 {
		return errors.New("Perm range error [0000,0777]")
	}

	out := &bytes.Buffer{}
	if err := feed.Encode(out); err != nil {
		return err
	}
	err = os.WriteFile(outpath, out.Bytes(), os.FileMode(perm32))
	if err != nil {
		return err
	}

	return nil
}
