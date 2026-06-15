.PHONY: shell-frontend shell-api dev
shell-frontend:
	kubectl exec -it deployments/todo-frontend -- sh

shell-api:
	kubectl exec -it deployments/todo-api -- sh

dev:
	tilt up
