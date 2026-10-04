package server

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Dien45/Mobile-agent/internal/agent"
	"github.com/Dien45/Mobile-agent/internal/core"
	"github.com/Dien45/Mobile-agent/internal/provider"
	"github.com/Dien45/Mobile-agent/internal/skills"
	"github.com/Dien45/Mobile-agent/internal/store"
)

//go:embed web/*
var webAssets embed.FS

type Server struct {
	Store *store.Store
	Agent *agent.Manager
	Skills skills.Manager
	Token string
	Workspace string
	http *http.Server
}

func New(st *store.Store, skillRoot, workspace string) *Server {
	var b [24]byte; _,_=rand.Read(b[:])
	return &Server{Store:st,Agent:agent.New(st),Skills:skills.Manager{Root:skillRoot},Token:hex.EncodeToString(b[:]),Workspace:workspace}
}

func (s *Server) Handler() http.Handler {
	mux:=http.NewServeMux()
	mux.HandleFunc("/health",func(w http.ResponseWriter,r *http.Request){writeJSON(w,200,map[string]any{"status":"ok","version":"0.1.0-dev","platform":runtime.GOOS+"/"+runtime.GOARCH})})
	mux.HandleFunc("/api/v1/sessions",s.sessions)
	mux.HandleFunc("/api/v1/sessions/",s.session)
	mux.HandleFunc("/api/v1/providers",s.providers)
	mux.HandleFunc("/api/v1/providers/",s.providerModels)
	mux.HandleFunc("/api/v1/tasks",func(w http.ResponseWriter,r *http.Request){writeJSON(w,200,s.Agent.AllStatus())})
	mux.HandleFunc("/api/v1/skills",s.skillList)
	mux.HandleFunc("/api/v1/skills/",s.skillDelete)
	mux.HandleFunc("/api/v1/terminal/exec",s.terminalExec)
	static,_:=fs.Sub(webAssets,"web"); files:=http.FileServer(http.FS(static))
	mux.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){
		if token:=r.URL.Query().Get("token");token!=""{if token!=s.Token{http.Error(w,"invalid local token",http.StatusForbidden);return};http.SetCookie(w,&http.Cookie{Name:"arka_token",Value:token,Path:"/",HttpOnly:true,SameSite:http.SameSiteStrictMode,MaxAge:86400});http.Redirect(w,r,"/",http.StatusSeeOther);return}
		files.ServeHTTP(w,r)
	})
	return s.security(mux)
}

func (s *Server) security(next http.Handler) http.Handler{return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
	host,_,err:=net.SplitHostPort(r.Host);if err!=nil{host=r.Host};host=strings.Trim(host,"[]")
	if host!="localhost"&&host!="127.0.0.1"&&host!="::1"{http.Error(w,"Arka only accepts loopback hosts",http.StatusForbidden);return}
	if strings.HasPrefix(r.URL.Path,"/api/"){cookie,err:=r.Cookie("arka_token");if err!=nil||cookie.Value!=s.Token{writeJSON(w,401,map[string]string{"error":"open Arka using the authenticated URL printed by the daemon"});return};if origin:=r.Header.Get("Origin");origin!=""&&!strings.Contains(origin,host){writeJSON(w,403,map[string]string{"error":"invalid origin"});return}}
	next.ServeHTTP(w,r)
})}

func (s *Server) Listen(addr string) error{s.http=&http.Server{Addr:addr,Handler:s.Handler(),ReadHeaderTimeout:10*time.Second};return s.http.ListenAndServe()}
func (s *Server) Shutdown(ctx context.Context)error{if s.http==nil{return nil};return s.http.Shutdown(ctx)}

