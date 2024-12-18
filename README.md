# Сервер с запросами на калькулятор
## Общее преставление
Это калькулятор Web-Server, когда вы делайте определленый запросы на локальный сервер, который создает Go, так скажем.
Он вам выдает результат... этот калькулятор умеет работать со скобками... с последотельностью знаком: *; /; +; -; вычесляет большество обычных длинных примером. Сервер и калькулятор был написан на языке Go.

## Как же его запустить?
Для этого вам нужно скачать этот репозиторий как исходный код там будет кнопка Codе, и там написано Download ZIP его и нужно скачать, далее нужно его сделать Unzip. Точнее разархивировать, скорее всего вы уже научины как - это делается, у вас Go обязательно должен быть установлен, если не установлен **[клик](https://go.dev/dl/)**, и вы передете на офицальный сайт Go там его и скачайте для вашей операционной системой.
Далее вам нужно найти приложение "Консоль" или "Терминал" в вашей операционной системе. Удостоверьтесь командой 

    go version

что Go реально установлен или установился. Если показывает что-то типо этого:
`go version go1.23.3 windows/amd64`.
У меня он установлен.
> В моем случае для Windows для x64bit.

Теперь вам нужно зайти в папку `caclator-web-service` примерно так может называться... используете в консоле команду `cd путь/до/caclucator-web-service` замените `путь/до/caclucator-web-service` путем к реальной папке к моему калькулятору на вашем компьютере или ноутбуке.
И после того когда зайдете в папку вам нужно написать команду `go run main.go` и сервер запуститься на таком адресе http://localhost:8080 вы можете нажать на него [здесь](http://localhost:8080), если у вас запущем сервер. Если появится сообщение `Server is running on http://localhost:8080`, то cервер успешно запустился.
# Установка Curl(пропустите, если уже установлен)
## Установка curl(только Windows)
Чтобы нам его установить нам в начале нужно зайти на их офицальный [сайт](https://curl.se/download.html). Нажмите на клавиатуре сочитание Ctrl + F и тогда откроется поисковик и введите туда **Windows**. Пролиствайте до вкладок **Windows** и выбирайте установщик там рядом название будет `the curl project` вы качайте не **установщик,** поэтому будем сами устанавливать... Вам нужен путь короткий поэтому подойдет для наших задач Ваш системный диск в диске С вам нужно создать папку Curl и разарихивировать папку именно туда... возможно версия уже есть более новая... обновлять тоже нужно самому. __*Не большой*__ минусик есть. И должно получится что-то типо этого `C:\Curl\curl-8.11.0_4-win64-mingw`, если с открытой папкой самого curl. И теперь - этот путь нужно добавить Path (Global) чтобы другие пользователи могли использовать... если не хотите добавляете Path (Local) для вашего пользователя. Чтобы добавить передите в Приложение Настройки (Windows 10/11) >> Система >> О системе >> Дополнительные настройки системы >> Переменные среды и тут два разветления внизу - это глобальные(для всех пользователей) настройки, вверху Локальные(только для этого пользователя) настройки. Выбираете для себя или для всех и ищите пункт Path он всегда будет на __*английском языке*__ и нажимайте Изменить, и дальше кнопку Новое или "...что-то в этом духе...". И вставляете этот самый путь, *возможно* он у вас будет другой. Теперь откройте терминал или консоль, если уже он включен перепустите его. Введите команду `curl www.google.com`, если вышло много **текста**, то curl
## Установка curl(Только MacOS и Linux)
Если случаем вы обнаружили что в вашей системе нету curl - это можно исправить, вот несколько основных Дистрибутивов Linux где я могу дать команду на скачивание установку curl:

  **Ubuntu/Debian**:
  `apt install curl`
  
  **OpenSUSE**:
  `zypper install curl`
  
  **Fedora/Red Hat**:
  `dnf install curl`
  
Для других вы найдете сами в интернете на свой дистрибутив, или же вам придется компилировать его, из его исходников. Обязательно для всех установок вы должны быть под root или используйте sudo вместе с командой.
А для **MacOS**
Он у вас по умолчанию установлен... открывайте терминал и пользуйтесь.
## Примеры использования моего калькулятора через команду curl
> **Примеры, пока отсутвуют, пока ещё не было тестов команд, чтобы удостоверится в их правильной команде.**
