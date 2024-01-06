// all input data of client id will be copied to the get code link
document.addEventListener('DOMContentLoaded', function() {
    // data
    var input_field_client_id = document.getElementById('client_id');
    var input_field_client_secret = document.getElementById('client_secret');
    var get_a_code_link = document.getElementById('get_a_code_link');
    var input_field_code = document.getElementById('code');
    var get_access_and_refresh_tokens_link = document.getElementById('get_access_and_refresh_tokens_link');
    var get_access_and_refresh_tokens_message = document.getElementById('get_access_and_refresh_tokens_message');
    var input_field_access_token = document.getElementById('access_token');
    var input_field_refresh_token = document.getElementById('refresh_token');
    var input_field_search_links = document.getElementById('search_links');
    var input_field_resume_id = document.getElementById('resume_id');
    var input_field_cover_letter = document.getElementById('cover_letter');
    var save_button = document.getElementById('save_button');
    var settings_was_stored = document.getElementById('settings_was_stored');
    var post_data = {
        ACCESS_TOKEN:   "",
        REFRESH_TOKEN:  "",
        SEARCH_LINKS:   "",
        RESUME_ID:      "",
        COVER_LETTER:   ""
    };
    
    var input_field_client_id_filled = false;
    var input_field_client_secret_filled = false;
    var input_field_code_filled = false;
    var input_field_access_token_filled = false;
    var input_field_refresh_token_filled = false;
    var input_field_search_links_filled = false;
    var input_field_resume_id_filled = false;
    
    get_a_code_link.classList.add('disabled');
    input_field_code.classList.add('disabled');
    get_access_and_refresh_tokens_link.classList.add('disabled');
    save_button.classList.add('disabled');
    // END data

    // attach input_field_client_id to the link "get_a_code_link"
    input_field_client_id.addEventListener('input', function() {
        var input_field_client_id_value = input_field_client_id.value;
        if (input_field_client_id_value.trim().length == 64) {
            get_a_code_link.href = "https://hh.ru/oauth/authorize?response_type=code&client_id=" + input_field_client_id_value;
            input_field_client_id_filled = true;
        } else {
            get_a_code_link.href = "";
            input_field_client_id_filled = false;
        }
    });

    // disable "get a code" link until data uppears in the client_id input field
    input_field_client_secret.addEventListener('input', function() {
        if (
            input_field_client_id_filled == true
            && input_field_client_secret.value.trim().length == 64
            && input_field_client_secret.value.trim() != input_field_client_id.value.trim()
        ) {
            get_a_code_link.classList.remove('disabled');
            input_field_code.classList.remove('disabled');
            input_field_client_secret_filled = true;
        } else {
            get_a_code_link.classList.add('disabled');
            input_field_code.classList.add('disabled');
            input_field_client_secret_filled = false;
        }
    });

    // attach input_field_client_id, input_field_client_secret, and input_field_code to the link "get_access_and_refresh_tokens"
    input_field_code.addEventListener('input', function() {
        var input_field_code_value = input_field_code.value;
        if (input_field_code_value.trim().length == 64) {
            // get_access_and_refresh_tokens_link.href = "get_access_and_refresh_tokens?client_id=" + input_field_client_id.value + "&client_secret=" + input_field_client_secret.value + "&code=" + input_field_code.value;
            get_access_and_refresh_tokens_link.classList.remove('disabled');
            input_field_code_filled = true;
        } else {
            get_access_and_refresh_tokens_link.href = "";
            get_access_and_refresh_tokens_link.classList.add('disabled');
            input_field_code_filled = false;
        }
    });

    // attach get request event to the get_access_and_refresh_tokens_link
    get_access_and_refresh_tokens_link.addEventListener('click', function() {
        var protocol = window.location.protocol + "//";
        var host = window.location.hostname + ":";
        var port = window.location.port || "8080";
        var left_side = protocol + host + port;
        var url = left_side + "/get_access_and_refresh_tokens?client_id=" + input_field_client_id.value + "&client_secret=" + input_field_client_secret.value + "&code=" + input_field_code.value;
        get_request(url);
    });

    // attach event to the input access_token
    input_field_access_token.addEventListener('input', function() {
        var input_field_access_token_value = input_field_access_token.value;
        if (input_field_access_token_value.trim().length == 64) {
            input_field_access_token_filled = true;
        } else {
            input_field_access_token_filled = false;
        }
    });

    // attach event to the input refresh_token
    input_field_refresh_token.addEventListener('input', function() {
        var input_field_refresh_token_value = input_field_refresh_token.value;
        if (input_field_refresh_token_value.trim().length == 64) {
            input_field_refresh_token_filled = true;
        } else {
            input_field_refresh_token_filled = false;
        }
    });

    // attach formatting data event to the search_links field
    input_field_search_links.addEventListener('input', function() {
        var input_field_search_links_value = input_field_search_links.value;
        if (input_field_search_links_value.length > 0) {
            input_field_search_links_filled = true;
        } else {
            input_field_search_links_filled = false;
        }
    });

    // attach event to the input_field_resume_id
    input_field_resume_id.addEventListener('input', function() {
        var input_field_resume_id_value = input_field_resume_id.value;
        if (input_field_resume_id_value.trim().length > 0) {
            input_field_resume_id_filled = true;
            // save_button.classList.remove('disabled');
        } else {
            input_field_resume_id_filled = false;
            // save_button.classList.add('disabled');
        }
    });

    // save button event
    save_button.addEventListener('click', function() {
        if (
            input_field_access_token_filled == true
            && input_field_refresh_token_filled == true
            && input_field_search_links_filled == true
            && input_field_resume_id_filled == true
        ) {
            // prepare post data
            post_data = {
                ACCESS_TOKEN:   input_field_access_token.value,
                REFRESH_TOKEN:  input_field_refresh_token.value,
                SEARCH_LINKS:   input_field_search_links.value.trim().split('\n').filter(line => line.trim() !== ''),
                RESUME_ID:      input_field_resume_id.value,
                COVER_LETTER:   encodeURIComponent(input_field_cover_letter.value)
            };
            console.log(post_data);

            // submit the form via post request
            var protocol = window.location.protocol + "//";
            var host = window.location.hostname + ":";
            var port = window.location.port || "8080";
            var left_side = protocol + host + port;
            post_request(left_side + "/set", post_data);
        }
    });

    // attach event to the the_form
    document.getElementById('the_form').addEventListener('input', function() {
        if (
            input_field_access_token_filled == true
            && input_field_refresh_token_filled == true
            && input_field_search_links_filled == true
            && input_field_resume_id_filled == true
        ) {
            // save_button.style.opacity = 1;
            save_button.classList.remove('disabled');
        } else {
            // save_button.style.opacity = .4;
            save_button.classList.add('disabled');
        }
    });

    // making a GET request using fetch
    function get_request (url) {
        fetch(url, {
            method: 'GET',
            headers: {
                'User-Agent': 'Auto Rectuiter/1.0',
                'Content-Type': 'application/json',
            },
        })
            .then(response => response.json()) // parse the JSON response
            .then(data => {
                // success:
                if (data.error) {
                    get_access_and_refresh_tokens_message.innerHTML = data.error_description;
                } else {
                    input_field_access_token.value = data.access_token;
                    input_field_refresh_token.value = data.refresh_token;
                    get_access_and_refresh_tokens_message.innerHTML = "&nbsp;";
                    console.log('Success:', data);
                }
            })
            .catch(error => {
                // error:
                console.error('Error:', error);
            });
    }

    // making a POST request using fetch
    function post_request (url, data) {
        fetch(url, {
            method: 'POST',
            headers: {
                'User-Agent': 'Auto Rectuiter/1.0',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data) // convert the data to JSON format
        })
            .then(response => response.json()) // parse the JSON response
            .then(data => {
                console.log('Success:', data);
                settings_was_stored.innerHTML = "data stored successfully";
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }
});