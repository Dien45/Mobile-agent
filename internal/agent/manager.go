package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Dien45/Mobile-agent/internal/core"
	"github.com/Dien45/Mobile-agent/internal/provider"
	"github.com/Dien45/Mobile-agent/internal/store"
	"github.com/Dien45/Mobile-agent/internal/tools"
)

type running struct { status core.TaskStatus; cancel context.CancelFunc }
type Manager struct {
	mu sync.RWMutex
	running map[string]*running
	store *store.Store
	provider *provider.Client
	maxTools int
}

func New(st *store.Store) *Manager { return &Manager{running:map[string]*running{}, store:st, provider:provider.New(), maxTools:1000} }

func (m *Manager) Status(sessionID string) core.TaskStatus {
	m.mu.RLock(); defer m.mu.RUnlock(); if r,ok:=m.running[sessionID];ok{return r.status}; return core.TaskStatus{SessionID:sessionID,Status:"idle"}
}
func (m *Manager) AllStatus() []core.TaskStatus { m.mu.RLock(); defer m.mu.RUnlock(); out:=make([]core.TaskStatus,0,len(m.running));for _,r:=range m.running{out=append(out,r.status)};return out }
func (m *Manager) Stop(sessionID string) bool { m.mu.Lock(); defer m.mu.Unlock(); if r,ok:=m.running[sessionID];ok{r.cancel();return true};return false }

func (m *Manager) Start(sessionID string) error {
	m.mu.Lock(); if r,ok:=m.running[sessionID];ok && (r.status.Status=="running"||r.status.Status=="queued"){m.mu.Unlock();return fmt.Errorf("session already running")}
	ctx,cancel:=context.WithCancel(context.Background()); r:=&running{status:core.TaskStatus{SessionID:sessionID,Status:"running"},cancel:cancel};m.running[sessionID]=r;m.mu.Unlock()
	go m.run(ctx,sessionID)
	return nil
}

func (m *Manager) update(id,status,errText string,count int){m.mu.Lock();defer m.mu.Unlock();if r,ok:=m.running[id];ok{r.status.Status=status;r.status.Error=errText;r.status.ToolCalls=count}}

func (m *Manager) run(ctx context.Context, sessionID string) {
	session,ok:=m.store.GetSession(sessionID);if !ok{m.update(sessionID,"failed","session not found",0);return}
	profile,ok:=m.store.Provider(session.Provider);if !ok||profile.BaseURL==""||session.Model==""{m.replyError(session,"Configure a provider and model in Settings before chatting.");return}
	registry:=tools.Registry{Workspace:session.Workspace}; count:=0
	for count <= m.maxTools {
		select{case <-ctx.Done():m.update(sessionID,"stopped","",count);return;default:}
		history:=m.store.Messages(sessionID); messages:=make([]provider.ChatMessage,0,len(history)+1)
		messages=append(messages,provider.ChatMessage{Role:"system",Content:"You are Arka, a local-first AI agent. Use tools to complete tasks, verify results, stay inside the workspace, and answer concisely. Tool execution is real and local; never claim an action succeeded without a tool result."})
		for _,v:=range history{
			msg:=provider.ChatMessage{Role:v.Role,Content:v.Content,Name:v.Name,ToolCallID:v.ToolCallID}
			if len(v.ToolCalls)>0{_ = json.Unmarshal(v.ToolCalls,&msg.ToolCalls)}
			messages=append(messages,msg)
		}
		res,err:=m.provider.Chat(ctx,profile,session.Model,messages,registry.Definitions());if err!=nil{m.replyError(session,"Provider error: "+err.Error());m.update(sessionID,"failed",err.Error(),count);return}
		content:="";if s,ok:=res.Message.Content.(string);ok{content=s}
		rawCalls,_:=json.Marshal(res.Message.ToolCalls)
		_ = m.store.AddMessage(core.Message{ID:core.NewID(),SessionID:sessionID,Role:"assistant",Content:content,ToolCalls:rawCalls,Provider:profile.ID,Model:session.Model,CreatedAt:time.Now().UTC()})
		if len(res.Message.ToolCalls)==0{m.update(sessionID,"completed","",count);return}
		for _,call:=range res.Message.ToolCalls{
			count++;m.update(sessionID,"running","",count)
			if count>m.maxTools{break}
			result,toolErr:=registry.Execute(ctx,call.Function.Name,call.Function.Arguments);if toolErr!=nil{if result!=""{result += "\n"};result += "ERROR: "+toolErr.Error()}
			if len(result)>1<<20{result=result[:1<<20]+"\n[truncated]"}
			_ = m.store.AddMessage(core.Message{ID:core.NewID(),SessionID:sessionID,Role:"tool",Name:call.Function.Name,ToolCallID:call.ID,Content:result,CreatedAt:time.Now().UTC()})
		}
	}
	m.replyError(session,fmt.Sprintf("Stopped after the safety ceiling of %d tool calls.",m.maxTools));m.update(sessionID,"failed","tool call ceiling reached",m.maxTools)
}

func (m *Manager) replyError(s core.Session,text string){_ = m.store.AddMessage(core.Message{ID:core.NewID(),SessionID:s.ID,Role:"assistant",Content:strings.TrimSpace(text),Provider:s.Provider,Model:s.Model,CreatedAt:time.Now().UTC()});m.update(s.ID,"failed",text,0)}
