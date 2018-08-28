package sshtmpl

var TokConfTmplStr = `Host {{.ProjectName}}.tok
    HostName localhost
    Port {{.DrushPort}}
    User tok
    IdentityFile ~/.ssh/tok_ssh.key
    IdentitiesOnly yes
    ForwardAgent yes
    StrictHostKeyChecking no
    UserKnownHostsFile /dev/null

`
