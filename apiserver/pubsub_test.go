// Copyright 2016 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	jujuhttp "github.com/juju/http/v2"
	"github.com/juju/loggo"
	"github.com/juju/names/v4"
	"github.com/juju/pubsub/v2"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/DavinZhang/juju/apiserver/params"
	"github.com/DavinZhang/juju/apiserver/websocket/websockettest"
	"github.com/DavinZhang/juju/state"
	coretesting "github.com/DavinZhang/juju/testing"
	"github.com/DavinZhang/juju/testing/factory"
)

type pubsubSuite struct {
	apiserverBaseSuite
	machineTag names.Tag
	password   string
	nonce      string
	hub        *pubsub.StructuredHub
	pubsubURL  string
}

var _ = gc.Suite(&pubsubSuite{})

func (s *pubsubSuite) SetUpTest(c *gc.C) {
	s.apiserverBaseSuite.SetUpTest(c)
	s.nonce = "nonce"
	m, password := s.Factory.MakeMachineReturningPassword(c, &factory.MachineParams{
		Nonce: s.nonce,
		Jobs:  []state.MachineJob{state.JobManageModel},
	})
	s.machineTag = m.Tag()
	s.password = password
	s.hub = s.config.Hub

	address := s.server.Listener.Addr().String()
	path := fmt.Sprintf("/model/%s/pubsub", s.State.ModelUUID())
	pubsubURL := &url.URL{
		Scheme: "wss",
		Host:   address,
		Path:   path,
	}
	s.pubsubURL = pubsubURL.String()
}

func (s *pubsubSuite) TestNoAuth(c *gc.C) {
	s.checkAuthFails(c, nil, http.StatusUnauthorized, "authentication failed: no credentials provided")
}

func (s *pubsubSuite) TestRejectsUserLogins(c *gc.C) {
	user := s.Factory.MakeUser(c, &factory.UserParams{Password: "sekrit"})
	header := jujuhttp.BasicAuthHeader(user.Tag().String(), "sekrit")
	s.checkAuthFails(c, header, http.StatusForbidden, "authorization failed: user username-.* is not a controller")
}

func (s *pubsubSuite) TestRejectsNonServerMachineLogins(c *gc.C) {
	m, password := s.Factory.MakeMachineReturningPassword(c, &factory.MachineParams{
		Nonce: "a-nonce",
		Jobs:  []state.MachineJob{state.JobHostUnits},
	})
	header := jujuhttp.BasicAuthHeader(m.Tag().String(), password)
	header.Add(params.MachineNonceHeader, "a-nonce")
	s.checkAuthFails(c, header, http.StatusForbidden, "authorization failed: machine .* is not a controller")
}

func (s *pubsubSuite) TestRejectsBadPassword(c *gc.C) {
	header := jujuhttp.BasicAuthHeader(s.machineTag.String(), "wrong")
	header.Add(params.MachineNonceHeader, s.nonce)
	s.checkAuthFails(c, header, http.StatusUnauthorized, "authentication failed: invalid entity name or password")
}

func (s *pubsubSuite) TestRejectsIncorrectNonce(c *gc.C) {
	header := jujuhttp.BasicAuthHeader(s.machineTag.String(), s.password)
	header.Add(params.MachineNonceHeader, "wrong")
	s.checkAuthFails(c, header, http.StatusUnauthorized, "authentication failed: machine 0 not provisioned")
}

func (s *pubsubSuite) checkAuthFails(c *gc.C, header http.Header, code int, message string) {
	conn, resp, err := s.dialWebsocketInternal(c, header)
	c.Assert(err, gc.Equals, websocket.ErrBadHandshake)
	c.Assert(conn, gc.IsNil)
	c.Assert(resp, gc.NotNil)
	defer resp.Body.Close()
	c.Check(resp.StatusCode, gc.Equals, code)
	out, err := ioutil.ReadAll(resp.Body)
	c.Assert(err, jc.ErrorIsNil)
	c.Assert(string(out), gc.Matches, message+"\n")
}

func (s *pubsubSuite) TestMessage(c *gc.C) {
	messages := []params.PubSubMessage{}
	done := make(chan struct{})
	loggo.GetLogger("pubsub").SetLogLevel(loggo.TRACE)
	loggo.GetLogger("juju.apiserver").SetLogLevel(loggo.TRACE)
	_, err := s.hub.SubscribeMatch(pubsub.MatchAll, func(topic string, data map[string]interface{}) {
		c.Logf("topic: %q, data: %v", topic, data)
		messages = append(messages, params.PubSubMessage{
			Topic: topic,
			Data:  data,
		})
		done <- struct{}{}
	})
	c.Assert(err, jc.ErrorIsNil)

	conn := s.dialWebsocket(c)
	defer conn.Close()

	// Read back the nil error, indicating that all is well.
	websockettest.AssertJSONInitialErrorNil(c, conn)

	message1 := params.PubSubMessage{
		Topic: "first",
		Data: map[string]interface{}{
			"origin":  "other",
			"message": "first message",
		}}
	err = conn.WriteJSON(&message1)
	c.Assert(err, jc.ErrorIsNil)

	message2 := params.PubSubMessage{
		Topic: "second",
		Data: map[string]interface{}{
			"origin": "other",
			"value":  false,
		}}
	err = conn.WriteJSON(&message2)
	c.Assert(err, jc.ErrorIsNil)

	select {
	case <-done:
	case <-time.After(coretesting.LongWait):
		c.Fatalf("no first message")
	}

	select {
	case <-done:
	case <-time.After(coretesting.LongWait):
		c.Fatalf("no second message")
	}

	// Close connection.
	err = conn.Close()
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(messages, jc.DeepEquals, []params.PubSubMessage{message1, message2})
}

func (s *pubsubSuite) dialWebsocket(c *gc.C) *websocket.Conn {
	conn, _, err := s.dialWebsocketInternal(c, s.makeAuthHeader())
	c.Assert(err, jc.ErrorIsNil)
	return conn
}

func (s *pubsubSuite) dialWebsocketInternal(c *gc.C, header http.Header) (*websocket.Conn, *http.Response, error) {
	return dialWebsocketFromURL(c, s.pubsubURL, header)
}

func (s *pubsubSuite) makeAuthHeader() http.Header {
	header := jujuhttp.BasicAuthHeader(s.machineTag.String(), s.password)
	header.Add(params.MachineNonceHeader, s.nonce)
	return header
}
