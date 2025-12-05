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
	name   string
	repo   string
	branch string
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

const ciTemplate = `
name: KohmeBot Plugin CI

on:
  push:
    branches: [ "<branch>" ]
  pull_request:
    branches: [ "<branch>" ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: 18

      - name: Install semver
        run: npm install semver

      - name: Get Plugin Version
        id: version
        shell: bash
        run: |
          MODULE=$(grep '^module ' go.mod | awk '{print $2}')
          PKGNAME=${MODULE##*/}
          FULL_IMPORT="$MODULE/$PKGNAME"

          cat <<EOF > version.go
          package main

          import (
              "fmt"
              "$FULL_IMPORT"
          )

          func main() {
              p := $PKGNAME.NewPlugin()
              fmt.Println(p.Version())
          }
          EOF

          go mod tidy
          go build -o get-version version.go
          VERSION=$(./get-version)
          echo "version=$VERSION" >> "$GITHUB_OUTPUT"
          echo "version=$VERSION"
          git checkout .

      - name: Check version increment
        uses: actions/github-script@v6
        id: check_version
        with:
          script: |
            const semver = require('semver');
            const newVersion = '${{ steps.version.outputs.version }}';

            const releases = await github.rest.repos.listReleases({
              owner: context.repo.owner,
              repo: context.repo.repo,
              per_page: 1,
            });

            if (releases.data.length === 0) {
              console.log('No previous release found, allow.');
              return;
            }

            const latestVersion = releases.data[0].tag_name;

            console.log("Latest release: ${latestVersion}");
            console.log("New version: ${newVersion}");

            if (!semver.valid(newVersion)) {
              core.setFailed("New version (${newVersion}) is not a valid semver.");
              return;
            }

            if (!semver.valid(latestVersion)) {
              console.log("Latest release version (${latestVersion}) is invalid semver, skipping check.");
              return;
            }

            if (semver.lte(newVersion, latestVersion)) {
              core.setFailed("New version (${newVersion}) is not greater than latest release version (${latestVersion})");
            }

      - name: Release
        if: success()  # 只有校验成功才执行
        uses: softprops/action-gh-release@v1
        with:
          tag_name: ${{ steps.version.outputs.version }}
          name: ${{ steps.version.outputs.version }}
          body: ${{ github.event.head_commit.message }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}


`

func init() {
	flag.StringVar(&name, "n", "", "plugin name")
	flag.StringVar(&repo, "r", "", "plugin repo")
	flag.StringVar(&branch, "b", "dev", "git default branch")
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

	if err := os.MkdirAll(path.Join(name, ".github", "workflows"), os.ModePerm); err != nil {
		return err
	}

	if err := runCMDDir(name, "go", "mod", "init", repo); err != nil {
		return err
	}

	structName := fmt.Sprintf("Plugin%s", strings.ToUpper(name[:1])+name[1:])
	f := strings.ReplaceAll(template, "<PluginStruct>", structName)
	f = strings.ReplaceAll(f, "<PluginName>", name)

	ci := strings.ReplaceAll(ciTemplate, "<branch>", branch)

	if err := writeFile(path.Join(name, ".github", "workflows", "kohme-plugin-ci.yml"), []byte(ci)); err != nil {
		return err
	}

	if err := writeFile(path.Join(name, name, "plugin.go"), []byte(f)); err != nil {
		return err
	}

	if err := runCMDDir(name, "go", "mod", "tidy"); err != nil {
		return err
	}

	if err := runCMDDir(name, "git", "init"); err != nil {
		return err
	}

	if err := runCMDDir(name, "git", "checkout", "-b", branch); err != nil {
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
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runCMDDir(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = path.Join(cmd.Dir, dir)
	return cmd.Run()
}
