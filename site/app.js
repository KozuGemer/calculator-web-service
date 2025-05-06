// Функция для получения значения cookie по имени
function getCookie(name) {
    let value = "; " + document.cookie;
    let parts = value.split("; " + name + "=");
    if (parts.length == 2) return parts.pop().split(";").shift();
}
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

document.getElementById('auth-form').addEventListener('submit', async function(e) {
    e.preventDefault();
    
    // Берем данные из формы
    const requestBody = JSON.stringify({
        login: document.getElementById('login').value,
        password: document.getElementById('password').value
    });
    console.log("Отправляемые данные:", requestBody);
    
    rlogin = document.getElementById('login').value;
    rpasswd = document.getElementById('password').value;
    
    try {
        // Отправляем запрос на авторизацию
        const response = await fetch('http://localhost:8080/api/v1/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',  // Устанавливаем заголовок Content-Type как application/json
            },
            body: JSON.stringify({
                login: rlogin,   // Используем переменную login
                password: rpasswd   // Используем переменную password
            })
        });
        
        console.log("Response status:", response.status);
        
        // Если HTTP-статус не ок, пытаемся получить ошибку из JSON
        if (!response.ok) {
            const errorData = await response.text(); // Получаем текст ошибки
            console.error("Ошибка: ", errorData);
            return;
        }
        
        // Парсим данные ответа
        const data = await response.json();
        console.log("Ответ от сервера:", data);
        
        // Проверяем, что ответ от сервера содержит поле "token"
        if (data.message === "Login successful") {
            // Сохраняем токен в cookie и localStorage
            // Здесь мы можем сохранить токен, если он возвращается сервером
            document.cookie = `token=${data.token}; path=/`;
            localStorage.setItem('token', data.token);
            console.log("Токен сохранен:", data.token);
            // Перенаправляем пользователя на страницу калькулятора
            window.location.href = "index.html";
        } else {
            alert('Ошибка авторизации');
        }
    } catch (error) {
        console.error("Ошибка:", error);
        alert('Ошибка: ' + error);
    }
});
