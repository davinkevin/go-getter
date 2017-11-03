package getter

import (
	"net/url"
	"os"
)

// FileGetter is a Getter implementation that will download a module from
// a file scheme.
type FileGetter struct {
	// Copy, if set to true, will copy data instead of using a symlink
	Copy bool

	// Used for calculating percent progress
	totalSize       int64
	PercentComplete int
	Done            chan int64
}

func (g *FileGetter) ClientMode(u *url.URL) (ClientMode, error) {
	path := u.Path
	if u.RawPath != "" {
		path = u.RawPath
	}

	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	// Check if the source is a directory.
	if fi.IsDir() {
		return ClientModeDir, nil
	}

	return ClientModeFile, nil
}

func (g *HttpGetter) CalcDownloadPercent(dst) {
	// stat file every n seconds to figure out the download progress
	var stop bool = false
	dstfile, err := os.Open(dst)
	defer dstfile.Close()

	if err != nil {
		log.Printf("couldn't open file for reading: %s", err)
		return
	}
	for {
		select {
		case <-g.Done:
			stop = true
		default:
			fi, err := dstfile.Stat()
			if err != nil {
				fmt.Printf("Error stating file: %s", err)
				return
			}
			size := fi.Size()

			// catch edge case that would break our percentage calc
			if size == 0 {
				size = 1
			}
			g.PercentComplete = int(float64(size) / float64(g.totalSize) * 100)
		}

		if stop {
			break
		}
		// repeat check once per second
		time.Sleep(time.Second)
	}
}
