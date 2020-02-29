package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/yottab/cli/config"
	ybApi "github.com/yottab/proto-api/proto"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	flagVariableArray = make([]string, 0, 8)
	flagIndex         int32
	flagAppName       string
)

func arrayFlagToMap(flags []string) map[string]string {
	varMap := make(map[string]string, len(flags))
	for _, v := range flags {
		index := strings.Index(v, "=")
		if index > 0 {
			key := v[:index]
			varMap[key] = v[index+1:]
		} else {
			log.Fatalf("Bad data entry format [%s], Enter the information in 'KEY=VALUE' format.", v)
		}
	}
	return varMap
}

func readFromConsole(inputAnswr string) (val string) {
	fmt.Print(inputAnswr)
	reader := bufio.NewReader(os.Stdin)
	val, _ = reader.ReadString('\n')
	return strings.TrimSpace(val)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func readPasswordFromConsole(inputAnswr string) (val string) {
	fmt.Print(inputAnswr)
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return ""
	}
	password := string(bytePassword)
	return strings.TrimSpace(password)
}

func grpcConnect() ybApi.Client {
	return ybApi.Connect(
		viper.GetString(config.KEY_HOST),
		ybApi.NewJwtAccess(func() string {
			return viper.GetString(config.KEY_TOKEN)
		}))
}

func toTime(t *ybApi.Timestamp) (out string) {
	if t != nil {
		out = time.Unix(t.Seconds, 0).Format(time.RFC3339)
	}
	return
}

func endpointTypeValid(etype string) error {
	switch etype {
	case "http":
		return nil
	case "grpc":
		return nil
	default:
		return errors.New("Endpoint type is invalid, valid values are http, grpc")
	}
}

func streamAppLog(args []string) {
	var client ybApi.Client
	firstTry := true
	for {
		if !firstTry {
			//Wait and retry
			time.Sleep(time.Millisecond * 500)
		}
		client = grpcConnect()
		req := getCliRequestIdentity(args, 0)
		logClient, err := client.V2().AppLog(context.Background(), req)
		uiCheckErr("Could not Get Application log", err)
		err = uiStreamLog(logClient)
		client.Close()
		if err != nil {
			log.Debug(err)
			if status.Code(err) == codes.ResourceExhausted {
				break
			}
			if strings.Contains(err.Error(), "RST_STREAM") {
				//Resume log streaming on proto related error
				log.Fatal("RST_STREAM")
				continue
			}
		} else {
			break
		}

	}

}

func streamBuildLog(appName, appTag string) {
	var client ybApi.Client
	firstTry := true
	id := getRequestIdentity(
		fmt.Sprintf(pushLogIDFormat, appName, appTag))
	for {
		if !firstTry {
			//Wait and retry
			time.Sleep(time.Millisecond * 500)
		}
		client = grpcConnect()
		logClient, err := client.V2().ImgBuildLog(context.Background(), id)
		uiCheckErr(fmt.Sprintf("Could not get build log right now!\nTry again in a few soconds using:\n$yb push log --name=%s --tag=%s", appName, appTag), err)
		err = uiStreamLog(logClient)
		client.Close()
		if err != nil {
			if status.Code(err) == codes.ResourceExhausted {
				break
			}
			if strings.Contains(err.Error(), "RST_STREAM") {
				//Resume log streaming on proto related error
				continue
			}
		} else {
			break
		}
	}

}
