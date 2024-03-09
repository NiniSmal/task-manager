let button = document.getElementById("buttonAuth")
button.addEventListener('click', () => {
    let login = document.getElementById('auth_login',).value
    let password = document.getElementById('auth_password',).value
    let body = {
        login: login,
        password: password
    }
    let bodyJson = JSON.stringify(body)
    fetch('http://localhost:8021/login', {
        method: "POST",
        body: bodyJson,
    }).then(response => {
        console.log(response.status)
        window.location.href ="http://localhost:63342/task-manager-gitlab/ui/get_all_task.html"

    })
});