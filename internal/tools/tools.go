package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/Dien45/Mobile-agent/internal/provider"
)

type Registry struct { Workspace string }

func (r Registry) Definitions() []provider.ToolDefinition {
	object := func(props map[string]any, required ...string) map[string]any { return map[string]any{"type":"object", "properties":props, "required":required, "additionalProperties":false} }
	str := func(desc string) map[string]any { return map[string]any{"type":"string", "description":desc} }
	return []provider.ToolDefinition{
		{Type:"function", Function:provider.ToolFunction{Name:"list_files", Description:"List files beneath a workspace directory.", Parameters:object(map[string]any{"path":str("Workspace-relative directory")})}},
		{Type:"function", Function:provider.ToolFunction{Name:"read_file", Description:"Read a UTF-8 text file in the local workspace.", Parameters:object(map[string]any{"path":str("Workspace-relative file path")}, "path")}},
		{Type:"function", Function:provider.ToolFunction{Name:"write_file", Description:"Atomically create or overwrite a UTF-8 text file in the local workspace.", Parameters:object(map[string]any{"path":str("Workspace-relative file path"), "content":str("Complete file content")}, "path", "content")}},
		{Type:"function", Function:provider.ToolFunction{Name:"terminal_exec", Description:"Run a command locally in the workspace. Use only when file tools are insufficient.", Parameters:object(map[string]any{"command":str("Shell command"), "timeout_seconds":map[string]any{"type":"integer", "minimum":1, "maximum":600}}, "command")}},
		{Type:"function", Function:provider.ToolFunction{Name:"git_status", Description:"Read local git status for the workspace.", Parameters:object(map[string]any{})}},
		{Type:"function", Function:provider.ToolFunction{Name:"git_diff", Description:"Read the local working tree and staged diff.", Parameters:object(map[string]any{})}},
		{Type:"function", Function:provider.ToolFunction{Name:"git_commit", Description:"Stage explicit workspace paths and create a local Git commit. This never pushes.", Parameters:object(map[string]any{"paths":map[string]any{"type":"array","items":map[string]any{"type":"string"},"minItems":1},"message":str("Commit message")},"paths","message")}}, 
	}
}

func (r Registry) safePath(name string) (string, error) {
	root, err := filepath.Abs(r.Workspace); if err != nil { return "", err }
	if resolved, resolveErr := filepath.EvalSymlinks(root); resolveErr == nil { root = resolved }
	if name == "" { name = "." }
	if filepath.IsAbs(name) { return "", fmt.Errorf("absolute paths are not allowed") }
	candidate, err := filepath.Abs(filepath.Join(root, filepath.Clean(name))); if err != nil { return "", err }
	if err := within(root, candidate); err != nil { return "", err }
	// Existing targets and existing parents are resolved so a workspace symlink cannot escape.
	probe := candidate
	if _, statErr := os.Lstat(probe); statErr != nil { probe = filepath.Dir(probe) }
	if resolved, resolveErr := filepath.EvalSymlinks(probe); resolveErr == nil {
		if err := within(root, resolved); err != nil { return "", fmt.Errorf("symlink %w", err) }
	}
	return candidate, nil
}

func within(root, candidate string) error {
	rel, err := filepath.Rel(root, candidate); if err != nil { return err }
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) { return fmt.Errorf("path escapes workspace") }
	return nil
}

