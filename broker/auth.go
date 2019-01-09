/* Copyright (c) 2018, joy.zhou <chowyu08@gmail.com>
 */
package broker

import (
	"strings"
)

const (
	PUB = 1
	SUB = 2
)

type Checker interface {
	CheckTopicAuth(typ int, client ClientInfo, topic string) bool
}

func (c *client) CheckTopicAuth(typ int, topic string) bool {
	if c.typ != CLIENT || !c.broker.config.Acl {
		return true
	}
	if strings.HasPrefix(topic, "$queue/") {
		topic = string([]byte(topic)[7:])
		if topic == "" {
			return false
		}
	}
	return c.broker.Auth.CheckTopicAuth(typ, c.info, topic)

}
