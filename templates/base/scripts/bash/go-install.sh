#!/bin/bash
if [[ ! -z $PRIVATE_REPO_URL ]] && [[ ! -z "$PRIVATE_REPO_SSH_KEY_BASE64" ]]; then
  echo $PRIVATE_REPO_SSH_KEY_BASE64 | base64 --decode > temp.key && chmod 0600 temp.key
  export GIT_SSH_COMMAND="ssh -i $(pwd)/temp.key"
  export GOPRIVATE=$PRIVATE_REPO_URL
  git config --global url.ssh://git@github.com:.insteadOf https://github.com
fi

go get .
go install golang.org/x/vuln/cmd/govulncheck@latest

if [[ ! -z $PRIVATE_REPO_URL ]] && [[ ! -z "$PRIVATE_REPO_SSH_KEY_BASE64" ]]; then
  rm temp.key
  unset GIT_SSH_COMMAND
  unset GOPRIVATE
  git config --global --remove url.ssh://git@github.com:
fi
