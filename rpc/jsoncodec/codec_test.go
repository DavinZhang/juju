// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package jsoncodec_test

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
	stdtesting "testing"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/rpc"
	"github.com/DavinZhang/juju/rpc/jsoncodec"
)

type suite struct {
	testing.LoggingSuite
}

var _ = gc.Suite(&suite{})

func TestPackage(t *stdtesting.T) {
	gc.TestingT(t)
}

type value struct {
	X string
}

func (*suite) TestRead(c *gc.C) {
	for i, test := range []struct {
		msg        string
		expectHdr  rpc.Header
		expectBody interface{}
	}{{
		msg: `{"RequestId": 1, "Type": "foo", "Id": "id", "Request": "frob", "Params": {"X": "param"}}`,
		expectHdr: rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:   "foo",
				Id:     "id",
				Action: "frob",
			},
		},
		expectBody: &value{X: "param"},
	}, {
		msg: `{"RequestId": 2, "Error": "an error", "ErrorCode": "a code"}`,
		expectHdr: rpc.Header{
			RequestId: 2,
			Error:     "an error",
			ErrorCode: "a code",
		},
		expectBody: new(map[string]interface{}),
	}, {
		msg: `{"RequestId": 3, "Response": {"X": "result"}}`,
		expectHdr: rpc.Header{
			RequestId: 3,
		},
		expectBody: &value{X: "result"},
	}, {
		msg: `{"RequestId": 4, "Type": "foo", "Version": 2, "Id": "id", "Request": "frob", "Params": {"X": "param"}}`,
		expectHdr: rpc.Header{
			RequestId: 4,
			Request: rpc.Request{
				Type:    "foo",
				Version: 2,
				Id:      "id",
				Action:  "frob",
			},
		},
		expectBody: &value{X: "param"},
	}, {
		msg: `{"request-id": 1, "type": "foo", "id": "id", "request": "frob", "params": {"X": "param"}}`,
		expectHdr: rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:   "foo",
				Id:     "id",
				Action: "frob",
			},
			Version: 1,
		},
		expectBody: &value{X: "param"},
	}, {
		msg: `{"request-id": 2, "error": "an error", "error-code": "a code"}`,
		expectHdr: rpc.Header{
			RequestId: 2,
			Error:     "an error",
			ErrorCode: "a code",
			Version:   1,
		},
		expectBody: new(map[string]interface{}),
	}, {
		msg: `{"request-id": 2, "error": "an error", "error-code": "a code", "error-info": {"foo": "bar", "baz": true}}`,
		expectHdr: rpc.Header{
			RequestId: 2,
			Error:     "an error",
			ErrorCode: "a code",
			ErrorInfo: map[string]interface{}{
				"foo": "bar",
				"baz": true,
			},
			Version: 1,
		},
		expectBody: new(map[string]interface{}),
	}, {
		msg: `{"request-id": 3, "response": {"X": "result"}}`,
		expectHdr: rpc.Header{
			RequestId: 3,
			Version:   1,
		},
		expectBody: &value{X: "result"},
	}, {
		msg: `{"request-id": 4, "type": "foo", "version": 2, "id": "id", "request": "frob", "params": {"X": "param"}}`,
		expectHdr: rpc.Header{
			RequestId: 4,
			Request: rpc.Request{
				Type:    "foo",
				Version: 2,
				Id:      "id",
				Action:  "frob",
			},
			Version: 1,
		},
		expectBody: &value{X: "param"},
	}} {
		c.Logf("test %d", i)
		codec := jsoncodec.New(&testConn{
			readMsgs: []string{test.msg},
		})
		var hdr rpc.Header
		err := codec.ReadHeader(&hdr)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(hdr, gc.DeepEquals, test.expectHdr)

		c.Assert(hdr.IsRequest(), gc.Equals, test.expectHdr.IsRequest())

		body := reflect.New(reflect.ValueOf(test.expectBody).Type().Elem()).Interface()
		err = codec.ReadBody(body, test.expectHdr.IsRequest())
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(body, gc.DeepEquals, test.expectBody)

		err = codec.ReadHeader(&hdr)
		c.Assert(err, gc.Equals, io.EOF)
	}
}

