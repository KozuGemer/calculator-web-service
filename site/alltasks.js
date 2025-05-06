async function fetchTasks() {
    try {
        // Запрашиваем список выражений с эндпоинта /api/v1/expressions
        const response = await fetch('/api/v1/expressions');  // Путь к API для получения списка задач
        const tasks = await response.json();
        const tableBody = document.querySelector("#taskTable tbody");
        const taskTable = document.querySelector("#taskTable");
        const noTasksMessage = document.querySelector("#noTasksMessage");
        if (tasks.length === 0) {
            // Если заданий нет, скрываем таблицу и показываем сообщение
            taskTable.style.display = 'none';
            noTasksMessage.style.display = 'block';
        } else {
            // Заполняем таблицу данными
            tasks.forEach(task => {
                const row = document.createElement("tr");

                const taskIdCell = document.createElement("td");
                taskIdCell.textContent = task.id;  // Предположим, что поле "id" в объекте задачи
                row.appendChild(taskIdCell);

                const expressionCell = document.createElement("td");
                expressionCell.textContent = task.expression; // Используем поле "expression"
                row.appendChild(expressionCell);

                const resultCell = document.createElement("td");
                resultCell.textContent = task.result;  // Используем поле "result"
                row.appendChild(resultCell);

                // Добавляем строку в таблицу
                tableBody.appendChild(row);
            });
        }
    } catch (error) {
        console.error("Ошибка при загрузке задач:", error);
        // Если заданий нет, скрываем таблицу и показываем сообщение
        const h1Element = document.querySelector('h1');
        const yourbody = document.querySelector('body');
        noTasksMessage.style.fontSize = '3em';
        yourbody.style.justifyContent = 'center';
        taskTable.style.display = 'none';
        h1Element.style.display = 'none';
        noTasksMessage.style.display = 'block';
                
    }
}
document.getElementById('Exit').addEventListener('click', function() {
    // Отправляем запрос на сервер для logout
    fetch('/logout', {
        method: 'POST',
        credentials: 'include'  // Важно! Это обеспечит отправку cookies с запросом
    })
    .then(response => {
        if (response.ok) {
            // Если logout успешен, перенаправляем пользователя на страницу логина
            window.location.href = '/login';  // Перенаправление на страницу логина
        } else {
            // Обрабатываем ошибку, если logout не удался
            console.error("Ошибка при выходе");
        }
    })
    .catch(error => {
        // Ошибка при отправке запроса
        console.error("Ошибка:", error);
    });
});

// Вызов функции fetchTasks, когда страница загрузится
document.addEventListener("DOMContentLoaded", function() {
    fetchTasks();
});