package telnet

import (
	"io"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Given a string reader of a 15-character-string", t, func() {
		rawdata := "Heeeelllllooooo"
		sr := strings.NewReader(rawdata)
		r := newReader(sr)
		Convey("With a small buffer of size 5", func() {
			l := 5
			buf := make([]byte, l)
			Convey("only 5 characters can be read for once", func() {
				n, err := r.read(buf)
				So(n, ShouldEqual, l)
				So(err, ShouldBeNil)
				So(string(buf[:n]), ShouldEqual, "Heeee")
			})
			Convey("finish reading to eof with 5 read calls", func() {
				count := 0
				expects := [5]string{
					"Heeee",
					"lllll",
					"ooooo",
				}
				for {
					count++
					n, err := r.read(buf)
					if err == io.EOF {
						So(count, ShouldEqual, 5)
						So(n, ShouldEqual, 0)
						break
					}
					if err == ErrEOS {
						So(n, ShouldEqual, 0)
					} else {
						So(n, ShouldEqual, 5)
						s := string(buf[:n])
						So(s, ShouldEqual, expects[count-1])
					}
				}
				So(count, ShouldEqual, 5)
			})
		})
		Convey("With a buffer of exact 15 bytes long", func() {
			l := 15
			buf := make([]byte, l)
			Convey("Reads the entire string on 1st read, EOP on 2nd read, and EOF on 3rd read", func() {
				n, err := r.read(buf)
				So(n, ShouldEqual, l)
				So(err, ShouldBeNil)
				So(string(buf), ShouldEqual, rawdata)
				n, err = r.read(buf)
				So(n, ShouldEqual, 0)
				So(err, ShouldEqual, ErrEOS)
				n, err = r.read(buf)
				So(n, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})
		})
		Convey("With a buffer larger than 15", func() {
			l := 20
			buf := make([]byte, l)
			Convey("Reads the entire string with ErrEOP, then returns empty buffer with EOF on next read", func() {
				n, err := r.read(buf)
				So(n, ShouldEqual, 15)
				So(err, ShouldEqual, ErrEOS)
				So(string(buf[:n]), ShouldEqual, rawdata)
				n, err = r.read(buf)
				So(n, ShouldEqual, 0)
				So(err, ShouldEqual, io.EOF)
			})
		})
	})
}
