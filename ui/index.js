// fetch('http://localhost:8021/getTaskByID?id=1', {
//     method: 'get'
// }).
// then(response => response.json()).
// then(task =>{
//
//     const tasksUl = document.getElementById('tasks') //взять div и сохранить в переменную
//     console.log(tasksUl)
//     tasksUl.innerHTML =`<li>${task.name}: ${task.status}</li>`
// })

fetch('http://localhost:8021/getAllTasks', {}).then(response => response.json()).then(tasks => {
    const tasksUl = document.getElementById('tasks')
    tasks.forEach(function (tasks) {
        tasksUl.innerHTML += `<li>${tasks.name}: ${tasks.status}</li>`
    })
})

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
    }).
    then(response => {
        console.log(response.status)
        location.reload();
    })
});



