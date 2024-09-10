package main

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

//go build -o conn_openvpn
//sudo ./conn_openvpn --conffile=/home/george/mabat/mabat_docs/vpn/sslvpn-gh106-client-config.ovpn --authfile='/home/george/mabat/mabat_docs/vpn/cred.txt'

var (
	conffile = flag.String("conffile", "", "openVPN client conf file")
	authfile = flag.String("authfile", "", "openVPN client credentials file")
	ipsStr = flag.String("ips", "182.48.202.50,10.255.0.1", "IPs to conection, listed separated by commas, default is <182.48.202.50,10.255.0.1>")
)

func main() {
	flag.Parse()
	slog.Info("Strted with:")
	slog.Info("Configuration", "file", conffile)
	slog.Info("Auth:", "file", authfile)

	closeCh := make(chan int, 2)
	ips := strings.Split(strings.ReplaceAll(*ipsStr, " ", ""), ",")

	go checkConnection(closeCh, ips)

	for {
		ctx, cancel := context.WithCancel(context.Background())
		go startOpenVPN(ctx)

		<- closeCh
		cancel()
	}
}

// checking connection using curl request
func checkConnection(closeCh chan int, ips []string) string {
	for {
		time.Sleep(20*time.Second)
		slog.Info("Checking Connection...")

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		cmd := exec.CommandContext(ctx, "curl", "ipinfo.io")

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stdout
	
		err := cmd.Run()
		if err != nil {
			slog.Error("CMD Error of checkConnection", "error", err)
			closeCh <- 1
			continue
		}

		res := out.String()

		println(res)

		if !strings.Contains(res, "hostname") {
			slog.Info("Reconnection...", "reason", "hostname not found")
			closeCh <- 1
			continue
		}

		isContainsAllIps := true
		// checkiscontains:
		for _, ip := range ips {
			if !strings.Contains(res, ip) {
				isContainsAllIps = false
				break
			}
		}

		if len(ips) > 0 && !isContainsAllIps {
			slog.Info("Reconnection...", "reason", "ip not found", "ips", ips)
			closeCh <- 1
			continue
		}

		slog.Info("Connected!")
		return "ok"
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