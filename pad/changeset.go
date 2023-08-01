package pad

import (
	"context"
	"fmt"
	"time"

	ch "github.com/ITLab-CC/bigbluebutton-bot/pad/changesetproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChangesetClient struct {
	addr   string
	conn   *grpc.ClientConn
	client ch.ChangesetClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewChangesetClient(addr string) *ChangesetClient {
	return &ChangesetClient{addr: addr}
}

func (cc *ChangesetClient) Connect() error {
	var err error
	cc.conn, err = grpc.Dial(cc.addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	cc.client = ch.NewChangesetClient(cc.conn)
	cc.ctx, cc.cancel = context.WithTimeout(context.Background(), time.Second * 5)
	return nil
}

func (cc *ChangesetClient) Close() {
	cc.cancel()
	cc.conn.Close()
}

func (cc *ChangesetClient) autoConnect() error {
	if cc.conn == nil {
		return cc.Connect()
	}
	return nil
}

func (cc *ChangesetClient) GenerateChangeset(oldtext string, newtext string, attribs string) (string, error) {
	err := cc.autoConnect()
	if err != nil {
		return "", err
	}

	// //remove all \n from oldtext and newtext
	// oldtext = strings.Replace(oldtext, "\n", "", -1)
	// newtext = strings.Replace(newtext, "\n", "", -1)

	r, err := cc.client.Generate(cc.ctx, &ch.GenerateRequest{
		Oldtext: oldtext,
		Newtext: newtext,
		Attribs: attribs,
	})
	if err != nil {
		return "", fmt.Errorf("could not greet: %v", err)
	}
	return r.Changeset, nil
}

// func main() {
// 	client := NewChangesetClient("localhost:50051")
// 	defer client.Close()

// 	changeset, err := client.GenerateChangeset("Hello\n", "Hello World\n")
// 	if err != nil {
// 		log.Fatalf("could not generate changeset: %v", err)
// 	}
// 	fmt.Println(changeset)
// }
