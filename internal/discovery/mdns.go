package discovery

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/grandcat/zeroconf"
)

// MDNSService wraps a zeroconf server instance.
type MDNSService struct {
	server *zeroconf.Server
}

// StartMDNS registers and starts the Glide mDNS service.
func StartMDNS(ctx context.Context, tcpPort int, udpPort int) (*MDNSService, error) {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Unknown"
	}

	instanceName := "Glide - " + hostName
	serviceType := "_glide._tcp"
	domain := "local."

	txtRecords := []string{
		"v=1",
		"tcp=" + itoa(tcpPort),
		"udp=" + itoa(udpPort),
		"name=" + hostName,
		"os=" + runtime.GOOS,
	}

	server, err := zeroconf.Register(
		instanceName,
		serviceType,
		domain,
		tcpPort,
		txtRecords,
		nil,
	)
	if err != nil {
		return nil, err
	}

	mdns := &MDNSService{server: server}

	log.Println("mDNS service started:", instanceName)

	// Shutdown on context cancel
	go func() {
		<-ctx.Done()
		log.Println("Stopping mDNS service")
		server.Shutdown()
	}()

	return mdns, nil
}

// Shutdown stops the mDNS service manually.
func (m *MDNSService) Shutdown() {
	if m.server != nil {
		m.server.Shutdown()
	}
}

// itoa avoids strconv import for tiny integers.
func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	buf := [10]byte{}
	i := len(buf)

	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}

	return string(buf[i:])
}
