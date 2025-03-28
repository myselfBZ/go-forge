package project

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)




func (p *Project) Build() error {
    dirName := getProjectName(p.Name)

    if err := os.Mkdir(dirName, 0755); err != nil {
        return err
    }
    
    os.Chdir(dirName)
    
    cmd := exec.Command("go", "mod", "init", p.Name)

    if _, err := cmd.CombinedOutput(); err != nil{
        return err
    }
    
    cwd, err := os.Getwd()

    if err != nil{
        return err
    }

    if err := p.traverseDirStructure(cwd, p.buildProjectStructure()); err != nil{
        return err
    }

    cmd = exec.Command("go", "mod", "tidy")

    if _, err := cmd.CombinedOutput(); err != nil{
        return err
    }

    return nil
}

func (p *Project) buildProjectStructure() Dir {
    frameworkPath := fmt.Sprintf("internal/src-files/%s/api/", p.Framework)
    dbPath := fmt.Sprintf("internal/src-files/%s/", p.DB)

    var root Dir
    root.Name = "root"
    var internals Dir
    internals.Name = "internal"
    var cmd Dir
    cmd.Name = "cmd"

    cmd.Dirs = append(cmd.Dirs, Dir{
        Name: "api",
        Files: []File{
            {Name: "api.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "api.go.boil"))},
            {Name: "api_test.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "api_test.go.boil"))},
            {Name: "auth.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "auth.go.boil"))},
            {Name: "errors.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "errors.go.boil"))},
            {Name: "users.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "users.go.boil"))},
            {Name: "users_test.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "users_test.go.boil"))},
            {Name: "health.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "health.go.boil"))},
            {Name: "json.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "json.go.boil"))},
            {Name: "middleware.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "middleware.go.boil"))},
            {Name: "test_utils.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "test_utils.go.boil"))},
            {Name: "main.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "main.go.boil"))},
        },
    })

    cmd.Dirs = append(cmd.Dirs, Dir{
        Name: "migrate",
    })

    internals.Dirs = append(internals.Dirs, Dir{
        Name: "auth", 
        Files: []File{
            {Name: "auth.go", Src: p.LoadFile("internal/src-files/auth/auth.go.boil")},
            {Name: "jwt.go", Src: p.LoadFile("internal/src-files/auth/jwt.go.boil")},
            {Name:"mocks.go", Src: p.LoadFile("internal/src-files/auth/mocks.go.boil")},
        },
    })


    store := Dir{
        Name: "store",
        Files: []File{
            {Name: "storage.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "storage.go.boil") )},
            {Name: "users.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "users.go.boil"))},
            {Name: "mocks.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "mocks.go.boil") )},
        },
    }

    store.Dirs = append(store.Dirs, Dir{
        Name: "cache",
        Files: []File{
            {Name: "storage.go", Src: p.LoadFile("internal/src-files/cache/storage.go.boil")},
            {Name: "mocks.go", Src: p.LoadFile("internal/src-files/cache/mocks.go.boil")},
            {Name: "users.go", Src: p.LoadFile("internal/src-files/cache/users.go.boil")},
            {Name: "redis.go", Src: p.LoadFile("internal/src-files/cache/redis.go.boil")},
        },
    })

    internals.Dirs = append(internals.Dirs, store)

    internals.Dirs = append(internals.Dirs, Dir{
        Name: "env",
        Files: []File{
            {Name: "env.go", Src: p.LoadFile("internal/src-files/env/env.go.boil")},
        },
    })


    internals.Dirs = append(internals.Dirs, Dir{
        Name: "db",
        Files: []File{
            {Name: "db.go", Src: p.LoadFile("internal/src-files/db/db.go.boil")},
        },
    })

    root.Dirs = append(root.Dirs, cmd)
    root.Dirs = append(root.Dirs, internals)

    root.Files = []File{
        {Name: "Dockerfile", Src: p.LoadFile("internal/src-files/root/Dockerfile")},
        {Name: ".gitignore", Src: p.LoadFile("internal/src-files/root/.gitignore")},
        {Name: "compose.yaml", Src: p.LoadFile("internal/src-files/root/compose.yaml")},
        {Name: "Makefile", Src: p.LoadFile("internal/src-files/root/Makefile")},
        {Name: "README.md", Src: p.LoadFile("internal/src-files/root/README.md")},
    }

    return root
}

func (p *Project) traverseDirStructure(path string, dir Dir) error {

    if dir.Name == "root" {
        dir.Name = "."
    } else {
        if err := os.Mkdir(path + "/" + dir.Name, 0755); err != nil{
            return err
        }
    }

    for _, f := range dir.Files{
        if err := os.WriteFile(fmt.Sprintf("%s/%s/%s", path, dir.Name, f.Name), f.Src, 0644); err != nil{
            return err
        }
    }

    for _, d := range dir.Dirs{
        if err := p.traverseDirStructure(path+"/"+dir.Name, d); err != nil{
            return err
        }
    }

    return nil
}

func (p *Project) LoadFile(path string) []byte {
    data, err := p.Fs.Open(path)
    if err != nil{
        log.Fatal(err)
    }

    buff := make([]byte, 4096)
    n, err := data.Read(buff)

    if err != nil{
        if !errors.Is(err, io.EOF) {
            log.Fatal(err)
        }
    }

    return buff[:n]
}

func getProjectName(name string) string {
    if strings.HasPrefix(name, "github.com/", ){
        chunks := strings.Split(name, "/")
        return chunks[len(chunks) - 1]
    }
    return name
}
