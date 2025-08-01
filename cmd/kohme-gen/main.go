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
	"fmt"
	"github.com/kohmebot/pkg/command"
	"github.com/kohmebot/pkg/version"
	"github.com/kohmebot/plugin"
	zero "github.com/wdvxdr1123/ZeroBot"
)

type <PluginStruct> struct {
	env plugin.Env
}

func NewPlugin() plugin.Plugin {
	return new(<PluginStruct>)
}

func (p *<PluginStruct>) Init(engine *zero.Engine, env plugin.Env) error {
	p.env = env

	return nil
}

func (p *<PluginStruct>) Name() string {
	return "<PluginName>"
}

func (p *<PluginStruct>) Description() string {
	return "插件描述"
}

func (p *<PluginStruct>) Commands() fmt.Stringer {
	return command.NewCommands()
}

func (p *<PluginStruct>) Version() uint64 {
	return uint64(version.NewVersion(0, 0, 10))
}

func (p *<PluginStruct>) OnBoot() {

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
