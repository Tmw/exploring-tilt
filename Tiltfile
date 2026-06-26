allow_k8s_contexts('orbstack')

def k8s_yaml_with_replacements(path):
    yaml = str(read_file(path))
    yaml = replacements(yaml)
    k8s_yaml(blob(yaml))

def replacements(blob):
    cwd = os.getcwd()
    return blob.replace('__REPO_PATH__', cwd)

k8s_yaml_with_replacements('k8s/todo-api.yml')
k8s_yaml_with_replacements('k8s/todo-frontend.yml')
k8s_yaml_with_replacements('k8s/nats.yml')
k8s_yaml_with_replacements('k8s/nginx.yml')

k8s_resource('todo-api', labels=["todo-api"])
k8s_resource('todo-frontend', labels=["todo-frontend"])
k8s_resource('nats', labels=["nats"])
k8s_resource('nginx', port_forwards=["8888"], labels=["nginx"])

docker_build(
    "todo-api",
    context=".",
    dockerfile="cmd/todo-api/Dockerfile.dev",
    ignore=['./*'],
)

docker_build(
    "todo-frontend",
    context="./frontend",
    dockerfile="frontend/Dockerfile.dev",
    ignore=['./*'],
)



