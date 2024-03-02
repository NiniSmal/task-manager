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

fetch('http://localhost:8021/getAllTasks',{}).
then(response => response.json()).
then(tasks => {
    const tasksUl = document.getElementById('tasks')
        tasks.forEach(function (tasks){
        tasksUl.innerHTML +=`<li>${tasks.name}: ${tasks.status}</li>`
    })
})
// fetch('http://localhost:8021//createTask', {}).
// then(response =>response.json()).
// then(creatTask =>{
//     let inputElement = document.getElementById('creatTask');//получаем введенный элемент
//     var value = inputElement.value;//получаем значение ,введенное пользователем
// })
// let inputElement = document.getElementById('creatTask');//получаем введенный элемент
// var value = inputElement.value;//получаем значение ,введенное пользователем


function createTask(){
    let name =document.getElementById('nameTask').value;
    let params =new URLSearchParams();
    params.set('nameTask', name);

    fetch('http://localhost:8021/createTask',{
        method:'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body:params
    }).then(
        response =>{
            return response.json();
        }
    ).then(
        text =>{
            document.getElementById('result').innerHTML =text;
        }
    )
}