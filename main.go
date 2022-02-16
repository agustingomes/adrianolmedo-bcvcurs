package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"
)

var cfg *Config

func main() {
	var (
		addr = flag.String("addr", getNetworkIP(), "Internal address of the container")
		port = flag.Int("port", 80, "Internal port of the container")
	)
	flag.Parse()

	cfg = &Config{
		Addr: *addr,
		Port: *port,
		Logger: NewDebug(func(d *Debug) {
			d.timefmt = "2006-01-02 15:04:05"
		}),
	}

	if err := run(); err != nil {
		cfg.Logger.Log("level", "error", "msg", err.Error())
		os.Exit(1)
	}
}

func run() error {
	http.HandleFunc("/", mGET(home()))
	http.HandleFunc("/v1/", mGET(home()))
	http.HandleFunc("/v1/dollar/", mGET(dollar(cfg.Logger)))

	s := &http.Server{
		Addr:           cfg.Addr + ":" + strconv.Itoa(cfg.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	cfg.Logger.Log("level", "info", "addr", cfg.Addr+":"+strconv.Itoa(cfg.Port), "msg", "listening")
	return s.ListenAndServe()
}

// getNetworkIP return local network IP. If you are not connected to IPv4 it will return 127.0.0.1.
func getNetworkIP() string {
	ip := "127.0.0.1"

	netInterfaceAddresses, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}

	for _, netInterfaceAddress := range netInterfaceAddresses {
		networkIP, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIP.IP.IsLoopback() && networkIP.IP.To4() != nil {
			ip := networkIP.IP.String()
			return ip
		}
	}

	return ip
}
