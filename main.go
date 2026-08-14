package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	logPath := filepath.Join(app.configDir, "De-Pulse.log")
	if logFile, logErr := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); logErr == nil {
		defer logFile.Close()
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
		log.Printf("starting %s %s on %s/%s", appName, appVersion, runtime.GOOS, runtime.GOARCH)
	} else {
		log.Printf("unable to open application log %s: %v", logPath, logErr)
	}
	if focusExisting(app.configDir) {
		return
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	server := &http.Server{Handler: app.routes(), ReadHeaderTimeout: 10 * time.Second}
	app.server = server
	rawURL := fmt.Sprintf("http://127.0.0.1:%d/", port)
	log.Printf("Local terminal: %s", rawURL)
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	}()
	time.Sleep(120 * time.Millisecond)
	exe, _ := os.Executable()
	iconPath := filepath.Join(filepath.Dir(filepath.Dir(exe)), "Resources", "DePulse.icns")
	windowPID := openAppWindow(rawURL, iconPath, app.configDir)
	writeInstance(app.configDir, rawURL, windowPID)
	if app.state.Settings.AutoStart {
		go func() { time.Sleep(time.Second); _ = app.engine.Start() }()
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	app.engine.Stop()
	if app.persistence != nil {
		_ = app.persistence.Close()
	}
	_ = os.Remove(instancePath(app.configDir))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
