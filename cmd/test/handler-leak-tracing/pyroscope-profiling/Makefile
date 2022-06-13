# Pyroscope's Makefile command list

########################################################################################################################

build-pyroscope-profiling:
	docker image build --no-cache -t pietroski/pyroscope-profiling-test -f deploy-sample/docker/pyroscope-profiling.Dockerfile ../../../../

tagging-pyroscope-profiling-image:
	docker tag pietroski/pyroscope-profiling-test pietroski/pyroscope-profiling-test:v0.0.6

update-docker-pyroscope-profiling-image:
	docker push pietroski/pyroscope-profiling-test:v0.0.6

kube-pyroscope-profiling-deployment:
	kubectl-24 apply -f deploy-sample/k8s/k8s-deployment.yml

kube-pyroscope-patch: kube-pyroscope-profiling-deployment

kube-pyroscope-cleanup:
	kubectl delete deployments.apps parca-hello-world

########################################################################################################################
