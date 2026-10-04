package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWorkspaceBoundaryAndFileRoundTrip(t *testing.T) {
	root := t.TempDir()
	r := Registry{Workspace: root}
	if _, err := r.Execute(context.Background(), "write_file", `{"path":"notes/a.txt","content":"hello"}`); err != nil { t.Fatal(err) }
	got, err := r.Execute(context.Background(), "read_file", `{"path":"notes/a.txt"}`)
	if err != nil || got != "hello" { t.Fatalf("got %q, err %v", got, err) }
	if _, err := os.Stat(filepath.Join(root,"notes","a.txt")); err != nil { t.Fatal(err) }
	if _, err := r.Execute(context.Background(), "read_file", `{"path":"../secret"}`); err == nil { t.Fatal("expected workspace escape to fail") }
	outside:=t.TempDir();if err:=os.WriteFile(filepath.Join(outside,"secret"),[]byte("nope"),0o600);err!=nil{t.Fatal(err)}
	if err:=os.Symlink(outside,filepath.Join(root,"escape"));err==nil{if _,err:=r.Execute(context.Background(),"read_file",`{"path":"escape/secret"}`);err==nil{t.Fatal("expected symlink escape to fail")}}
}

func TestUnknownToolFails(t *testing.T) {
	r := Registry{Workspace:t.TempDir()}
	if _,err:=r.Execute(context.Background(),"does_not_exist",`{}`);err==nil{t.Fatal("expected error")}
}