func (r Registry) Execute(ctx context.Context, name, raw string) (string, error) {
	var args map[string]any
	if raw == "" { raw = "{}" }
	if err := json.Unmarshal([]byte(raw), &args); err != nil { return "", fmt.Errorf("invalid tool arguments: %w", err) }
	switch name {
	case "list_files":
		p, err := r.safePath(text(args,"path")); if err != nil { return "", err }
		entries, err := os.ReadDir(p); if err != nil { return "", err }
		lines := make([]string, 0, len(entries)); for _, e := range entries { suffix := ""; if e.IsDir() { suffix = "/" }; lines = append(lines, e.Name()+suffix) }
		sort.Strings(lines); return strings.Join(lines, "\n"), nil
	case "read_file":
		p, err := r.safePath(text(args,"path")); if err != nil { return "", err }
		b, err := os.ReadFile(p); if err != nil { return "", err }; if len(b) > 1<<20 { return "", fmt.Errorf("file exceeds 1 MiB limit") }; return string(b), nil
	case "write_file":
		p, err := r.safePath(text(args,"path")); if err != nil { return "", err }
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil { return "", err }
		tmp := p+".arka.tmp"; if err := os.WriteFile(tmp, []byte(text(args,"content")), 0o644); err != nil { return "", err }; if err := replaceFile(tmp,p); err != nil { _ = os.Remove(tmp); return "", err }; return "file written", nil
	case "terminal_exec":
		seconds := integer(args,"timeout_seconds",120); if seconds > 600 { seconds = 600 }
		childCtx, cancel := context.WithTimeout(ctx, time.Duration(seconds)*time.Second); defer cancel()
		var cmd *exec.Cmd; if runtime.GOOS == "windows" { cmd = exec.CommandContext(childCtx,"powershell.exe","-NoProfile","-Command",text(args,"command")) } else { cmd = exec.CommandContext(childCtx,"/bin/sh","-lc",text(args,"command")) }
		cmd.Dir = r.Workspace; out, err := cmd.CombinedOutput(); if len(out)>1<<20 { out=out[:1<<20] }
		if childCtx.Err()!=nil { return string(out), fmt.Errorf("command timed out") }; if err != nil { return string(out), fmt.Errorf("command failed: %w",err) }; return string(out),nil
	case "git_status":
		cmd := exec.CommandContext(ctx,"git","status","--short","--branch"); cmd.Dir=r.Workspace; out,err:=cmd.CombinedOutput(); return string(out),err
	case "git_diff":
		working:=exec.CommandContext(ctx,"git","diff","--");working.Dir=r.Workspace;wout,werr:=working.CombinedOutput();if werr!=nil{return string(wout),werr}
		staged:=exec.CommandContext(ctx,"git","diff","--cached","--");staged.Dir=r.Workspace;sout,serr:=staged.CombinedOutput();if serr!=nil{return string(sout),serr}
		return "WORKING TREE\n"+string(wout)+"\nSTAGED\n"+string(sout),nil
	case "git_commit":
		message:=strings.TrimSpace(text(args,"message"));if message==""{return "",fmt.Errorf("commit message is required")}
		rawPaths,ok:=args["paths"].([]any);if !ok||len(rawPaths)==0{return "",fmt.Errorf("at least one explicit path is required")}
		pathArgs:=[]string{"add","--"};for _,rawPath:=range rawPaths{path,ok:=rawPath.(string);if !ok||path==""{return "",fmt.Errorf("invalid commit path")};if _,err:=r.safePath(path);err!=nil{return "",err};pathArgs=append(pathArgs,filepath.Clean(path))}
		add:=exec.CommandContext(ctx,"git",pathArgs...);add.Dir=r.Workspace;if out,err:=add.CombinedOutput();err!=nil{return string(out),fmt.Errorf("git add failed: %w",err)}
		commit:=exec.CommandContext(ctx,"git","commit","-m",message);commit.Dir=r.Workspace;out,err:=commit.CombinedOutput();if err!=nil{return string(out),fmt.Errorf("git commit failed: %w",err)}
		hash:=exec.CommandContext(ctx,"git","rev-parse","HEAD");hash.Dir=r.Workspace;hout,_:=hash.Output();return strings.TrimSpace(string(hout))+"\n"+string(out),nil
	default: return "", fmt.Errorf("unknown tool %q",name)
	}
}

func replaceFile(tmp, target string) error {
	backup:=target+".arka.bak";_ = os.Remove(backup)
	if _,err:=os.Stat(target);err==nil{if err:=os.Rename(target,backup);err!=nil{return err}}
	if err:=os.Rename(tmp,target);err!=nil{_ = os.Rename(backup,target);return err}
	_ = os.Remove(backup);return nil
}
func text(m map[string]any, key string) string { if v,ok:=m[key].(string); ok{return v}; return "" }
func integer(m map[string]any,key string, fallback int) int { if v,ok:=m[key].(float64);ok{return int(v)}; return fallback }
