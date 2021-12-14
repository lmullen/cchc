package messages_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lmullen/cchc/common/messages"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/rabbitmq"
	"github.com/stretchr/testify/require"
)

func TestRoundtrip(t *testing.T) {
	t.Parallel()

	user := "gnomock"
	pass := "strong-passwords-are-the-best"

	// gnomock setup
	p := rabbitmq.Preset(rabbitmq.WithUser(user, pass))

	container, err := gnomock.Start(p)
	require.NoError(t, err)
	defer func() { require.NoError(t, gnomock.Stop(container)) }()

	// Connect to RabbitMQ and make sure we can disconnect, both without errors
	connstr := fmt.Sprintf("amqp://%s:%s@%s", user, pass, container.DefaultAddress())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var repo messages.Repository

	repo, err = messages.NewRabbitMQ(ctx, connstr, "testq", 20)

	require.NoError(t, err)
	defer func() { require.NoError(t, repo.Close()) }()

	send := messages.NewFullTextMsg(
		uuid.New(),
		"http://www.loc.gov/item/magbell.00510307/",
		"Behold, I am sending you a message.")

	err = repo.Send(context.TODO(), send)
	require.NoError(t, err)

	consumer := repo.Consume()
	incoming := <-consumer
	var receive *messages.FullTextPredict = &messages.FullTextPredict{}
	err = json.Unmarshal(incoming.Body, receive)
	require.NoError(t, err)

	require.Equal(t, send, receive)

}
