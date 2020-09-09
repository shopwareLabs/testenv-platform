build:
	docker build -t shopware/testenv:6.3 images/6
	docker build -t shopware/testenv:5.6 images/5

push:
	docker push shopware/testenv:6.3
	docker push shopware/testenv:5.6