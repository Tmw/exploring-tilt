# load('ext://restart_process', 'docker_build_with_restart')
allow_k8s_contexts('orbstack')

def replacements(blob):
    cwd = os.getcwd()
    return blob.replace('__REPO_PATH__', cwd)

yaml = str(read_file('./k8s/todo-api.yml'))
yaml = replacements(yaml)

k8s_yaml(blob(yaml))

k8s_yaml(['k8s/todo-frontend.yml'])

k8s_resource('todo-api', port_forwards="9191", labels=["todo-api"])
k8s_resource('todo-frontend', port_forwards="9090", labels=["todo-frontend"])

docker_build(
    "todo-api",
    context=".",
    dockerfile="cmd/todo-api/Dockerfile",
    ignore=['./data/'],

    #live_update=[
        #sync("cmd/todo-api/", "/app"),
        #sync("internal/todo-service/", "/app"),
        #sync("pkg", "/app"),
        #run("go run cmd/todo-api/.")
    #]
)

docker_build(
    "todo-frontend",
    context="./frontend",
    dockerfile="frontend/Dockerfile",
    ignore=['./data/'],

    #live_update=[
        #sync("cmd/todo-api/", "/app"),
        #sync("internal/todo-service/", "/app"),
        #sync("pkg", "/app"),
        #run("go run cmd/todo-api/.")
    #]
)



