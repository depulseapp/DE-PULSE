package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type instanceInfo struct {
	URL        string `json:"url"`
	BackendPID int    `json:"backendPid"`
	WindowPID  int    `json:"windowPid"`
}

func instancePath(configDir string) string { return filepath.Join(configDir, "instance.json") }
func focusProcess(pid int) {
	if pid <= 0 {
		return
	}
	switch runtime.GOOS {
	case "darwin":
		_ = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "System Events" to set frontmost of first process whose unix id is %d to true`, pid)).Run()
	case "windows":
		_ = exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf(`$p=Get-Process -Id %d -ErrorAction SilentlyContinue; if($p){(New-Object -ComObject WScript.Shell).AppActivate($p.Id)}`, pid)).Run()
	}
}
func focusExisting(configDir string) bool {
	data, err := os.ReadFile(instancePath(configDir))
	if err != nil {
		return false
	}
	var in instanceInfo
	if json.Unmarshal(data, &in) != nil || in.URL == "" {
		return false
	}
	c := http.Client{Timeout: 600 * time.Millisecond}
	resp, err := c.Get(strings.TrimRight(in.URL, "/") + "/api/health")
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	if resp.StatusCode != 200 {
		return false
	}
	focusProcess(in.WindowPID)
	return true
}
func writeInstance(configDir, rawURL string, windowPID int) {
	data, _ := json.Marshal(instanceInfo{URL: rawURL, BackendPID: os.Getpid(), WindowPID: windowPID})
	_ = atomicWrite(instancePath(configDir), data, 0600)
}
func openAppWindow(rawURL, iconPath, configDir string) int {
	if os.Getenv("DEPULSE_HEADLESS") == "1" {
		return 0
	}
	if os.Getenv("PMT_NO_BROWSER") == "1" {
		return 0
	}
	switch runtime.GOOS {
	case "darwin":

		script := fmt.Sprintf(`ObjC.import('Cocoa'); ObjC.import('WebKit');
const app=$.NSApplication.sharedApplication;
app.setActivationPolicy($.NSApplicationActivationPolicyRegular);
let mainWindow=null;
ObjC.registerSubclass({
  name:'DePulseDelegateV121',
  superclass:'NSObject',
  // Do not declare NSApplicationDelegate formally here. On current macOS JXA
  // the bridge can fail protocol-name resolution with "protocol does not exist"
  // before the window starts. NSApplication accepts an NSObject delegate that
  // implements the selectors below; formal protocol conformance is unnecessary.
  methods:{
    'applicationShouldTerminateAfterLastWindowClosed:':{
      types:['bool',['id']],
      implementation:function(sender){return true;}
    },
    'applicationShouldHandleReopen:hasVisibleWindows:':{
      types:['bool',['id','bool']],
      implementation:function(sender,hasVisible){
        if(mainWindow){mainWindow.makeKeyAndOrderFront(null);app.activateIgnoringOtherApps(true);}
        return true;
      }
    }
  }
});
const delegate=$.DePulseDelegateV121.alloc.init;
app.delegate=delegate;
const mainMenu=$.NSMenu.alloc.init;
const appMenuItem=$.NSMenuItem.alloc.init;
mainMenu.addItem(appMenuItem);
const appMenu=$.NSMenu.alloc.initWithTitle('DE.PULSE');
const quitItem=$.NSMenuItem.alloc.initWithTitleActionKeyEquivalent('Quit DE.PULSE','terminate:','q');
appMenu.addItem(quitItem);
appMenuItem.submenu=appMenu;
// WKWebView text fields rely on the Cocoa responder chain for standard
// Command-X/C/V/A shortcuts. Expose a normal Edit menu so API-key/password
// fields support copy/paste/select-all just like a regular macOS app.
const editMenuItem=$.NSMenuItem.alloc.init;
mainMenu.addItem(editMenuItem);
const editMenu=$.NSMenu.alloc.initWithTitle('Edit');
const cutItem=$.NSMenuItem.alloc.initWithTitleActionKeyEquivalent('Cut','cut:','x');
const copyItem=$.NSMenuItem.alloc.initWithTitleActionKeyEquivalent('Copy','copy:','c');
const pasteItem=$.NSMenuItem.alloc.initWithTitleActionKeyEquivalent('Paste','paste:','v');
const selectAllItem=$.NSMenuItem.alloc.initWithTitleActionKeyEquivalent('Select All','selectAll:','a');
editMenu.addItem(cutItem);
editMenu.addItem(copyItem);
editMenu.addItem(pasteItem);
editMenu.addItem($.NSMenuItem.separatorItem);
editMenu.addItem(selectAllItem);
editMenuItem.submenu=editMenu;
app.mainMenu=mainMenu;
const frame=$.NSMakeRect(0,0,1440,900);
const style=$.NSWindowStyleMaskTitled|$.NSWindowStyleMaskClosable|$.NSWindowStyleMaskMiniaturizable|$.NSWindowStyleMaskResizable;
const win=$.NSWindow.alloc.initWithContentRectStyleMaskBackingDefer(frame,style,$.NSBackingStoreBuffered,false);
mainWindow=win;
win.title='DE.PULSE v%s';
win.minSize=$.NSMakeSize(820,650);
win.center;
const config=$.WKWebViewConfiguration.alloc.init;
const content=win.contentView;
const web=$.WKWebView.alloc.initWithFrameConfiguration(content.bounds,config);
web.autoresizingMask=$.NSViewWidthSizable|$.NSViewHeightSizable;
content.addSubview(web);
const u=$.NSURL.URLWithString('%s');
web.loadRequest($.NSURLRequest.requestWithURL(u));
try{const icon=$.NSImage.alloc.initWithContentsOfFile('%s');if(icon)app.applicationIconImage=icon;}catch(e){}
app.finishLaunching;
win.makeKeyAndOrderFront(null);
app.activateIgnoringOtherApps(true);
app.run;`, appVersion, strings.ReplaceAll(rawURL, "'", "%27"), strings.ReplaceAll(iconPath, "'", "\\'"))
		windowLogPath := filepath.Join(configDir, "native-window.log")
		windowLog, err := os.OpenFile(windowLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			log.Printf("native window log unavailable: %v", err)
		}
		cmd := exec.Command("/usr/bin/osascript", "-l", "JavaScript", "-e", script)
		if windowLog != nil {
			_, _ = fmt.Fprintf(windowLog, "\n--- DE.PULSE %s window start %s ---\n", appVersion, time.Now().Format(time.RFC3339))
			cmd.Stdout = windowLog
			cmd.Stderr = windowLog
		}
		if err := cmd.Start(); err == nil {
			go func() {
				err := cmd.Wait()
				if err != nil {
					log.Printf("native window exited with error: %v (details: %s)", err, windowLogPath)
				}
				if windowLog != nil {
					_ = windowLog.Close()
				}
				p, _ := os.FindProcess(os.Getpid())
				_ = p.Signal(os.Interrupt)
			}()
			return cmd.Process.Pid
		}
		if windowLog != nil {
			_ = windowLog.Close()
		}
		log.Printf("unable to start native macOS window")
	case "windows":
		paths := []string{filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Microsoft", "Edge", "Application", "msedge.exe"), filepath.Join(os.Getenv("PROGRAMFILES"), "Microsoft", "Edge", "Application", "msedge.exe")}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				cmd := exec.Command(path, "--app="+rawURL, "--new-window", "--class=De-Pulse")
				if cmd.Start() == nil {
					return cmd.Process.Pid
				}
			}
		}
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", rawURL).Start()
	default:
		for _, name := range []string{"chromium", "google-chrome"} {
			if path, err := exec.LookPath(name); err == nil {
				cmd := exec.Command(path, "--app="+rawURL, "--new-window")
				if cmd.Start() == nil {
					return cmd.Process.Pid
				}
		}
		_ = exec.Command("xdg-open", rawURL).Start()
	}
	return 0
}
