async function calculate() {
    const expression = document.getElementById('expression').value;
    try {
        const response = await fetch('/api/v1/tasks', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ "expression": expression }),
        });
        const data = await response.json();
        console.log("Task status check:", data);
        if (data.id) {
            await checkStatus(data.id);
        } else {
            document.getElementById('result').textContent = 'Ошибка: Некорректный ответ сервера';
        }
    } catch (error) {
        document.getElementById('result').textContent = 'Ошибка: ' + error;
    }
}

async function checkStatus(taskId, delay = 500) {
    try {
        const response = await fetch(`/api/v1/tasks/complete?id=${taskId}`);
        const task = await response.json();
        if (task.status === 'done') {
            document.getElementById('result').textContent = 'Результат: ' + task.result;
        } else {
            await new Promise(resolve => setTimeout(resolve, delay)); // Ждем delay мс
            await checkStatus(taskId, Math.min(delay * 3, 3000)); // Рекурсивно вызываем с увеличенной задержкой
        }
    } catch (error) {
        document.getElementById('result').textContent = 'Ошибка: ' + error;
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
