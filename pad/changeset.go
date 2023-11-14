package pad

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	ch "github.com/ITLab-CC/bigbluebutton-bot/pad/changesetproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/connectivity"
	git "gopkg.in/src-d/go-git.v4" // with go modules disabled
)

type ChangesetClient struct {
	ip   string
	port string

	Changsetserverpath string

	changsetServerProcess *exec.Cmd

	Downloadeurl string

	conn   *grpc.ClientConn
	client ch.ChangesetClient
	ctx    context.Context
	cancel context.CancelFunc
}

func NewChangesetClient(ip string, port string) *ChangesetClient {
	return &ChangesetClient{
		ip:   ip,
		port: port,

		Changsetserverpath: "./.changsetserver",
		Downloadeurl:       "https://github.com/bigbluebutton-bot/changeset-grpc",
	}
}

type submoduleInfo struct {
	Path string
	URL  string
}

func (cc *ChangesetClient) extractInfoFromSubmoduleFile(filename string) (*submoduleInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var info submoduleInfo

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "	path = ") {
			info.Path = strings.TrimPrefix(line, "	path = ")
		} else if strings.HasPrefix(line, "	url = ") {
			info.URL = strings.TrimPrefix(line, "	url = ")
		}
	}

	if info.Path == "" || info.URL == "" {
		return nil, fmt.Errorf("Pfad oder URL nicht gefunden")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &info, nil
}

// Downloade changeset server files
func (cc *ChangesetClient) downloadAndInstallChangesetServer() error {
	// if folder cc.Changsetserverpath error
	if _, err := os.Stat(cc.Changsetserverpath); err == nil {
		// if folder exits return error
		return fmt.Errorf("folder %s already exists", cc.Changsetserverpath)
	}

	// download changeset server files
	_, err := git.PlainClone(cc.Changsetserverpath, false, &git.CloneOptions{
		URL:      cc.Downloadeurl,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("could not download changeset server files (%s): %v", cc.Downloadeurl, err)
	}

	// get url from submodule file
	submoduleinfo, err := cc.extractInfoFromSubmoduleFile(cc.Changsetserverpath + "/.gitmodules")
	if err != nil {
		return err
	}

	// downloade submodule files
	_, err = git.PlainClone(cc.Changsetserverpath+"/"+submoduleinfo.Path, false, &git.CloneOptions{
		URL:      submoduleinfo.URL,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("could not download changeset server files (%s): %v", submoduleinfo.URL, err)
	}

	// test if npm is installed
	_, err = exec.LookPath("npm")
	if err != nil {
		return fmt.Errorf("npm is not installed")
	}

	// install node modules
	cmd := exec.Command("npm", "install")
	cmd.Dir = cc.Changsetserverpath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error while installing node modules: %v", err)
	}

	// install etherpad
	if err := installEtherpad(cc.Changsetserverpath); err != nil {
		return err
	}

	return nil
}

func (cc *ChangesetClient) StartChangesetServer() error {
	// test if node is installed
	_, err := exec.LookPath("node")
	if err != nil {
		return fmt.Errorf("node is not installed")
	}

	// if file exists
	if _, err := os.Stat(cc.Changsetserverpath + "/server.js"); err != nil {
		// if file not exists download it
		if err := cc.downloadAndInstallChangesetServer(); err != nil {
			return err
		}
	}

	// Command to execute and start the changeset server
	cc.changsetServerProcess = exec.Command("node", "server.js", cc.ip, cc.port)
	cc.changsetServerProcess.Dir = cc.Changsetserverpath
	cc.changsetServerProcess.Stdout = os.Stdout
	cc.changsetServerProcess.Stderr = os.Stderr
	go func() {
		// run the command
		err = cc.changsetServerProcess.Run()
		if err != nil {
			fmt.Println("error while starting the changeset server: ", err)
			return
		}
	}()

	//Kill the changeset node server if the golang programm exits
	signals := make(chan os.Signal, 1)
	// Hören Sie auf SIGINT und SIGTERM
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		if cc.changsetServerProcess.Process != nil {
			cc.changsetServerProcess.Process.Kill()
		}
		os.Exit(1)
	}()

	// try to connect 10 times
	for i := 0; i < 10; i++ {
		err = cc.autoConnect()
		if err != nil {
			fmt.Println("could not connect to changeset server: ", err)
		}
		// send ping
		_, err = cc.client.Ping(cc.ctx, &ch.Nothing{})
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)

		if i >= 9 {
			return fmt.Errorf("could not start and connect to changeset server")
		}
	}

	fmt.Println("Changeset server started successfully")
	return nil
}

func (cc *ChangesetClient) StopChangesetServer() {
	if cc.changsetServerProcess.Process != nil {
		cc.changsetServerProcess.Process.Kill()
	}
}

func (cc *ChangesetClient) Connect() error {
	var err error
	cc.conn, err = grpc.Dial(cc.ip+":"+cc.port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	cc.client = ch.NewChangesetClient(cc.conn)
	return nil
}

func (cc *ChangesetClient) Close() {
	cc.cancel()
	cc.conn.Close()
}

func (cc *ChangesetClient) autoConnect() error {
	if cc.conn == nil {
		err := cc.Connect()
		if err != nil {
			return err
		}
	}
	if cc.conn.GetState() != connectivity.Ready {
		cc.conn.Connect()
	}
	cc.ctx, cc.cancel = context.WithTimeout(context.Background(), time.Second*5)
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
