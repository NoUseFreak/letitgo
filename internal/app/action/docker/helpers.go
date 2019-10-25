package docker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type image struct {
	Repo string
	Tag  string
}

func parseImageName(name string) image {
	parts := strings.Split(name, ":")
	resp := image{Repo: parts[0]}

	if len(parts) > 1 {
		resp.Tag = parts[1]
	}

	return resp
}

func newProgressWriter() (io.ReadCloser, io.WriteCloser) {
	r, w, _ := os.Pipe()
	go func(r io.Reader) {
		fmt.Print("  progress: ")
		scanner := bufio.NewScanner(r)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			fmt.Print(".")
			// fmt.Println(scanner.Text())
			if scanner.Err() != nil {
				return
			}
		}
		fmt.Print("\n")
	}(r)

	return r, w
}
