sudo kind delete cluster
sudo ./cluster_w_registry.sh
sudo make local-image
yes | sudo docker builder prune
sudo docker rmi $(sudo docker images -f "dangling=true" -q)
sudo docker push localhost:5001/kube-scheduler
sudo docker push localhost:5001/controller
sudo helm install scheduler-plugins manifests/install/charts/as-a-second-scheduler/
sudo kubectl apply -f network_info_crd.yaml
# sudo kubectl apply -f sample_crd_obj.yaml
# sudo kubectl apply -f crontab_crd.yaml
sudo kubectl apply -f ./config/rbac/role.yaml 
sudo kubectl apply -f ./config/rbac/role_binding.yaml 
sudo kubectl apply -f sample_pod.yaml
sudo kubectl get pods -n scheduler-plugins