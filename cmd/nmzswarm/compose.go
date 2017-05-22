package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TODO: more template variables
const composeTemplate = `# generated by nmzswarm
version: "3"

services:
  master:
    image: "{{.Image}}"
    entrypoint: /.nmzswarm-agent.master
    command: ["-worker-service=worker", "-input=/.nmzswarm-agent.master.txt", "-chunks={{.Chunks}}", "-shuffle=true", "-rand-seed={{.RandSeed}}"]
    networks:
      - net
    deploy:
      mode: replicated
      replicas: 1
      restart_policy:
        condition: none
      placement:
        constraints: [node.id == {{.SelfNodeID}}]

  worker:
    image: "{{.Image}}"
    entrypoint: /.nmzswarm-agent.worker
    command: ["-executor={{.WorkerScript}}"]
    networks:
      - net
    deploy:
      mode: replicated
      replicas: {{.Replicas}}
      restart_policy:
# The restart condition needs to be any for funker function
        condition: any

networks:
  net:
`

type composeOptions struct {
	Replicas     int
	Chunks       int
	Image        string
	WorkerScript string
	SelfNodeID   string
	RandSeed     int64
}

// createCompose creates "dir/docker-compose.yml".
func createCompose(dir string, opts composeOptions) (string, error) {
	if opts.Replicas == 0 || opts.Chunks == 0 || opts.Image == "" || opts.WorkerScript == "" || opts.SelfNodeID == "" {
		return "", fmt.Errorf("invalid opts: %+v", opts)
	}
	composeFilePath := filepath.Join(dir, "docker-compose.yml")
	tmpl, err := template.New("").Parse(composeTemplate)
	if err != nil {
		return "", err
	}
	f, err := os.Create(composeFilePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	if err = tmpl.Execute(f, opts); err != nil {
		return "", err
	}
	return composeFilePath, nil
}