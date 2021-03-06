/*
MIT License

Copyright (c) 2018 KIM KeepInMind Gmbh/srl

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package socks5_test

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/booster-proj/proxy/socks5"
)

type conn struct {
	bytes.Buffer
}

// protocol stubs
func (c *conn) Close() error                       { return nil }
func (c *conn) LocalAddr() net.Addr                { return nil }
func (c *conn) RemoteAddr() net.Addr               { return nil }
func (c *conn) SetDeadline(t time.Time) error      { return nil }
func (c *conn) SetReadDeadline(t time.Time) error  { return nil }
func (c *conn) SetWriteDeadline(t time.Time) error { return nil }

func TestReadAddress(t *testing.T) {
	conn := new(conn)

	var tests = []struct {
		in  []byte
		out string
		err bool
	}{
		{in: []byte{1, 93, 184, 216, 34, 1, 187},
			out: "93.184.216.34:443",
			err: false}, // ipv4

		{in: []byte{4, 42, 3, 176, 192, 0, 3, 0, 208, 0, 0, 0, 0, 72, 136, 160, 1, 1, 187},
			out: "[2a03:b0c0:3:d0::4888:a001]:443",
			err: false}, // ipv6

		{in: []byte{3, 21, 111, 117, 116, 108, 111, 111, 107, 46, 111, 102, 102, 105, 99, 101, 51, 54, 53, 46, 99, 111, 109, 1, 187},
			out: "outlook.office365.com:443",
			err: false}, // FQDN

		{in: []byte{0, 93, 184, 216, 34, 1, 187},
			out: "",
			err: true}, // wrong address type

		{in: []byte{5, 93, 184, 216, 34, 1, 187},
			out: "",
			err: true}, // wrong address type
	}

	for _, test := range tests {
		if _, err := conn.Write(test.in); err != nil {
			t.Fatal(err)
		}

		s, err := socks5.ReadAddress(conn)
		if err != nil {
			// only fail if not expecting an error
			if !test.err {
				t.Fatal(err)
			} else {
				return
			}
		}

		t.Log("Address Read: " + s)

		if s != test.out {
			t.Fatal(err)
		}
	}
}
