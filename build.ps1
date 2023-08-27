param([String]$randomHash="JWHQZM9Z4HQOYICDHW4OCJAXPPNHBA")
garble build -ldflags="-X main.randomHash=$randomHash -H=windowsgui";