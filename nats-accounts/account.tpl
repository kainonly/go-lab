accounts: {
{{- range .}}
    {{ .Name}} {
        jetstream: enabled
        users: [
        {{- range .Users}}
            { nkey: {{ .NKey }} }
        {{- end }}
        ]
    }
{{- end }}
}
