# install git
sudo apt-get update && sudo apt-get install -y git

# set up git
git clone https://github.com/yagikota/danewebperf.git

# install go
curl -LO https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz && rm go1.21.1.linux-amd64.tar.gz
echo "PATH=$PATH:/usr/local/go/bin" >> ~/.profile
source ~/.profile

# install docker
# https://docs.docker.com/engine/install/debian/#install-using-the-repository
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# user docker command without sudo
sudo gpasswd -a $USER docker
sudo systemctl restart docker

# install tree command
sudo apt-get update && sudo apt-get install -y tree

# install zip and unzip command
sudo apt-get update && sudo apt-get install -y zip unzip

# install beautifulsoup4
sudo apt install -y python3-bs4

sudo apt-get update && sudo apt-get install -y htop

echo "Finished installing all dependencies"
echo "Please reboot your machine to use docker command without sudo"
