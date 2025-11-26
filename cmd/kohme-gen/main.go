package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	name string
	repo string
)

const template = `
package <PluginName>

import (
	"github.com/kohmebot/plugin/v2"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type <PluginStruct> struct {
	env plugin.Env
}

func NewPlugin() plugin.Plugin {
	return new(<PluginStruct>)
}

func (p *<PluginStruct>) OnInit(engine plugin.Engine, env plugin.Env) error {
	p.env = env

	return nil
}

func (p *<PluginStruct>) OnBoot() {

}

func (p *<PluginStruct>) OnHelp(ctx *zero.Ctx) {
	ctx.Send("Hello World!")
}

func (p *<PluginStruct>) Name() string {
	return "<PluginName>"
}


func (p *<PluginStruct>) Version() string {
	return "v1.0.0"
}



`

func init() {
	flag.StringVar(&name, "n", "", "plugin name")
	flag.StringVar(&repo, "r", "", "plugin repo")
	flag.Parse()
	if len(name) == 0 {
		panic(fmt.Errorf("plugin name is empty"))
	}
	if len(repo) == 0 {
		panic(fmt.Errorf("plugin repo is empty"))
	}
}

func main() {
	if err := gen(); err != nil {
		panic(err)
	}
}

func gen() error {

	if err := os.MkdirAll(path.Join(name, name), os.ModePerm); err != nil {
		return err
	}

	if err := runCMDDir(name, "go", "mod", "init", repo); err != nil {
		return err
	}

	structName := fmt.Sprintf("Plugin%s", strings.ToUpper(name[:1])+name[1:])
	f := strings.ReplaceAll(template, "<PluginStruct>", structName)
	f = strings.ReplaceAll(f, "<PluginName>", name)

	if err := writeFile(path.Join(name, name, "plugin.go"), []byte(f)); err != nil {
		return err
	}

	if err := runCMDDir(name, "go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}

func writeFile(path string, b []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

func runCMD(name string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runCMDDir(dir string, name string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = path.Join(cmd.Dir, dir)
	return cmd.Run()
}
