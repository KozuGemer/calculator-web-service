// Функция для получения значения cookie по имени
function getCookie(name) {
    let value = "; " + document.cookie;
    let parts = value.split("; " + name + "=");
    if (parts.length == 2) return parts.pop().split(";").shift();
}

document.getElementById('submit').addEventListener('click', async function(e) {
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
        const response = await fetch('/api/v1/login', {
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