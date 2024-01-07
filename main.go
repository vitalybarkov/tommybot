// the server with the robot
package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson" // need to install by terminal by: "go get -u github.com/tidwall/gjson"
)

// global variables
var ROBOT_STARTED bool = false
var ROBOT_RESTARTED bool = false
var RESTARTS_COUNT int = 0
var ACCESS_TOKEN string = ""
var REFRESH_TOKEN string = ""
var ERROR_CODE int = 0 // 0 — good, 1 — already_applied, 2 — limit_exceeded, 3 — token_expired
var MAX_REQUESTS int = 15000

var SEARCH_LINKS []string = []string{}
var RESUME_ID = ""
var COVER_LETTER string = ""
var API_CREDENTIALS_FILE_LOCATION string = "api_credentials.json"

type store_data struct {
	ACCESS_TOKEN  string   `json:"ACCESS_TOKEN"`
	REFRESH_TOKEN string   `json:"REFRESH_TOKEN"`
	SEARCH_LINKS  []string `json:"SEARCH_LINKS"`
	RESUME_ID     string   `json:"RESUME_ID"`
	COVER_LETTER  string   `json:"COVER_LETTER"`
}

// END global variables

// base functionality
func request(method string, url string) (response string) { // base request function
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Auto Rectuiter/1.0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+ACCESS_TOKEN)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}

func get_ids(data string) (ids []string) { // get ids of vacancies without applies
	values := gjson.Get(data, `items.#(relations.#==0)#.id`)
	// values := gjson.Get(data, "items.#.id") // just get ids of all vacancies
	var list []string

	r := gjson.Parse(values.Raw).Array()
	for _, value := range r {
		list = append(list, value.Str)
	}

	return list
}

func get_error_code(data string) (outcome_error_code int) { // if exist, get an error code from the response
	values := gjson.Get(data, "errors.0.value")
	var error_code int = 0 // all_good_code

	r := gjson.Parse(values.Raw).Str
	if r == "limit_exceeded" {
		error_code = 2 // limit_exceeded_code
	} else if r == "token_expired" || r == "bad_authorization" {
		error_code = 3 // token_expired_code
	}
	//  else if r == "already_applied" {
	// 	error_code = 1 // already_applied_code
	// }

	return error_code
}

func search_request(search_string string, page string, per_page string) (ids_from_data []string, pages_from_data int) { // search request
	var right_side_of_the_request string = prepare_search_string_for_search_request(search_string, page, per_page)
	var url string = "https://api.hh.ru/vacancies?" + right_side_of_the_request
	var data string = request("GET", url)
	ERROR_CODE = get_error_code(data)
	var pages int = int(gjson.Get(data, "pages").Num)
	time.Sleep(500 * time.Millisecond)

	return get_ids(data), pages
}

func prepare_search_string_for_search_request(search_string string, page string, per_page string) string {
	// split the search_string based on the "?"
	split_result := strings.Split(search_string, "?")
	var right_side_of_the_request string = ""
	if len(split_result) > 1 {
		right_side_of_the_request = split_result[1]
	} else {
		fmt.Println("wrong search string format, please, check your search strings for '?' symbol included, a search string need to be looks like a link with attributes")
		return ""
	}
	// END split the search_string based on the "?"

	// check existence "page=" in the right_side_of_the_request, than replace or add it with the income page variable
	if strings.Contains(right_side_of_the_request, "page=") {
		// replace "page=" with the value of the income page variable
		re := regexp.MustCompile(`page=(\d+)`)
		match := re.FindStringSubmatch(right_side_of_the_request)
		if len(match) > 1 {
			right_side_of_the_request = strings.Replace(right_side_of_the_request, match[1], page, 1)
		}
	} else {
		// add "&page=" with the value of the income page variable to the end
		right_side_of_the_request = right_side_of_the_request + "&page=" + page
	}
	// END check existence "page=" in the right_side_of_the_request, than replace or add it with the income page variable

	// check existence "per_page=" in the right_side_of_the_request, than replace or add it with the income per_page variable
	if strings.Contains(right_side_of_the_request, "per_page=") {
		// replace "per_page=" with the value of the income per_page variable
		re := regexp.MustCompile(`per_page=(\d+)`)
		match := re.FindStringSubmatch(right_side_of_the_request)
		if len(match) > 1 {
			right_side_of_the_request = strings.Replace(right_side_of_the_request, match[1], per_page, 1)
		}
	} else {
		// add "&per_page=" with the value of the income per_page variable to the end
		right_side_of_the_request = right_side_of_the_request + "&per_page=" + per_page
	}
	// END check existence "per_page=" in the right_side_of_the_request, than replace or add it with the income per_page variable

	return right_side_of_the_request
}

