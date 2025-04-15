package project

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
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
    frameworkPath := fmt.Sprintf("%s/api/", p.Framework)
    dbPath := fmt.Sprintf("%s/", p.DB)

    var root Dir
    root.Name = "root"
    var internals Dir
    internals.Name = "internal"
    var cmd Dir
    cmd.Name = "cmd"

    cmd.Dirs = append(cmd.Dirs, Dir{
        Name: "api",
        Files: []File{
            // files that contain project's name
            {Name: "api.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "api.go.boil"))},
            {Name: "auth.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "auth.go.boil"))},
            {Name: "middleware.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "middleware.go.boil"))},
            {Name: "main.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "main.go.boil"))},
            {Name: "test_utils.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "test_utils.go.boil"))},
            {Name: "users.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "users.go.boil"))},

            {Name: "api_test.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "api_test.go.boil"))},
            {Name: "errors.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "errors.go.boil"))},
            {Name: "users_test.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "users_test.go.boil"))},
            {Name: "health.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "health.go.boil"))},
            {Name: "json.go", Src: p.LoadFile(fmt.Sprintf(frameworkPath + "%s", "json.go.boil"))},
        },
    })

    cmd.Dirs = append(cmd.Dirs, Dir{
        Name: "migrate",
        Files: []File{
            { Name: "00001_create_roles.down.sql", Src: p.LoadFile("migrations/00001_create_roles.down.sql") },
            { Name: "00001_create_roles.up.sql", Src: p.LoadFile("migrations/00001_create_roles.up.sql") },
            { Name: "00002_create_user.down.sql", Src: p.LoadFile("migrations/00002_create_user.down.sql") },
            { Name: "00002_create_user.up.sql", Src: p.LoadFile("migrations/00002_create_user.up.sql") },
            { Name: "00003_seed_users_with_roles.down.sql", Src: p.LoadFile("migrations/00003_seed_users_with_roles.down.sql") },
            { Name: "00003_seed_users_with_roles.up.sql", Src: p.LoadFile("migrations/00003_seed_users_with_roles.up.sql") },
        },
    })

    internals.Dirs = append(internals.Dirs, Dir{
        Name: "auth", 
        Files: []File{
            {Name: "auth.go", Src: p.LoadFile("auth/auth.go.boil")},
            {Name: "jwt.go", Src: p.LoadFile("auth/jwt.go.boil")},
            {Name:"mocks.go", Src: p.LoadFile("auth/mocks.go.boil")},
        },
    })


    store := Dir{
        Name: "store",
        Files: []File{
            {Name: "storage.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "storage.go.boil") )},
            {Name: "users.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "users.go.boil"))},
            {Name: "mocks.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "mocks.go.boil") )},
            {Name: "roles.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "roles.go.boil"))},
        },
    }

    store.Dirs = append(store.Dirs, Dir{
        Name: "cache",
        Files: []File{
            {Name: "storage.go", Src: p.LoadFile("cache/storage.go.boil")},
            {Name: "users.go", Src: p.LoadFile("cache/users.go.boil")},
            {Name: "mocks.go", Src: p.LoadFile("cache/mocks.go.boil")},

            {Name: "redis.go", Src: p.LoadFile("cache/redis.go.boil")},
        },
    })

    internals.Dirs = append(internals.Dirs, store)

    internals.Dirs = append(internals.Dirs, Dir{
        Name: "env",
        Files: []File{
            {Name: "env.go", Src: p.LoadFile("env/env.go.boil")},
        },
    })


    internals.Dirs = append(internals.Dirs, Dir{
        Name: "db",
        Files: []File{
            {Name: "db.go", Src: p.LoadFile(fmt.Sprintf(dbPath + "%s", "db/db.go.boil"))},
        },
    })

    root.Dirs = append(root.Dirs, cmd)
    root.Dirs = append(root.Dirs, internals)

    root.Files = []File{
        {Name: "Dockerfile", Src: p.LoadFile("root/Dockerfile")},
        {Name: ".gitignore", Src: p.LoadFile("root/.gitignore")},
        {Name: "compose.yaml", Src: p.LoadFile("root/compose.yaml")},
        {Name: "Makefile", Src: p.LoadFile("root/Makefile")},
        {Name: "README.md", Src: p.LoadFile("root/README.md")},
        {Name: ".env", Src: p.LoadFile("root/.env")},
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
        
        outFile, err := os.Create(fmt.Sprintf("%s/%s/%s", path, dir.Name, f.Name))
        if err != nil{
            log.Fatal("couldn't create file", err)
        }
        
        if err := f.Src.Execute(outFile, p.Name); err != nil{
            log.Fatal("couldn't execute the template:", err)
        }

        outFile.Close()
    }

    for _, d := range dir.Dirs{
        if err := p.traverseDirStructure(path+"/"+dir.Name, d); err != nil{
            return err
        }
    }

    return nil
}

func (p *Project) LoadFile(path string) *template.Template {
    content, err := p.Fs.ReadFile(path)
    if err != nil{
        log.Fatal("couldn't open file", err)
    }
    tmpl, err := template.New(path).Parse(string(content))
    if err != nil{
        log.Fatal("couldn't parse the file", err)
    }

    return tmpl
}

func getProjectName(name string) string {
    if strings.HasPrefix(name, "github.com/", ){
        chunks := strings.Split(name, "/")
        return chunks[len(chunks) - 1]
    }
    return name
}
