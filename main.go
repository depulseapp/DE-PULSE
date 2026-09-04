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
	if handled, code := runDeveloperCommand(os.Args[1:], os.Stdout, os.Stderr); handled {
		if code != 0 {
			os.Exit(code)
		}
		return
	}
	hosted := isHostedRuntime()
	if hosted {
		if err := validateHostedEnvironment(); err != nil {
			log.Fatalf("hosted environment validation failed: %v", err)
		}
	}
	app, err := NewApplication()
	if err != nil {
		log.Fatal(err)
	}
	if !hosted {
		logPath := filepath.Join(app.configDir, "De-Pulse.log")
		if logFile, logErr := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); logErr == nil {
			defer logFile.Close()
			log.SetOutput(io.MultiWriter(os.Stderr, logFile))
		} else {
			log.Printf("unable to open application log %s: %v", logPath, logErr)
		}
	}
	log.Printf("starting %s %s on %s/%s (%s)", appName, appVersion, runtime.GOOS, runtime.GOARCH, runtimeMode())
	if !hosted && focusExisting(app.configDir) {
		return
	}
	listenAddr := "127.0.0.1:0"
	if hosted {
		listenAddr = hostedListenAddress()
	}
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	server := &http.Server{Handler: securityPerimeter(app.protectDocumentationHTTP(hostedManagedSecretBoundary(providerCredentialBoundary(app, app.routes())))), ReadHeaderTimeout: 10 * time.Second}
	app.server = server
	rawURL := fmt.Sprintf("http://127.0.0.1:%d/", port)
	if hosted {
		log.Printf("Hosted HTTP listener: %s", listener.Addr().String())
	} else {
		log.Printf("Local terminal: %s", rawURL)
	}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	}()
	if !hosted {
		time.Sleep(120 * time.Millisecond)
		exe, _ := os.Executable()
		iconPath := filepath.Join(filepath.Dir(filepath.Dir(exe)), "Resources", "DePulse.icns")
		windowPID := openAppWindow(rawURL, iconPath, app.configDir)
		writeInstance(app.configDir, rawURL, windowPID)
	}
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
	if !hosted {
		_ = os.Remove(instancePath(app.configDir))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