func (s *Server) sessions(w http.ResponseWriter,r *http.Request){switch r.Method{
case http.MethodGet:writeJSON(w,200,s.Store.ListSessions(r.URL.Query().Get("trash")=="1"))
case http.MethodPost:
	var in struct{Title,Provider,Model,Workspace string};if !decode(w,r,&in){return};if in.Title==""{in.Title="New session"};if in.Workspace==""{in.Workspace=s.Workspace};now:=time.Now().UTC();v:=core.Session{ID:core.NewID(),Title:in.Title,Provider:in.Provider,Model:in.Model,Workspace:in.Workspace,Status:"idle",CreatedAt:now,UpdatedAt:now};if err:=s.Store.SaveSession(v);err!=nil{fail(w,err);return};writeJSON(w,201,v)
default:w.WriteHeader(http.StatusMethodNotAllowed)}}

func (s *Server) session(w http.ResponseWriter,r *http.Request){rest:=strings.TrimPrefix(r.URL.Path,"/api/v1/sessions/");parts:=strings.Split(strings.Trim(rest,"/"),"/");if len(parts)==0||parts[0]==""{http.NotFound(w,r);return};id:=parts[0]
	if len(parts)==2{switch parts[1]{case "messages":s.messages(w,r,id);return;case "stop":if r.Method==http.MethodPost{writeJSON(w,200,map[string]bool{"stopped":s.Agent.Stop(id)});return};case "restore":if r.Method==http.MethodPost{if err:=s.Store.RestoreSession(id);err!=nil{fail(w,err);return};writeJSON(w,200,map[string]bool{"ok":true});return}}}
	switch r.Method{
	case http.MethodGet:v,ok:=s.Store.GetSession(id);if !ok{http.NotFound(w,r);return};writeJSON(w,200,map[string]any{"session":v,"messages":s.Store.Messages(id),"task":s.Agent.Status(id)})
	case http.MethodPatch:
		v,ok:=s.Store.GetSession(id);if !ok{http.NotFound(w,r);return};var in struct{Title,Provider,Model string};if !decode(w,r,&in){return};if in.Title!=""{v.Title=in.Title};if in.Provider!=""{v.Provider=in.Provider};if in.Model!=""{v.Model=in.Model};v.UpdatedAt=time.Now().UTC();if err:=s.Store.SaveSession(v);err!=nil{fail(w,err);return};writeJSON(w,200,v)
	case http.MethodDelete:if s.Agent.Status(id).Status=="running"{writeJSON(w,409,map[string]string{"error":"stop the running task before deleting this session"});return};if err:=s.Store.DeleteSession(id,r.URL.Query().Get("purge")=="1");err!=nil{fail(w,err);return};w.WriteHeader(204)
	default:w.WriteHeader(http.StatusMethodNotAllowed)}
}

func (s *Server) messages(w http.ResponseWriter,r *http.Request,id string){if r.Method==http.MethodGet{writeJSON(w,200,s.Store.Messages(id));return};if r.Method!=http.MethodPost{w.WriteHeader(405);return};var in struct{Content string};if !decode(w,r,&in)||strings.TrimSpace(in.Content)==""{return};v:=core.Message{ID:core.NewID(),SessionID:id,Role:"user",Content:strings.TrimSpace(in.Content),CreatedAt:time.Now().UTC()};if err:=s.Store.AddMessage(v);err!=nil{fail(w,err);return};if err:=s.Agent.Start(id);err!=nil{writeJSON(w,409,map[string]string{"error":err.Error()});return};writeJSON(w,202,v)}

func (s *Server) providers(w http.ResponseWriter,r *http.Request){if r.Method==http.MethodGet{writeJSON(w,200,map[string]any{"profiles":s.Store.ListProviders(),"presets":[]string{"openai","anthropic","gemini","openrouter","xai","groq","mistral","azure-openai","ollama","custom-openai"}});return};if r.Method!=http.MethodPost{w.WriteHeader(405);return};var p core.ProviderProfile;if !decode(w,r,&p){return};if p.ID==""{p.ID=core.NewID()};if p.Name==""{p.Name="Custom OpenAI"};if p.Type==""{p.Type="openai-compatible"};if err:=s.Store.SaveProvider(p);err!=nil{fail(w,err);return};p.APIKey="";writeJSON(w,201,p)}
func (s *Server) providerModels(w http.ResponseWriter,r *http.Request){id:=strings.TrimSuffix(strings.TrimPrefix(r.URL.Path,"/api/v1/providers/"),"/models");if !strings.HasSuffix(r.URL.Path,"/models"){http.NotFound(w,r);return};p,ok:=s.Store.Provider(id);if !ok{http.NotFound(w,r);return};models,err:=provider.New().Models(r.Context(),p);if err!=nil{fail(w,err);return};writeJSON(w,200,models)}

