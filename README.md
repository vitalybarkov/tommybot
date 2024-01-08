# tommybot set up manual
1. at the beginning, you need to get "client id" and "client secret" at https://dev.hh.ru/admin
2. next, deploy:
   - download zip file of the repository https://github.com/vitalybarkov/tommybot/zipball/master/, or clone https://github.com/vitalybarkov/tommybot/
   - unzip with a folder name "tommybot" or something
   - install golang on your machine (https://go.dev/doc/install)
   - use terminal (or command prompt on win)
      - "cd /Users/you_are/Downloads/tommybot" (change to location of your unzipped tommybot)
      - "go build main.go"
      - "./main" (or "main.exe" on win, if you need run your app with another port, just add your port number like "./main 8081")
      - open your browser at "localhost:8080/set" (if you deployed remotely, than "location_of_your_deploy/set")
      - open your browser and go to https://dev.hh.ru/admin and edit "Redirect URI" field of your app with https://tommybot-code.onrender.com/get_code (if you deployed remotely, use "location_of_your_deploy/get_code") 
3. next, prepare data for the tommybot set up, and fill the set up form:
   - prepare "client id" and "client secret" of your app at https://dev.hh.ru/admin
   - prepare search links collection for the tommybot (every line is a search link, can be many lines), by your browser just use search at hh.ru, than copy the top line browser link, looks like https://moscow.hh.ru/search/vacancy?hhtmFrom=main&hhtmFromLabel=vacancy_search_line&search_field=name&search_field=company_name&search_field=description&enable_snippets=false&L_save_area=true&area=1&text=java
   - "resume id" — follow you resume link, and get resume id in the url of your top line browser after https://moscow.hh.ru/resume/, looks like f51d091fff033004f00039ed1f5876516d6e42
   - optional: write a cover letter, what is less than 10_000 symbols, and it must be unified for all your vacancy requests (the search links)
> [!IMPORTANT]
> be aware, you can set up all data only once, to set up again, you'll need to delete "api_credentials.json" file and rename ./static/_index.html file to ./static/index.html (if you deployed remotely, then you'll need to redeploy the repository)
4. now, you can start or restart the tommybot by opening "localhost:8080/restart" ("location_of_your_deploy/restart")
