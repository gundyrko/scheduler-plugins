apt install make -y
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.0/kind-linux-amd64
chmod +x ./kind
mv ./kind /bin/kind
snap install helm --classic
wget https://github.com/kubeedge/kubeedge/releases/download/v1.10.1/keadm-v1.10.1-linux-amd64.tar.gz
tar -xzvf ./keadm-v1.10.1-linux-amd64.tar.gz
mv ./keadm-v1.10.1-linux-amd64/keadm/keadm /bin/keadm
rm ./keadm-v1.10.1-linux-amd64.tar.gz
rm ./keadm-v1.10.1-linux-amd64 -r