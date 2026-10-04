package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/Dien45/Mobile-agent/internal/server"
	"github.com/Dien45/Mobile-agent/internal/store"
)

const version = "0.1.0-dev"

func main(){if err:=run();err!=nil{fmt.Fprintln(os.Stderr,"arka:",err);os.Exit(1)}}
func run()error{
	command:="start";if len(os.Args)>1{command=os.Args[1]}
	switch command{
	case "version","--version","-v":fmt.Println("arka",version,runtime.GOOS+"/"+runtime.GOARCH);return nil
	case "doctor":return doctor()
	case "start":return start(hasArg("--no-open"))
	case "open":return openBrowser("http://127.0.0.1:7331")
	default:return fmt.Errorf("unknown command %q (use start, open, doctor, or version)",command)
	}
}
func start(noOpen bool)error{
	home:=server.StateDir();if err:=os.MkdirAll(home,0o700);err!=nil{return err};st,err:=store.Open(filepath.Join(home,"state.json"));if err!=nil{return err};workspace,err:=os.Getwd();if err!=nil{return err};app:=server.New(st,filepath.Join(home,"skills"),workspace);addr:="127.0.0.1:7331";url:=app.URL(addr)
	fmt.Println("Arka local daemon is starting")
	fmt.Println("Workspace:",workspace)
	fmt.Println("Open:",url)
	if !noOpen{go func(){time.Sleep(350*time.Millisecond);_ = openBrowser(url)}()}
	errCh:=make(chan error,1);go func(){errCh<-app.Listen(addr)}();sigCh:=make(chan os.Signal,1);signal.Notify(sigCh,os.Interrupt,syscall.SIGTERM)
	select{case sig:=<-sigCh:fmt.Println("Stopping after",sig);ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second);defer cancel();return app.Shutdown(ctx);case err:=<-errCh:if err==http.ErrServerClosed{return nil};return err}
}
func doctor()error{home:=server.StateDir();workspace,_:=os.Getwd();fmt.Println("Arka doctor");fmt.Println("  version:   ",version);fmt.Println("  platform:  ",runtime.GOOS+"/"+runtime.GOARCH);fmt.Println("  state:     ",home);fmt.Println("  workspace: ",workspace);fmt.Println("  git:       ",lookup("git"));fmt.Println("  shell:     ",shell());fmt.Println("  loopback:   127.0.0.1:7331");return nil}
func lookup(name string)string{p,err:=exec.LookPath(name);if err!=nil{return "not found"};return p}
func shell()string{if runtime.GOOS=="windows"{return lookup("powershell.exe")};if v:=os.Getenv("SHELL");v!=""{return v};return "/bin/sh"}
func openBrowser(url string)error{var cmd *exec.Cmd;switch runtime.GOOS{case "windows":cmd=exec.Command("cmd","/c","start","",url);case "darwin":cmd=exec.Command("open",url);default:if os.Getenv("TERMUX_VERSION")!=""{cmd=exec.Command("termux-open-url",url)}else{cmd=exec.Command("xdg-open",url)}};return cmd.Start()}
func hasArg(want string)bool{for _,v:=range os.Args[2:]{if v==want{return true}};return false}
