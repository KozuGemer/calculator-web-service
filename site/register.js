document.getElementById('register-form').addEventListener('submit', async function(e) {
    e.preventDefault();

    // Получаем данные из формы
    const login = document.getElementById('login').value;
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    // Проверяем, что пароли совпадают
    if (password !== confirmPassword) {
        alert("Пароли не совпадают. Пожалуйста, попробуйте снова.");
        return;  // Останавливаем выполнение, чтобы не отправлять форму
    }

    const requestBody = JSON.stringify({
        login: login,
        password: password
    });

    try {
        const response = await fetch('http://localhost:8080/api/v1/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: requestBody
        });

        if (!response.ok) {
            const errorData = await response.json();
            alert(`Ошибка: ${errorData.error}`);
            return;
        }

        const data = await response.json();
        console.log("Registration successful:", data);
        alert("Регистрация прошла успешно! Теперь вы можете войти.");
        window.location.href = "login.html"; // Перенаправление на страницу логина
    } catch (error) {
        console.error("Error:", error);
        alert('Ошибка регистрации!');
    }
});
