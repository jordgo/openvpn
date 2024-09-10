package openvpn

import (
	"bytes"
	"context"
	"flag"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)


var (
	conffile = flag.String("conffile", "", "openVPN client conf file")
	authfile = flag.String("authfile", "", "openVPN client credentials file")
)

func main() {
	closeCh := make(chan int, 2)

	go checkConnection(closeCh)

	for {
		ctx, cancel := context.WithCancel(context.Background())
		go startOpenVPN(ctx)

		<- closeCh
		cancel()
	}
}

// checking connection using curl request
func checkConnection(closeCh chan int) string {
	for {
		time.Sleep(15*time.Second)
		slog.Info("Checking Connection...")

		cmd := exec.Command("curl", "ipinfo.io")

		var out bytes.Buffer
		cmd.Stdout = &out
	
		err := cmd.Run()
		if err != nil {
			slog.Error("CMD Error of checkConnection", "error", err)
		}

		res := out.String()

		println(res)

		if strings.Contains(res, "hostname") {
			slog.Info("Connected!")
			return "ok"
		}
		
		closeCh <- 1
	}


}

// execute cmd command to start openVPN
func startOpenVPN(ctx context.Context) error {
	cmd := exec.CommandContext(ctx,
		"openvpn", 
		"--config", *conffile, 
		"--auth-user-pass", *authfile,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		slog.Error("CMD Error of startOpenVPN", "error", err)
	}

	return err
} 