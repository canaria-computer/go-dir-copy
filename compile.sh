mv dist dist_old
rm dist_old

GOOS=darwin    GOARCH=386   go build -o ./dist/darwin/386/
GOOS=darwin    GOARCH=amd64 go build -o ./dist/darwin/amd64/
GOOS=freebsd   GOARCH=386   go build -o ./dist/freebsd/386/
GOOS=freebsd   GOARCH=amd64 go build -o ./dist/freebsd/amd64/
GOOS=freebsd   GOARCH=arm   go build -o ./dist/freebsd/arm/
GOOS=linux     GOARCH=386   go build -o ./dist/linux/386/
GOOS=linux     GOARCH=amd64 go build -o ./dist/linux/amd64/
GOOS=linux     GOARCH=arm   go build -o ./dist/linux/arm/
GOOS=netbsd    GOARCH=386   go build -o ./dist/netbsd/386/
GOOS=netbsd    GOARCH=amd64 go build -o ./dist/netbsd/amd64/
GOOS=netbsd    GOARCH=arm   go build -o ./dist/netbsd/arm/
GOOS=openbsd   GOARCH=386   go build -o ./dist/openbsd/386/
GOOS=openbsd   GOARCH=amd64 go build -o ./dist/openbsd/amd64/
GOOS=windows   GOARCH=386   go build -o ./dist/windows/386/
GOOS=windows   GOARCH=amd64 go build -o ./dist/windows/amd64/

7z a -mx=5 ./dist/darwin.zip ./dist/darwin/
7z a -mx=5 ./dist/darwin.7z ./dist/darwin/
7z a -mx=5 ./dist/freebsd.zip ./dist/freebsd/
7z a -mx=5 ./dist/freebsd.7z ./dist/freebsd/
7z a -mx=5 ./dist/linux.zip ./dist/linux/
7z a -mx=5 ./dist/linux.7z ./dist/linux/
7z a -mx=5 ./dist/netbsd.zip ./dist/netbsd/
7z a -mx=5 ./dist/netbsd.7z ./dist/netbsd/
7z a -mx=5 ./dist/openbsd.zip ./dist/openbsd/
7z a -mx=5 ./dist/openbsd.7z ./dist/openbsd/
7z a -mx=5 ./dist/windows.zip ./dist/windows/
7z a -mx=5 ./dist/windows.7z ./dist/windows/