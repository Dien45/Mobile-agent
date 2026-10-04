package skills

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/Dien45/Mobile-agent/internal/core"
)

var validName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,79}$`)

type Manager struct{ Root string }

func (m Manager) List() ([]core.Skill,error){entries,err:=os.ReadDir(m.Root);if os.IsNotExist(err){return []core.Skill{},nil};if err!=nil{return nil,err};out:=[]core.Skill{};for _,e:=range entries{if !e.IsDir()||strings.HasPrefix(e.Name(),"."){continue};desc:=description(filepath.Join(m.Root,e.Name(),"SKILL.md"));info,_:=e.Info();out=append(out,core.Skill{Name:e.Name(),Description:desc,Source:"local",Enabled:true,InstalledAt:info.ModTime()})};sort.Slice(out,func(i,j int)bool{return out[i].Name<out[j].Name});return out,nil}

func (m Manager) Install(ctx context.Context,source,name string)(core.Skill,error){if err:=os.MkdirAll(m.Root,0o700);err!=nil{return core.Skill{},err};stage,err:=os.MkdirTemp(m.Root,".install-");if err!=nil{return core.Skill{},err};defer os.RemoveAll(stage)
	payload:=filepath.Join(stage,"payload")
	switch{
	case strings.HasPrefix(source,"https://")||strings.HasPrefix(source,"ssh://")||strings.HasPrefix(source,"git@"):
		cmd:=exec.CommandContext(ctx,"git","clone","--depth","1","--",source,payload);if out,err:=cmd.CombinedOutput();err!=nil{return core.Skill{},fmt.Errorf("git clone: %w: %s",err,string(out))}
	case strings.HasSuffix(strings.ToLower(source),".zip"):
		if err:=extractZip(source,payload);err!=nil{return core.Skill{},err}
	default:
		if err:=copyTree(source,payload);err!=nil{return core.Skill{},err}
	}
	if name==""{name=filepath.Base(strings.TrimSuffix(source,".git"))};if !validName.MatchString(name){return core.Skill{},fmt.Errorf("invalid skill name")}
	if _,err:=os.Stat(filepath.Join(payload,"SKILL.md"));err!=nil{return core.Skill{},fmt.Errorf("SKILL.md is required at package root")}
	target:=filepath.Join(m.Root,name);if _,err:=os.Stat(target);err==nil{return core.Skill{},fmt.Errorf("skill already exists")}
	if err:=os.Rename(payload,target);err!=nil{return core.Skill{},err}
	return core.Skill{Name:name,Description:description(filepath.Join(target,"SKILL.md")),Source:source,Enabled:true,InstalledAt:time.Now().UTC()},nil
}

func (m Manager) Remove(name string)error{if !validName.MatchString(name){return fmt.Errorf("invalid skill name")};return os.RemoveAll(filepath.Join(m.Root,name))}
func description(path string)string{f,err:=os.Open(path);if err!=nil{return ""};defer f.Close();scan:=bufio.NewScanner(f);for scan.Scan(){line:=strings.TrimSpace(scan.Text());if line==""||strings.HasPrefix(line,"#")||strings.HasPrefix(line,"---"){continue};return line};return ""}

func copyTree(src,dst string)error{root,err:=filepath.Abs(src);if err!=nil{return err};info,err:=os.Stat(root);if err!=nil{return err};if !info.IsDir(){return fmt.Errorf("skill source must be a directory or zip")};return filepath.Walk(root,func(path string,info os.FileInfo,err error)error{if err!=nil{return err};rel,err:=filepath.Rel(root,path);if err!=nil{return err};if rel==".git"||strings.HasPrefix(rel,".git"+string(filepath.Separator)){if info.IsDir(){return filepath.SkipDir};return nil};target:=filepath.Join(dst,rel);if info.Mode()&os.ModeSymlink!=0{return fmt.Errorf("symlinks are not allowed in skill packages")};if info.IsDir(){return os.MkdirAll(target,0o755)};if info.Size()>16<<20{return fmt.Errorf("skill file too large: %s",rel)};in,err:=os.Open(path);if err!=nil{return err};defer in.Close();out,err:=os.OpenFile(target,os.O_CREATE|os.O_WRONLY|os.O_TRUNC,info.Mode().Perm());if err!=nil{return err};_,copyErr:=io.Copy(out,io.LimitReader(in,16<<20+1));closeErr:=out.Close();if copyErr!=nil{return copyErr};return closeErr})}
func extractZip(path,dst string)error{z,err:=zip.OpenReader(path);if err!=nil{return err};defer z.Close();if len(z.File)>2000{return fmt.Errorf("archive has too many files")};var total uint64;for _,f:=range z.File{total+=f.UncompressedSize64;if total>128<<20{return fmt.Errorf("archive exceeds 128 MiB")};clean:=filepath.Clean(f.Name);if filepath.IsAbs(clean)||clean==".."||strings.HasPrefix(clean,".."+string(filepath.Separator)){return fmt.Errorf("archive path escapes package")};target:=filepath.Join(dst,clean);if f.FileInfo().Mode()&os.ModeSymlink!=0{return fmt.Errorf("archive symlinks are not allowed")};if f.FileInfo().IsDir(){if err:=os.MkdirAll(target,0o755);err!=nil{return err};continue};if err:=os.MkdirAll(filepath.Dir(target),0o755);err!=nil{return err};in,err:=f.Open();if err!=nil{return err};out,err:=os.OpenFile(target,os.O_CREATE|os.O_WRONLY|os.O_TRUNC,f.Mode().Perm());if err!=nil{in.Close();return err};_,e1:=io.Copy(out,io.LimitReader(in,16<<20+1));e2:=out.Close();e3:=in.Close();if e1!=nil{return e1};if e2!=nil{return e2};if e3!=nil{return e3}};return nil}
