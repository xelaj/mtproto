package serialize

import (
	"bytes"
	"compress/gzip"

	"github.com/xelaj/go-dry"
)

func decompressData(data []byte) ([]byte, error) {
	decompressed := make([]byte, 0, 4096)

	var buf bytes.Buffer
	_, _ = buf.Write(data)
	gz, err := gzip.NewReader(&buf)
	dry.PanicIfErr(err)
	b := make([]byte, 4096)
	for {
		n, _ := gz.Read(b)

		decompressed = append(decompressed, b[0:n]...)
		if n <= 0 {
			break
		}
	}

	return decompressed, nil
	//? это то что я пытался сделать
	// data := d.PopMessage()
	// gz, err := gzip.NewReader(bytes.NewBuffer(data))
	// dry.PanicIfErr(err)

	// decompressed, err := ioutil.ReadAll(gz)
	// dry.PanicIfErr(err)

	// return decompressed
}
