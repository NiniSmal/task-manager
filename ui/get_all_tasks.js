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

let button = document.getElementById('buttonGetAll');
button.addEventListener('click', () => {
    fetch('http://localhost:8021/getAllTasks', {}).then(response => response.json()).then(tasks => {
        const tasksUl = document.getElementById('tasks')
        tasks.forEach(function (tasks) {
            tasksUl.innerHTML += `<li>${tasks.name}: ${tasks.status}</li>`
        })
    })

})



