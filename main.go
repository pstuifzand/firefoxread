package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/pierrec/lz4"
)

type Entry struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	ID    int    `json:"id"`
}
type Tab struct {
	Entries []Entry `json:"entries"`
	Index   int     `json:"index"`
}

type Window struct {
	Tabs     []Tab `json:"tabs"`
	Selected int   `json:"selected"`
}

type Session struct {
	Windows        []Window `json:"windows"`
	SelectedWindow int      `json:"selectedWindow"`
}

func main() {
	// Compress and uncompress an input string.
	f, err := os.Open("/home/peter/.mozilla/firefox/bcwmsgew.default-release/sessionstore-backups/recovery.jsonlz4")
	if err != nil {
		log.Fatal(err)
	}

	var header [12]byte
	f.Read(header[:])

	l := binary.LittleEndian.Uint32(header[8:])

	output := make([]byte, l)
	input := make([]byte, l)

	_, err = f.Read(input)
	if err != nil {
		log.Fatal(err)
	}

	lz4.UncompressBlock(input, output)

	br := bytes.NewReader(output)

	if len(os.Args) > 1 && os.Args[1] == "-json" {
		io.Copy(os.Stdout, br)
	} else {
		var session Session

		err = json.NewDecoder(br).Decode(&session)
		if err != nil {
			log.Fatal(err)
		}

		for _, window := range session.Windows {
			for tabIndex, tab := range window.Tabs {
				if tabIndex == window.Selected-1 {
					entries := tab.Entries
					spew.Dump(entries[tab.Index-1])
				}
			}
		}

		// window := session.Windows[session.SelectedWindow-1]
		// tab := window.Tabs[window.Selected-1]
		// entries := tab.Entries
		// spew.Dump(entries[tab.Index-1])
	}
}
