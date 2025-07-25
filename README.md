![poster](static/poster_03_4_6.png)

# tommybot делает отклики на вакансии за вас, нужно только правильно установить, используя это руководство:
1. для начала, нужно получить "client id" и "client secret", для этого:
   - почитайте условия использования https://dev.hh.ru/admin/developer_agreement
   - если со всем согласны, зайдите на https://dev.hh.ru/, нажмите на кнопку "Добавить приложение"
   - далее нажмите на кнопку "Регистрация нового приложения", откроется форма "Заявка на регистрацию приложения"
   - в поле "Название приложения" заполните какое-нибудь название на английском допустим "tommybot"
   - укажите в поле "Redirect URI" строчку https://tommybot-code.onrender.com/get_code (эту строчку потом можно будет изменить)
   - заполните ваши имя, фамилию, отчество на русском языке
   - в поле "Приложением будут пользоваться:" выделите "Только соискатели"
   - в поле "Информация о создателе приложения:" напишите вашу профессию на русском или английском языке, допустим "повар" или "software developer"
   - в поле "Кто будет его использовать:" напишите "только лично я" или "все кто даст доступ к своему аккаунту"
   - в поле "Какие задачи должно решать приложение:" напишите "ежедневно откликаться на вакансии"
   - в поле "Опишите все функциональные возможности приложения и укажите используемые методы API:" напишите "GET /vacancies; POST /negotiations"
   - нажмите кнопку "Добавить"
   - теперь, ждите несколько дней письма на почту с одобрением вашей заявки, также она будет видна по ссылке https://dev.hh.ru/admin
2. далее развернуть приложение:
   - скачайте zip файл с репозитория tommybot https://github.com/vitalybarkov/tommybot/zipball/master/, или склонируйте https://github.com/vitalybarkov/tommybot
   - распакуйте с названием папки "tommybot" или что-то подобное
   - установите golang на ваш компьютер (https://go.dev/doc/install)
   - используйте терминал (или command prompt для windows):
      - "cd /Users/your_name/Downloads/tommybot" (изменить рабочий каталог на распакованный вами tommybot)
      - "go build main.go"
      - "./main" (или "main.exe" на windows, если вам нужно запустить приложение на другом порте, просто добавьте номер порта, например "./main 8081")
      - откройте браузер и перейдите на "localhost:8080/set" (если вы развернули приложение на удаленном сервере, то тогда "location_of_your_deploy/set")
      - откройте браузер во второй вкладке и перейдите на https://dev.hh.ru/admin и отредактировать поле вашего приложение "Redirect URI" на https://tommybot-code.onrender.com/get_code (если вы развернули приложение на удаленном сервере, используйте "location_of_your_deploy/get_code")
3. и далее подготовить данные для настройки tommybot и заполнить форму настройки tommybot "setting up tommybot":
   - скопируйте "client id" и "client secret" вашего приложения на https://dev.hh.ru/admin, и вставьте на форме настроек в поля "client id" и "client secret"
   - затем кликните на ссылку "GET A CODE"
   - авторизуйтесь под вашим hh.ru аккаунтом
   - подождите (ожидание может затянутся примерно на минуту) пока появится код, скопируйте его и вернитесь на форму настроек и вставьте в поле "code"
   - нажмите на ссылку "GET ACCESS AND REFRESH TOKENS", затем поля "access token" и "refresh token" заполняться автоматически (код может быть использован только единожды, чтобы взять новый код, нажмите заново на ссылку "GET A CODE")
   - подготовьте и заполните поле "search links" (каждая линия это поисковая ссылка, может быть несколько несколько линий) — в вашем браузере используйте поиск на hh.ru, затем скопируйте ссылку сверху в адресной строке браузера, выглядит примерно как https://moscow.hh.ru/search/vacancy?hhtmFrom=main&hhtmFromLabel=vacancy_search_line&search_field=name&search_field=company_name&search_field=description&enable_snippets=false&L_save_area=true&area=1&text=java
   - подготовьте и заполните поле "resume id" — в браузере откройте свое резюме на hh.ru, и скопируйте id резюме сверху в адресной строке браузера после https://moscow.hh.ru/resume/, выглядит примерно как f51d091fff033004f00039ed1f5876516d6e42
   - необязательно: напишите сопроводительное письмо и заполните поле "cover letter" (письмо должно быть до 10_000 символов, и должно быть универсальным для всех откликов по вакансиям (поисковым ссылкам))
> [!IMPORTANT]
> имейте ввиду, вы можете настроить все данные только один раз — для повторной настройки вам нужно будет удалить файл «api_credentials.json» и переименовать файл ./static/_index.html в ./static/index.html (если вы развернули удаленно, то вам потребуется повторно развернуть репозиторий)
4. все, теперь вы можете запустить tommybot, открыв "localhost:8080/hi" ("location_of_your_deploy/hi"), или перезапустить "localhost:8080/restart" ("location_of_your_deploy/restart")
