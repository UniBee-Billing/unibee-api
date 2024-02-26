package onetime

import (
	"context"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/redismq"

	"unibee/api/onetime/mock"
)

func (c *ControllerMock) MockMessageSend(ctx context.Context, req *mock.MockMessageSendReq) (res *mock.MockMessageSendRes, err error) {
	_, err = redismq.Send(&redismq.Message{
		Topic: redismq2.TopicTest1.Topic,
		Tag:   redismq2.TopicTest1.Tag,
		Body:  req.Message,
	})
	if err != nil {
		return nil, err
	}
	return &mock.MockMessageSendRes{}, nil
}
