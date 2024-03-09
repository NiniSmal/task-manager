let button = document.getElementById('buttonCreate');
button.addEventListener('click', () => {
    let task = document.getElementById('nameTask').value
    let body = {
        name: task
    }
    let bodyJson = JSON.stringify(body)
    fetch('http://localhost:8021/createTask', {
        method: "POST",
        body: bodyJson,
    }).then(response => {
        console.log(response.status)
        location.reload();
    })
});