func (*suite) TestErrorAfterClose(c *gc.C) {
	conn := &testConn{
		err: errors.New("some error"),
	}
	codec := jsoncodec.New(conn)
	var hdr rpc.Header
	err := codec.ReadHeader(&hdr)
	c.Assert(err, gc.ErrorMatches, "error receiving message: some error")

	err = codec.Close()
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(conn.closed, jc.IsTrue)

	err = codec.ReadHeader(&hdr)
	c.Assert(err, gc.Equals, io.EOF)
}

func (*suite) TestWrite(c *gc.C) {
	for i, test := range []struct {
		hdr       *rpc.Header
		body      interface{}
		isRequest bool
		expect    string
	}{{
		hdr: &rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:   "foo",
				Id:     "id",
				Action: "frob",
			},
		},
		body:   &value{X: "param"},
		expect: `{"RequestId": 1, "Type": "foo","Id":"id", "Request": "frob", "Params": {"X": "param"}}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 2,
			Error:     "an error",
			ErrorCode: "a code",
		},
		expect: `{"RequestId": 2, "Error": "an error", "ErrorCode": "a code"}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 2,
			Error:     "an error",
			ErrorCode: "a code",
			ErrorInfo: map[string]interface{}{
				"ignored": "for version0",
			},
		},
		expect: `{"RequestId": 2, "Error": "an error", "ErrorCode": "a code"}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 3,
		},
		body:   &value{X: "result"},
		expect: `{"RequestId": 3, "Response": {"X": "result"}}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 4,
			Request: rpc.Request{
				Type:    "foo",
				Version: 2,
				Id:      "",
				Action:  "frob",
			},
		},
		body:   &value{X: "param"},
		expect: `{"RequestId": 4, "Type": "foo", "Version": 2, "Request": "frob", "Params": {"X": "param"}}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:   "foo",
				Id:     "id",
				Action: "frob",
			},
			Version: 1,
		},
		body:   &value{X: "param"},
		expect: `{"request-id": 1, "type": "foo","id":"id", "request": "frob", "params": {"X": "param"}}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 2,
			Error:     "an error",
			ErrorCode: "a code",
			Version:   1,
		},
		expect: `{"request-id": 2, "error": "an error", "error-code": "a code"}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 2,
			Error:     "an error",
			ErrorCode: "a code",
			ErrorInfo: map[string]interface{}{
				"foo": "bar",
				"baz": true,
			},
			Version: 1,
		},
		expect: `{"request-id": 2, "error": "an error", "error-code": "a code", "error-info": {"foo": "bar", "baz": true}}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 3,
			Version:   1,
		},
		body:   &value{X: "result"},
		expect: `{"request-id": 3, "response": {"X": "result"}}`,
	}, {
		hdr: &rpc.Header{
			RequestId: 4,
			Request: rpc.Request{
				Type:    "foo",
				Version: 2,
				Id:      "",
				Action:  "frob",
			},
			Version: 1,
		},
		body:   &value{X: "param"},
		expect: `{"request-id": 4, "type": "foo", "version": 2, "request": "frob", "params": {"X": "param"}}`,
	}} {
		c.Logf("test %d", i)
		var conn testConn
		codec := jsoncodec.New(&conn)
		err := codec.WriteMessage(test.hdr, test.body)
		c.Assert(err, jc.ErrorIsNil)
		c.Assert(conn.writeMsgs, gc.HasLen, 1)

		assertJSONEqual(c, conn.writeMsgs[0], test.expect)
	}
}

