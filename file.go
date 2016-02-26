package analytics

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

func (ts *Series) SavePlot(path string, name string) {
	seriestxt := make([]string, ts.Len)
	for i := range seriestxt {
		seriestxt[i] = fmt.Sprint(int64(ts.x[i]), ts.y[i])
	}
	writeLines(seriestxt, path+"/"+name)
}

// writeLines writes the lines to the given file.
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

func (ts *Series) Save(name string) {
	x := new(bytes.Buffer)
	enc := gob.NewEncoder(x)
	enc.Encode(ts.x)
	y := new(bytes.Buffer)
	enc2 := gob.NewEncoder(y)
	enc2.Encode(ts.y)
	ioutil.WriteFile(name+".x", x.Bytes(), 0600)
	ioutil.WriteFile(name+".y", y.Bytes(), 0600)
}
