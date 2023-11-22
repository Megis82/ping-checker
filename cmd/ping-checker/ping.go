package main

import (
	"context"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Megis82/ping-checker/internal/config"
	logger "github.com/Megis82/ping-checker/internal/log"
	probing "github.com/Megis82/pro-bing"
	"go.uber.org/zap"
)

// var usage = `
// Usage:

//     ping [-c count] [-i interval] [-t timeout] [-rt receiveTimeout] [--privileged] host

// Examples:

//     # ping google continuously
//     ping www.google.com

//     # ping google 5 times
//     ping -c 5 www.google.com

//     # ping google 5 times at 500ms intervals
//     ping -c 5 -i 500ms www.google.com

//     # ping google for 10 seconds
//     ping -t 10s www.google.com

// 	# ping google with received timeout 2 seconds
//     ping -rt 2s www.google.com

//     # Send a privileged raw ICMP ping
//     sudo ping --privileged www.google.com

//     # Send ICMP messages with a 100-byte payload
//     ping -s 100 1.1.1.1
// `

// type pingerParameters struct {
// 	timeout        time.Duration
// 	receiveTimeout time.Duration
// 	interval       time.Duration
// 	count          int
// 	size           int
// 	ttl            int
// 	privileged     bool
// }

func main() {

	cfg, err := config.Init()
	if err != nil {
		return
	}

	logger, err := logger.NewLogger(cfg.LogFileName)
	if err != nil {
		return
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	for _, host := range cfg.RequestsAddresses {
		loopHost := strings.TrimSpace(host)
		fmt.Println(loopHost)
		go func() {
			runPinger(ctx, loopHost, cfg.ReceiveTimeout, logger)
		}()
	}
	<-ctx.Done()
}

func runPinger(ctx context.Context, host string, receiveTimeout time.Duration, logger *zap.Logger) {

	pinger, err := probing.NewPinger(host)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}

	// pinger.OnRecv = func(pkt *probing.Packet) {
	// 	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v\n",
	// 		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	// }
	pinger.OnLost = func(pkt *probing.Packet) {

		slog := fmt.Sprintf("%d bytes to %s (%s) was lost: icmp_seq=%d\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Addr, pkt.Seq)
		logger.Info(slog)
	}
	// pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
	// 	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
	// 		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	// }
	// pinger.OnFinish = func(stats *probing.Statistics) {
	// 	fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
	// 	fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
	// 		stats.PacketsSent, stats.PacketsRecv, stats.PacketsRecvDuplicates, stats.PacketLoss)
	// 	// fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
	// 	// 	stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	// }

	pinger.Count = -1
	pinger.Size = 32
	pinger.Interval = time.Second
	pinger.Timeout = time.Second * 100000
	pinger.ReceiveTimeout = time.Millisecond * 300
	pinger.TTL = 64
	pinger.SetPrivileged(true)

	// fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	err = pinger.RunWithContext(ctx)
	if err != nil {
		fmt.Println("Failed to ping target host:", err)
	}

}
