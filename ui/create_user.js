let buttonCreateUser = document.getElementById("buttonCreateUser")
buttonCreateUser.addEventListener('click', () => {
    let login = document.getElementById('create_login',).value
    let password = document.getElementById('create_password',).value
    let body = {
        login: login,
        password: password
    }
    let bodyJson = JSON.stringify(body)
    fetch('http://localhost:8021/createUser', {
        method: "POST",
        body: bodyJson,
    }).then(response => {
        console.log(response.status)
        alert("Регистрация прошла успешно")
        window.location.href ="http://localhost:63342/task-manager-gitlab/ui/auth.html"
    })

})