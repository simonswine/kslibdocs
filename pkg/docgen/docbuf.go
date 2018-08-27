package docgen

import "bytes"

type docbuf struct {
	err error
	buf bytes.Buffer
}

func (d *docbuf) Err() error {
	return d.err
}

func (d *docbuf) WriteString(s string) {
	if d.err != nil {
		return
	}

	_, d.err = d.buf.WriteString(s)
}

func (d *docbuf) Write(data []byte) (int, error) {
	if d.err != nil {
		return 0, d.err
	}

	var n int
	n, d.err = d.buf.Write(data)
	return n, d.err
}

func (d *docbuf) String() string {
	if d.err != nil {
		return ""
	}

	return d.buf.String()
}

func (d *docbuf) Bytes() []byte {
	if d.err != nil {
		return nil
	}

	return d.buf.Bytes()
}