func (*suite) TestDumpRequest(c *gc.C) {
	for i, test := range []struct {
		hdr    rpc.Header
		body   interface{}
		expect string
	}{{
		hdr: rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:   "Foo",
				Id:     "id",
				Action: "Something",
			},
		},
		body:   struct{ Arg string }{Arg: "an arg"},
		expect: `{"RequestId":1,"Type":"Foo","Id":"id","Request":"Something","Params":{"Arg":"an arg"}}`,
	}, {
		hdr: rpc.Header{
			RequestId: 2,
		},
		body:   struct{ Ret string }{Ret: "return value"},
		expect: `{"RequestId":2,"Response":{"Ret":"return value"}}`,
	}, {
		hdr: rpc.Header{
			RequestId: 3,
		},
		expect: `{"RequestId":3}`,
	}, {
		hdr: rpc.Header{
			RequestId: 4,
			Error:     "an error",
			ErrorCode: "an error code",
		},
		expect: `{"RequestId":4,"Error":"an error","ErrorCode":"an error code"}`,
	}, {
		hdr: rpc.Header{
			RequestId: 5,
		},
		body:   make(chan int),
		expect: `"marshal error: json: unsupported type: chan int"`,
	}, {
		hdr: rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:    "Foo",
				Version: 2,
				Id:      "id",
				Action:  "Something",
			},
		},
		body:   struct{ Arg string }{Arg: "an arg"},
		expect: `{"RequestId":1,"Type":"Foo","Version":2,"Id":"id","Request":"Something","Params":{"Arg":"an arg"}}`,
	}, {
		hdr: rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:   "Foo",
				Id:     "id",
				Action: "Something",
			},
			Version: 1,
		},
		body:   struct{ Arg string }{Arg: "an arg"},
		expect: `{"request-id":1,"type":"Foo","id":"id","request":"Something","params":{"Arg":"an arg"}}`,
	}, {
		hdr: rpc.Header{
			RequestId: 2,
			Version:   1,
		},
		body:   struct{ Ret string }{Ret: "return value"},
		expect: `{"request-id":2,"response":{"Ret":"return value"}}`,
	}, {
		hdr: rpc.Header{
			RequestId: 3,
			Version:   1,
		},
		expect: `{"request-id":3}`,
	}, {
		hdr: rpc.Header{
			RequestId: 4,
			Error:     "an error",
			ErrorCode: "an error code",
			Version:   1,
		},
		expect: `{"request-id":4,"error":"an error","error-code":"an error code"}`,
	}, {
		hdr: rpc.Header{
			RequestId: 5,
			Version:   1,
		},
		body:   make(chan int),
		expect: `"marshal error: json: unsupported type: chan int"`,
	}, {
		hdr: rpc.Header{
			RequestId: 1,
			Request: rpc.Request{
				Type:    "Foo",
				Version: 2,
				Id:      "id",
				Action:  "Something",
			},
			Version: 1,
		},
		body:   struct{ Arg string }{Arg: "an arg"},
		expect: `{"request-id":1,"type":"Foo","version":2,"id":"id","request":"Something","params":{"Arg":"an arg"}}`,
	}} {
		c.Logf("test %d; %#v", i, test.hdr)
		data := jsoncodec.DumpRequest(&test.hdr, test.body)
		c.Check(string(data), gc.Equals, test.expect)
	}
}

// assertJSONEqual compares the json strings v0
// and v1 ignoring white space.
func assertJSONEqual(c *gc.C, v0, v1 string) {
	var m0, m1 interface{}
	err := json.Unmarshal([]byte(v0), &m0)
	c.Assert(err, jc.ErrorIsNil)
	err = json.Unmarshal([]byte(v1), &m1)
	c.Assert(err, jc.ErrorIsNil)
	data0, err := json.Marshal(m0)
	c.Assert(err, jc.ErrorIsNil)
	data1, err := json.Marshal(m1)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(data0), gc.Equals, string(data1))
}

type testConn struct {
	readMsgs  []string
	err       error
	writeMsgs []string
	closed    bool
}

func (c *testConn) Receive(msg interface{}) error {
	if len(c.readMsgs) > 0 {
		s := c.readMsgs[0]
		c.readMsgs = c.readMsgs[1:]
		return json.Unmarshal([]byte(s), msg)
	}
	if c.err != nil {
		return c.err
	}
	return io.EOF
}

func (c *testConn) Send(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	c.writeMsgs = append(c.writeMsgs, string(data))
	return nil
}

func (c *testConn) Close() error {
	c.closed = true
	return nil
}
