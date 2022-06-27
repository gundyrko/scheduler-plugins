sudo kind delete cluster
sudo ./cluster_w_registry.sh
# sudo docker push localhost:5002/kube-scheduler
# sudo docker push localhost:5002/controller
# sudo docker push localhost:5002/appedge
sudo helm install scheduler-plugins manifests/install/charts/as-a-second-scheduler/
sudo kubectl apply -f network_info_crd.yaml
# sudo kubectl apply -f sample_crd_obj.yaml
# sudo kubectl apply -f crontab_crd.yaml
sudo kubectl apply -f ./config/rbac/role.yaml 
sudo kubectl apply -f ./config/rbac/role_binding.yaml
sudo kubectl apply -f config/rbac/app_proxy_service.yaml
sudo kubectl apply -f sample_pod_deployment.yaml
# sudo kubectl apply -f sample_pod.yaml
sudo kubectl get pods -n scheduler-plugins