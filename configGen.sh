apt install golang-cfssl
cat > ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "kubernetes": {
        "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ],
        "expiry": "87600h"
      }
    }
  }
}
EOF
cat > ca-csr.json <<EOF
{
  "CN": "kubernetes",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
EOF
cfssl gencert -initca ca-csr.json | cfssljson -bare ca
kubectl config set-cluster kubernetes \
  --certificate-authority=ca.pem \
  --embed-certs=true \
  --server=$SERVER_IPPORT \
  --kubeconfig=k8sconfig
kubectl config set-credentials client \
  --client-key=ca-key.pem \
  --client-certificate=ca.pem \
  --embed-certs=true \
  --kubeconfig=k8sconfig
kubectl config set-context client-context \
  --cluster=kubernetes \
  --user=client \
  --kubeconfig=k8sconfig
kubectl config use-context client-context --kubeconfig=k8sconfig
mv k8sconfig ~/.kube/config
export KUBECONFIG=~/.kube/config