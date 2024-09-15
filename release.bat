@ECHO OFF
SET version="v1.3.8"

git tag %version%
git push origin %version%
go list -m github.com/go-per/simpkg@%version%

PAUSE