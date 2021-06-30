git pull
go build .
cd client
go build .
systemctl restart goschoolapi
systemctl restart goschoolclient