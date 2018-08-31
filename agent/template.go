package agent

import (
	"io"
	"text/template"
)

const hacfgTmpl = `
global
  daemon
  maxconn 2048

defaults
  mode tcp
  balance roundrobin
  timeout connect 5000ms
  timeout client 50000ms
  timeout server 50000ms

listen http-in
  stick-table type ip size 200k expire 30m
  stick on src
  bind *:80
  {{range .}}server {{.Name}} {{.IP}}:{{.Port}}
  {{end}}
`

//Endpoint from cloud-provider
type Endpoint struct {
	Name string
	IP   string
	Port int
}

//WriteHaproxyCfg output haproxy.cfg
func WriteHaproxyCfg(wr io.Writer, eps []Endpoint) {
	tmpl := template.Must(template.New("t1").Parse(hacfgTmpl))
	tmpl.Execute(wr, eps)
}
