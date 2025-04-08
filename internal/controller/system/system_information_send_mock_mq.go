package system

import (
	"context"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/utility"

	"unibee/api/system/information"
)

func (c *ControllerInformation) SendMockMQ(ctx context.Context, req *information.SendMockMQReq) (res *information.SendMockMQRes, err error) {
	utility.Assert(len(req.Message) > 0, "invalid message")
	_, _ = redismq.Send(&redismq.Message{
		Topic:      redismq2.TopicTest1.Topic,
		Tag:        redismq2.TopicTest1.Tag,
		Body:       req.Message,
		CustomData: map[string]interface{}{"CreateFrom": utility.ReflectCurrentFunctionName()},
	})
	return &information.SendMockMQRes{}, nil
}
