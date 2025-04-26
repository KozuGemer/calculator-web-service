document.getElementById('auth-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    const login = document.getElementById('login').value;
    const password = document.getElementById('password').value;
    
    try {
        const response = await fetch('http://localhost:8080/api/v1/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                login: login,
                password: password
            }),
        });

        if (!response.ok) {
            const error = await response.json();
            console.error('Ошибка: ', error);
            return;
        }
        
        const data = await response.json();
        console.log("Ответ от сервера:", data);  // Логирование ответа сервера

        if (data.token) {
            // Сохраняем токен в localStorage
            localStorage.setItem('token', data.token);
            console.log("Токен сохранен:", data.token); // Логирование сохраненного токена
            // Переход к калькулятору
            window.location.href = "index.html";
        } else {
            alert('Ошибка авторизации');
        }
    } catch (error) {
        alert('Ошибка: ' + error);
    }
});
