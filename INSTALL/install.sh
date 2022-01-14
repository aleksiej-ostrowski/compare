# sudo snap install --classic go

# or

ver="1.17.5"
curl -O https://dl.google.com/go/go{$ver}.linux-amd64.tar.gz

sudo rm -rf ./go

tar xvf go$ver.linux-amd64.tar.gz

sudo rm -rf /usr/local/go
sudo mv ./go /usr/local

# -- sudo chown -R root:root ./go

go env -w GO111MODULE=auto

export GOROOT=/usr/local/go 

# export PATH=$GOPATH/bin:$GOROOT/bin:$PATH 

# sudo apt-get install sqlite3
# sudo apt-get install sqlitebrowser
