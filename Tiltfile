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

k8s_resource('todo-api', port_forwards="9191", labels=["todo-api"])
k8s_resource('todo-frontend', port_forwards="9090", labels=["todo-frontend"])
k8s_resource('nats', port_forwards=["4222", "8222"], labels=["nats"])

docker_build(
    "todo-api",
    context=".",
    dockerfile="cmd/todo-api/Dockerfile",
    ignore=['./data/', './frontend/'],
)

docker_build(
    "todo-frontend",
    context="./frontend",
    dockerfile="frontend/Dockerfile",
    ignore=['./data/'],
)



