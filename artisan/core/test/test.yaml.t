---
{{range $key, $value := group "PORT_VALUE" }}
port:
    name: {{ $key }}
    value: {{ $value }}
{{end}}
...