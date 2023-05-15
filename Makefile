run:
	docker build --tag forum .
	docker run --name app -d --rm -p 8081:8081 forum

stop:
	docker stop app
	docker rmi forum