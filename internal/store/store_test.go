package store

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/Dien45/Mobile-agent/internal/core"
)

func TestSessionLifecyclePersists(t *testing.T) {
	path:=filepath.Join(t.TempDir(),"state.json")
	s,err:=Open(path);if err!=nil{t.Fatal(err)}
	now:=time.Now().UTC();session:=core.Session{ID:"one",Title:"Original",CreatedAt:now,UpdatedAt:now}
	if err:=s.SaveSession(session);err!=nil{t.Fatal(err)}
	if err:=s.RenameSession("one","Renamed");err!=nil{t.Fatal(err)}
	if err:=s.AddMessage(core.Message{ID:"m1",SessionID:"one",Role:"user",Content:"hello",CreatedAt:now});err!=nil{t.Fatal(err)}
	if err:=s.DeleteSession("one",false);err!=nil{t.Fatal(err)}
	if len(s.ListSessions(false))!=0{t.Fatal("deleted session should be hidden")}
	if err:=s.RestoreSession("one");err!=nil{t.Fatal(err)}
	reopened,err:=Open(path);if err!=nil{t.Fatal(err)}
	got,ok:=reopened.GetSession("one");if !ok||got.Title!="Renamed"{t.Fatalf("unexpected session: %#v",got)}
	if len(reopened.Messages("one"))!=1{t.Fatal("message did not persist")}
}
