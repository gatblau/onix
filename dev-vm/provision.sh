# Assign arguments
export VM_USER=$1
export VM_USER_PUBKEY=$2
export VM_UID=$3
export VM_GID=$4
export VM_NAME=$5

echo ================================================================
echo Updating and upgrading base OS ...
sudo apt update
# DBG tmp disable
#sudo apt upgrade -y

echo ================================================================
echo Adding group/user for ${VM_USER} ...
groupadd -g ${VM_GID} ${VM_USER}
useradd -m -u ${VM_UID} -g ${VM_GID} -s "/bin/bash" ${VM_USER}
echo "${VM_USER} ALL=(ALL:ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/${VM_USER}-override

echo ================================================================
echo Creating SSH key ...
mkdir -p /home/${VM_USER}/.ssh
chown ${VM_USER}:${VM_USER} /home/${VM_USER}/.ssh
chmod 700 /home/${VM_USER}/.ssh
ssh-keygen -q -t rsa -N '' -C "${VM_USER}@${VM_NAME}" -f /home/${VM_USER}/.ssh/id_rsa <<<y 2>&1 >/dev/null

echo ================================================================
echo Adding self and public key to authorised ...
cat /home/${VM_USER}/.ssh/id_rsa.pub >> /home/${VM_USER}/.ssh/authorized_keys
chown ${VM_USER}:${VM_USER} /home/${VM_USER}/.ssh/authorized_keys
chmod 600 /home/${VM_USER}/.ssh/authorized_keys
echo ${VM_USER_PUBKEY} >> /home/${VM_USER}/.ssh/authorized_keys

echo ================================================================
echo Installing additional tools ...
sudo apt install -y git docker.io docker-compose gnupg
sudo usermod -aG docker ${VM_USER}

echo ================================================================
echo Installing Kubernetes and Tekton CLI ...
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
rm kubectl
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 3EFE0E0A2F2F60AA
echo "deb http://ppa.launchpad.net/tektoncd/cli/ubuntu focal main"|sudo tee /etc/apt/sources.list.d/tektoncd-ubuntu-cli.list
sudo apt update && sudo apt install -y tektoncd-cli

echo ================================================================
echo Adding OpenShift CLI ...
curl -O https://mirror.openshift.com/pub/openshift-v4/clients/ocp/stable/openshift-client-linux.tar.gz
sudo tar xvf openshift-client-linux.tar.gz -C /usr/local/bin/ oc
rm openshift-client-linux.tar.gz

echo ================================================================
echo Adding Kustomize ...
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
sudo install -o root -g root -m 0755 kustomize /usr/local/bin/kustomize
rm kustomize

echo ================================================================
echo Adding Helm ...
curl -s "https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3" | bash

echo ================================================================
echo Cleaning up ...
sudo apt -y autoremove

echo ================================================================
echo Finished
echo Use "ssh ${VM_USER}@127.0.0.1 -p 2222" to login normally
echo Use "vagrant ssh" to connect as the standard vagrant user