func (s *Server) skillList(w http.ResponseWriter,r *http.Request){if r.Method==http.MethodGet{v,err:=s.Skills.List();if err!=nil{fail(w,err);return};writeJSON(w,200,v);return};if r.Method!=http.MethodPost{w.WriteHeader(405);return};var in struct{Source,Name string};if !decode(w,r,&in){return};ctx,cancel:=context.WithTimeout(r.Context(),2*time.Minute);defer cancel();v,err:=s.Skills.Install(ctx,in.Source,in.Name);if err!=nil{fail(w,err);return};writeJSON(w,201,v)}
func (s *Server) skillDelete(w http.ResponseWriter,r *http.Request){if r.Method!=http.MethodDelete{w.WriteHeader(405);return};name:=strings.Trim(strings.TrimPrefix(r.URL.Path,"/api/v1/skills/"),"/");if err:=s.Skills.Remove(name);err!=nil{fail(w,err);return};w.WriteHeader(204)}

func (s *Server) terminalExec(w http.ResponseWriter,r *http.Request){if r.Method!=http.MethodPost{w.WriteHeader(405);return};var in struct{Command,Cwd string};if !decode(w,r,&in){return};cwd:=s.Workspace;if in.Cwd!=""{abs,err:=filepath.Abs(in.Cwd);if err!=nil{fail(w,err);return};root,_:=filepath.Abs(s.Workspace);rel,_:=filepath.Rel(root,abs);if rel==".."||strings.HasPrefix(rel,".."+string(filepath.Separator)){writeJSON(w,403,map[string]string{"error":"cwd escapes workspace"});return};cwd=abs};ctx,cancel:=context.WithTimeout(r.Context(),10*time.Minute);defer cancel();var cmd *exec.Cmd;if runtime.GOOS=="windows"{cmd=exec.CommandContext(ctx,"powershell.exe","-NoProfile","-Command",in.Command)}else{cmd=exec.CommandContext(ctx,"/bin/sh","-lc",in.Command)};cmd.Dir=cwd;out,err:=cmd.CombinedOutput();if len(out)>2<<20{out=out[:2<<20]};code:=0;if err!=nil{code=1};writeJSON(w,200,map[string]any{"output":string(out),"exit_code":code,"error":errorText(err)})}

func decode(w http.ResponseWriter,r *http.Request,v any)bool{r.Body=http.MaxBytesReader(w,r.Body,2<<20);if err:=json.NewDecoder(r.Body).Decode(v);err!=nil{writeJSON(w,400,map[string]string{"error":err.Error()});return false};return true}
func writeJSON(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
func fail(w http.ResponseWriter,err error){writeJSON(w,500,map[string]string{"error":err.Error()})}
func errorText(err error)string{if err==nil{return ""};return err.Error()}
func ParsePort(addr string)int{_,p,_:=net.SplitHostPort(addr);v,_:=strconv.Atoi(p);return v}
func (s *Server) URL(addr string)string{return fmt.Sprintf("http://%s/?token=%s",addr,s.Token)}
func StateDir()string{
	if v:=os.Getenv("ARKA_HOME");v!=""{return v}
	home,_:=os.UserHomeDir()
	if runtime.GOOS=="windows"{if v:=os.Getenv("LOCALAPPDATA");v!=""{return filepath.Join(v,"Arka")}}
	if os.Getenv("TERMUX_VERSION")!=""{return filepath.Join(home,".local","state","arka")}
	if runtime.GOOS=="darwin"{return filepath.Join(home,"Library","Application Support","Arka")}
	if v:=os.Getenv("XDG_STATE_HOME");v!=""{return filepath.Join(v,"arka")}
	return filepath.Join(home,".local","state","arka")
}
