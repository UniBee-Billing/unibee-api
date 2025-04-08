package mock

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	redismq "github.com/jackyang-hk/go-redismq"
	redismq2 "unibee/internal/cmd/redismq"
	"unibee/utility"
)

type TestMessageListener struct {
}

func (t TestMessageListener) GetTopic() string {
	return redismq2.TopicTest1.Topic
}

func (t TestMessageListener) GetTag() string {
	return redismq2.TopicTest1.Tag
}

func (t TestMessageListener) Consume(ctx context.Context, message *redismq.Message) redismq.Action {
	utility.Assert(len(message.Body) > 0, "body is nil")
	utility.Assert(len(message.Body) != 0, "body length is 0")
	g.Log().Infof(ctx, "TestMessageListener Receive Message:%s", message.Body)
	return redismq.CommitMessage
}

func init() {
	redismq.RegisterListener(New())
	fmt.Println("TestMessageListener RegisterListener")
}

func New() *TestMessageListener {
	return &TestMessageListener{}
}