func save_credentials_to_file(filename string) bool {
	if len(ACCESS_TOKEN) == 64 && len(REFRESH_TOKEN) == 64 && len(RESUME_ID) > 0 && len(SEARCH_LINKS) > 0 {
		data := store_data{
			ACCESS_TOKEN:  ACCESS_TOKEN,
			REFRESH_TOKEN: REFRESH_TOKEN,
			RESUME_ID:     RESUME_ID,
			SEARCH_LINKS:  SEARCH_LINKS,
			COVER_LETTER:  COVER_LETTER,
		}
		content, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		err = os.WriteFile(filename, content, 0644)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		return false
	}
}

func load_credentials_from_file(filename string) (string, string, []string, string, string, error) {
	r, err := os.ReadFile(filename)
	if err != nil {
		return "", "", nil, "", "", err
	}

	data := store_data{}
	err = json.Unmarshal(r, &data)
	if err != nil {
		return "", "", nil, "", "", err
	}

	return data.ACCESS_TOKEN, data.REFRESH_TOKEN, data.SEARCH_LINKS, data.RESUME_ID, data.COVER_LETTER, err
}

// END base functionality

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
		if ROBOT_RESTARTED == false {
			// stopping the current robot
			ROBOT_RESTARTED = true
			// END stopping the current robot
			// restarting a new one
			RESTARTS_COUNT++
			ROBOT_STARTED = false
			var protocol string = "http://"
			if r.TLS != nil {
				protocol = "https://"
			}
			var requestURL string = fmt.Sprintf(`%s%s/hi`, protocol, r.Host)
			request("GET", requestURL)
			// END restarting a new one robot
			// reopen restart function after 12 minutes
			go func() { // using this thread for sleep 12 minutes
				time.Sleep(720 * time.Second)
				ROBOT_RESTARTED = false
			}()
			// END reopen restart function after 12 minutes
			if ROBOT_STARTED == false {
				fmt.Fprintf(w, "the robot do not have settings, please, set up the tommybot first, at "+`%s%s/set`, protocol, r.Host)
			} else {
				fmt.Fprintf(w, "the robot was restarted, for next restart come back after 12 minutes")
			}
		} else {
			fmt.Fprintf(w, "the robot's restart function on hold now, please, come back later (approximately 12 minutes after the robot restart)")
		}
	})

	http.HandleFunc("/get_code", func(w http.ResponseWriter, r *http.Request) {
		param1 := r.URL.Query().Get("code")
		if param1 != "" {
			fmt.Fprintf(w, "%s", param1)
		} else {
			fmt.Fprintf(w, "no code")
		}
	})

	http.HandleFunc("/get_access_and_refresh_tokens", func(w http.ResponseWriter, r *http.Request) {
		client_id := r.URL.Query().Get("client_id")
		client_secret := r.URL.Query().Get("client_secret")
		code := r.URL.Query().Get("code")

		client := &http.Client{}
		req, err := http.NewRequest("POST", "https://hh.ru/oauth/token?grant_type=authorization_code&client_id="+client_id+"&client_secret="+client_secret+"&code="+code, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("User-Agent", "Auto Rectuiter/1.0")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(string(body))
		fmt.Fprintf(w, "%s", string(body))
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		// change the working directory to the directory containing the executable
		exe_path, err := os.Executable()
		if err != nil {
			fmt.Println("Error getting executable path:", err)
			return
		}
		exe_dir := filepath.Dir(exe_path)
		if err := os.Chdir(exe_dir); err != nil {
			fmt.Println("Error changing working directory:", err)
			return
		}
		// END change the working directory to the directory containing the executable

		filename := "./static/index.html"
		_, err = os.Stat(filename)
		if err == nil {
			if r.Method == http.MethodGet { // GET
				// check if file exists, then return the html file
				file_settings := API_CREDENTIALS_FILE_LOCATION
				if _, err := os.Stat(file_settings); err == nil {
					fmt.Fprintf(w, "the tommybot was set, to reset you need to redeploy the application")
				} else if os.IsNotExist(err) {
					http.ServeFile(w, r, filename)
				} else {
					fmt.Println("Error:", err)
				}
				// END check if file exists, then return the html file
			} else if r.Method == http.MethodPost { // POST
				// unmarshal json data
				data := store_data{}
				err := json.NewDecoder(r.Body).Decode(&data)
				if err != nil {
					http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
					return
				}
				// END unmarshal json data

				// store the data in the settings file
				ACCESS_TOKEN = data.ACCESS_TOKEN
				REFRESH_TOKEN = data.REFRESH_TOKEN
				SEARCH_LINKS = data.SEARCH_LINKS
				RESUME_ID = data.RESUME_ID
				COVER_LETTER = data.COVER_LETTER
				if save_credentials_to_file(API_CREDENTIALS_FILE_LOCATION) {
					fmt.Printf("the data was store in the setting\n")
				} else {
					http.Error(w, "the data values has wrong sizes", http.StatusBadRequest)
					return
				}
				// END store the data in the settings file

				// get API credentials from the JSON file
				ACCESS_TOKEN, REFRESH_TOKEN, SEARCH_LINKS, RESUME_ID, COVER_LETTER, err = load_credentials_from_file(API_CREDENTIALS_FILE_LOCATION)
				if err != nil {
					fmt.Println("Error getting API credentials:", err)
					return
				}
				// END get API credentials from the JSON file

				// rename the index file
				old_file_name := filename
				new_file_name := "./static/_index.html"
				err = os.Rename(filename, new_file_name)
				if err != nil {
					fmt.Println("Error renaming file:", err)
					return
				}
				fmt.Printf("file %s renamed to %s\n", old_file_name, new_file_name)
				// END rename the index file

				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"message": "data stored successfully"}`))
			}
			return
		} else if os.IsNotExist(err) {
			fmt.Printf("file %s does not exist\n", filename)
			http.Error(w, "the tommybot was set, to reset you need to redeploy the application", http.StatusBadRequest)
			return
		} else {
			fmt.Printf("error checking file existence: %v\n", err)
			http.Error(w, "Error checking file existence", http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		// the robot
		// check if settings file exists
		if ACCESS_TOKEN == "" || REFRESH_TOKEN == "" || SEARCH_LINKS == nil || RESUME_ID == "" || COVER_LETTER == "" {
			file_settings := API_CREDENTIALS_FILE_LOCATION
			if _, err := os.Stat(file_settings); err == nil {
				// get API credentials from the JSON file
				ACCESS_TOKEN, REFRESH_TOKEN, SEARCH_LINKS, RESUME_ID, COVER_LETTER, err = load_credentials_from_file(file_settings)
				if err != nil {
					fmt.Println("Error getting API credentials:", err)
					return
				}
			} else if os.IsNotExist(err) {
				fmt.Println("Error:", err)
				return
			} else {
				fmt.Println("Error:", err)
				return
			}
		}
		// END check if settings file exists

		go func() { // using this thread for the fast way response
			if ROBOT_STARTED == false {
				var requests_counter int = 1
				ROBOT_STARTED = true
				for { // infinity loop using for start the work every 24 hours
					if requests_counter < MAX_REQUESTS { // maximum requests using for defence from unstoppable requesting if the responses don't give key phrases for proper working the daily plan work
						// daily plan work
						log.Println("work started // " + time.Now().Format("02-01-2006 15:04:05"))

						// testing the access_token for a expiration
						var test_data string = request("GET", "https://api.hh.ru/me")
						ERROR_CODE = get_error_code(test_data)
						// END testing the access_token for a expiration

						// doing the work for every search strings
						for i := 0; i < len(SEARCH_LINKS); i++ {
							// if got error "limit_exceeded" than break the loop
							if ERROR_CODE == 2 {
								break
							}
							// END if got error "limit_exceeded" than break the loop

							// if got error code "token_expired", than get new access_token and refresh_token and store new ones on relevant variables
							if ERROR_CODE == 3 {
								var data string = request("POST", "https://hh.ru/oauth/token?grant_type=refresh_token&refresh_token="+REFRESH_TOKEN) // getting new access_token and refresh_token pair
								if gjson.Get(data, "error").Exists() {                                                                               // if answer with error, than break the loop
									log.Println(data)
									break
								} else { // if good answer, than refresh the collected data into relevant variables
									log.Println(data)
									ACCESS_TOKEN = gjson.Parse(gjson.Get(data, "access_token").Raw).Str
									REFRESH_TOKEN = gjson.Parse(gjson.Get(data, "refresh_token").Raw).Str

									// save API credentials to a JSON file
									if save_credentials_to_file(API_CREDENTIALS_FILE_LOCATION) {
										fmt.Printf("the data was store in the setting\n")
									} else {
										http.Error(w, "the data values has wrong sizes", http.StatusBadRequest)
										return
									}
									// END save API credentials to a JSON file
									ERROR_CODE = 0
								}

							}
							// END if got error "token_expired", than get new access_token and refresh_token and store new ones on relevant variables

							// get ids collection of all pages of the search string
							var search_string string = SEARCH_LINKS[i]
							var page string = "0"
							var per_page string = "100"
							ids_collection := make([][]string, 1)
							ids, pages := search_request(search_string, page, per_page)
							ids_collection[0] = ids
							if pages > 0 {
								for i := 1; i < pages; i++ {
									local_ids, _ := search_request(search_string, strconv.Itoa(i), per_page)
									ids_collection = append(ids_collection, local_ids)

								}
							}
							log.Println(search_string)
							log.Println(len(ids_collection))
							log.Println(ids_collection)
							// END get ids collection of all pages of the search_string

							// make threads for all cells of ids_collection
							var wg sync.WaitGroup
							for i := 0; i < len(ids_collection); i++ { // enumerate pages of search_string
								wg.Add(1)
								go func(local_i int) { // create a new thread
									defer wg.Done()
									// applying to all of vacancies
									for i := 0; i < len(ids_collection[local_i]); i++ { // enumerate vacancy ids
										// if got error "limit_exceeded" or "token_expired" than break the loop
										if ERROR_CODE == 2 || ERROR_CODE == 3 {
											break
										}
										// END if got error "limit_exceeded" or "token_expired" than break the loop

										var data string = (request("POST", "https://api.hh.ru/negotiations?vacancy_id="+ids_collection[local_i][i]+"&resume_id="+RESUME_ID+"&message="+COVER_LETTER))
										requests_counter++
										log.Println(data)
										ERROR_CODE = get_error_code(data)
										time.Sleep(1500 * time.Millisecond) // sleep for 1.5 seconds after a one request
									}
									// END applying to all of vacancies
								}(i)
								time.Sleep(3000 * time.Millisecond) // sleep for 1.5 seconds after a one thread finished the work
							}
							wg.Wait()
							// END ake threads for all cells of ids_collection
						}
						// END doing the work for every SEARCH_LINKS

						ERROR_CODE = 0 // reset the ERROR_CODE to all_good_code, needs for every day plan work
						log.Println("work done // " + time.Now().Format("02-01-2006 15:04:05"))
						// END daily plan work
					}
					// waiting for 24.15 hours for starting a new daily plan work
					var k int64 = 0
					for {
						var protocol string = "http://"
						if r.TLS != nil {
							protocol = "https://"
						}
						var requestURL string = fmt.Sprintf(`%s%s/`, protocol, r.Host)
						log.Printf("%s %s", strconv.FormatInt(k, 10), request("GET", requestURL)) // this request to this web server needs for awake the server every 10 minutes (it's just a render.com rule for free web services)
						k++
						if k > 145 {
							break
						}
						time.Sleep(600 * time.Second)
						if RESTARTS_COUNT > 0 {
							RESTARTS_COUNT = RESTARTS_COUNT - 1
							// ROBOT_RESTARTED = false
							return // stop the current robot
						}
					}
					// END waiting for 24.15 hours for starting a new daily plan work
				}
			}
		}()
		if ROBOT_STARTED == false {
			fmt.Fprintf(w, "the robot was started and working now")
		} else {
			fmt.Fprintf(w, "the robot working now, and you can't start more instances of the robot")
		}
		// END the robot
	})

	// run the server
	port := os.Getenv("PORT")
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	if port == "" {
		port = "8080"
	}
	fmt.Println("tommybot server runs at port: " + port)
	fmt.Println("// to change the port, run the app with first argument of your port number")
	log.Fatal(http.ListenAndServe(":"+port, nil))
	// END run the server
}

// END the server with the robot
