package grok_patterns 

const (
  // Mcollective_patterns ...
  Mcollective_patterns string = `
# Remember, these can be multi-line events.
MCOLLECTIVE ., \[%{TIMESTAMP_ISO8601:timestamp} #%{POSINT:pid}\]%{SPACE}%{LOGLEVEL:event_level}

MCOLLECTIVEAUDIT %{TIMESTAMP_ISO8601:timestamp}:
`
)
