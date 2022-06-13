# Parca-profiling's Makefile command list

########################################################################################################################

build-parca-profiling:
	docker image build --no-cache -t pietroski/parca-profiling-test -f deploy-sample/docker/parca-profiling.Dockerfile ../../../../

tagging-parca-profiling-image:
	docker tag pietroski/parca-profiling-test pietroski/parca-profiling-test:v0.0.4

update-docker-parca-profiling-image:
	docker push pietroski/parca-profiling-test:v0.0.4

kube-parca-profiling-deployment:
	kubectl-24 apply -f deploy-sample/k8s/k8s-deployment.yml

kube-parca-profiling-service:
	kubectl-24 apply -f deploy-sample/k8s/k8s-service.yml

kube-parca-agnt:
	kubectl-24 apply -f deploy-sample/parca/k8s-parca-agent.yaml

kube-parca-svr:
	kubectl-24 apply -f deploy-sample/parca/k8s-parca-server.yaml

kube-parca-namespace:
	kubectl-24 create namespace parca

kube-parca-patch: kube-parca-svr kube-parca-agnt kube-parca-profiling-deployment kube-parca-profiling-service

kube-parca-cleanup:
	kubectl delete services parca-hello-world-service
	kubectl delete deployments.apps parca-hello-world

########################################################################################################################
