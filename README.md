# tommybot
1. at the beginning, you need to get "client id" and "client secret" at https://dev.hh.ru/admin
2. next, to set up the tommybot use deploy_location/set
   [!WARNING]
   be aware, you can set up all data only once, to set up again, you'll need to redeploy the repository
4. next, prepare data for the tommybot set up:
   - prepare search links collection for the tommybot (every line is a search link, can be many lines), by your browser, just use search at hh.ru, than copy the top line broser link, looks like https://moscow.hh.ru/search/vacancy?hhtmFrom=main&hhtmFromLabel=vacancy_search_line&search_field=name&search_field=company_name&search_field=description&enable_snippets=false&L_save_area=true&area=1&text=java
   - "resume id" — follow you resume link, and get resume id in the url of your top line browser after https://moscow.hh.ru/resume/, looks like f51d091fff033004f00039ed1f5876516d6e42
   - optional: write a cover letter, what less than 10000 symbols, and it must be unified for all your vacancy requests (the search links)
5. now you can start the tommybot by using deploy_location/restart
