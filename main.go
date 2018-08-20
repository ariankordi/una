package main

import (
  "github.com/gorilla/mux"
  "github.com/gorilla/sessions"
  "github.com/gorilla/csrf"

  "github.com/gobuffalo/pop"

  "github.com/pkg/errors"
  "golang.org/x/crypto/bcrypt"
  "github.com/gobuffalo/pop/nulls"

  "database/sql"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "regexp"
  "bytes"
  "strings"
  "math/rand"
  "time"

  "./models"

  "log"
)

var secret = []byte("sdfjhskjhfakljdfhineorfnaoaaaaAAAA_=-=;f=sd;f=ak5l390tgr.e0-a.9t4,g0.r")
var sessionStore sessions.Store
var db *pop.Connection

func Login(w http.ResponseWriter, r *http.Request) {
  username := r.PostFormValue("username")
  password := r.PostFormValue("password")
  // Check the login form
  if username == "" || password == "" {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "You have to enter a username, and password.",
    })
    return
  }
  var user models.User
	err := db.Select("id", "password").Where("username = ?", username).First(&user)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
      w.WriteHeader(http.StatusNotFound)
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{
        "error": "Your account couldn't be found.",
      })
			return
		}
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Something went wrong while searching for your account.\n" + err.Error(),
    })
		return
	}
	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
      w.WriteHeader(http.StatusUnauthorized)
      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(map[string]string{
        "error": "Your password didn't work.",
      })
			return
		}
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Something went wrong while verifying your password.\n" + err.Error(),
    })
		return
	}
  session, _ := sessionStore.Get(r, "una_session")
  session.Values["user"] = user.ID
  session.Save(r, w)
  w.WriteHeader(http.StatusAccepted)
}

func Signup(w http.ResponseWriter, r *http.Request) {
  username := r.PostFormValue("username")
  nickname := r.PostFormValue("nickname")
  password := r.PostFormValue("password")
  passwordAgain := r.PostFormValue("password_again")
  // Check the signup form
  if username == "" || nickname == "" || password == "" || passwordAgain == "" {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "You have to enter a username, nickname, and password.",
    })
		return
	}
  if password != passwordAgain {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Your passwords don't match.",
    })
		return
	}
  nameSyntaxValid, _ := regexp.MatchString("^[A-Za-z0-9][^/]{1,32}$", username)
  if !nameSyntaxValid {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Your username contains a forward slash (/), has no letters, or is too long or too short.",
    })
    return
  }
  if nickname != "" && len([]rune(nickname)) > 64 {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Your nickname is too long.",
    })
    return
  }
  // End form "syntax" validation.
  nameExists, err := db.Where("username = ?", username).Count(&models.User{})
	if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Something went wrong while trying to verify your username.\n" + err.Error(),
    })
		return
	}
	if nameExists > 0 {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Your username is already taken.",
    })
		return
	}
  // End all checks, make account.
  newPassword, err := bcrypt.GenerateFromPassword([]byte(password), 13)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Something went wrong while trying to store your password.\n" + err.Error(),
    })
    return
  }
  newUser := models.User{
    Username: username,
    Nickname: nickname,
    Password: newPassword,
  }
  err = db.Create(&newUser)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Something went wrong while trying to create your account.\n" + err.Error(),
    })
    return
  }

  session, _ := sessionStore.Get(r, "una_session")
  session.Values["user"] = newUser.ID
  session.Save(r, w)

  w.WriteHeader(http.StatusCreated)
}

// CreateAnonymousUser takes a request and a session, and if there's no user on the session,
// and no anonymous user, an anonymous user is created and put on the session.
func CreateAnonymousUser(w http.ResponseWriter, r *http.Request, session *sessions.Session) {
  if session.Values["user"] == nil && session.Values["anonymous_user"] == nil {
    source := rand.NewSource(time.Now().UnixNano())
    rand1 := rand.New(source)
    // Generate a random number for their name ID
    newUser := models.AnonymousUser{
      NameID: 1000 + rand1.Intn(8999),
    }
    db.Create(&newUser)
    // No error handling lol
    session.Values["anonymous_user"] = newUser.ID
    session.Save(r, w)
  }
}

func Logout(w http.ResponseWriter, r *http.Request) {
  session, _ := sessionStore.Get(r, "una_session")
  session.Values["user"] = nil
  session.Save(r, w)
  w.WriteHeader(http.StatusAccepted)
}


func CreateLobby(w http.ResponseWriter, r *http.Request) {
  name := r.PostFormValue("name")
  if name == "" {
    w.WriteHeader(http.StatusBadRequest)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "You have to enter a name.",
    })
    return
  }
  session, _ := sessionStore.Get(r, "una_session")
  CreateAnonymousUser(w, r, session)
  creatorID := 0
  anonymousCreatorID := 0
  if session.Values["user"] != nil {
    creatorID = session.Values["user"].(int)
  } else {
    anonymousCreatorID = session.Values["anonymous_user"].(int)
  }
  newLobby := models.Lobby{
    Name: name,
    CreatorID: nulls.Int{Int: creatorID, Valid: session.Values["user"] != nil},
    AnonymousCreatorID: nulls.Int{Int: anonymousCreatorID, Valid: session.Values["user"] == nil},
  }
  err := db.Create(&newLobby)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Something went wrong while trying to create your lobby.\n" + err.Error(),
    })
    return
  }
  w.WriteHeader(http.StatusCreated)
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(newLobby)
}

func GetOpenLobbies(w http.ResponseWriter, r *http.Request) {
  var lobbies []models.Lobby
  err := db.Order("id desc").All(&lobbies)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
      "error": "Something went wrong while searching for open lobbies.\n" + err.Error(),
    })
		return
	}
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(lobbies)
}

func main() {
  // Connect to the Pop database
  dbConnect, err := pop.Connect("development")
  if err != nil {
    if err.Error() == "tried to load pop configuration file, but couldn't find it" {
      log.Fatal(err.Error() + "\n\033[1mDoes database.yml exist?\033(B\033[m Make sure to create a database.yml from database-example.yml.")
    }
    log.Fatal(err)
  }
  db = dbConnect
  // Initialize the session cookie store
  sessionStore = sessions.NewCookieStore(secret)
  // Initialize the router
  router := mux.NewRouter()
  // Make another router for API routes
  app := router.PathPrefix("/app").Subrouter()
  app.HandleFunc("/login", Login).Methods("POST")
  app.HandleFunc("/signup", Signup).Methods("POST")
  app.HandleFunc("/logout", Logout).Methods("POST")
  // Game routes
  app.HandleFunc("/lobbies", GetOpenLobbies)
  app.HandleFunc("/lobby_create", CreateLobby)

  // Serve static files at /static/
  router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
  // Read index.html
  html, err := ioutil.ReadFile("./index.html")
  if err != nil {
    html = []byte(err.Error())
  }
  // Set the not found handler to serve index.html, so that Vue can handle routing
  router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Return an actual 404 if /static/ or /app/ is in the URL
    if strings.HasPrefix(r.URL.Path, "/static/") || strings.HasPrefix(r.URL.Path, "/app/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
    // Get the user
    session, _ := sessionStore.Get(r, "una_session")
    userID := session.Values["user"]
    user := []byte("null")
    if userID != nil {
      var meUser models.User
      db.Select("username", "nickname", "avatar", "created_at", "updated_at").Where("id = ?", userID).First(&meUser)
      meUser.ID = userID.(int)
      user, _ = json.Marshal(meUser)
    }
    template := bytes.Replace(html, []byte("__CSRF_TOKEN"), []byte(csrf.Token(r)), 1)
    template = bytes.Replace(template, []byte("__USER_DATA"), user, 1)
    w.Header().Set("Content-Type", "text/html")
    w.Write(template)
  })
  // Run the server
  log.Println("Listening on :8080")
  log.Fatal(http.ListenAndServe(":8080", csrf.Protect(secret, csrf.Secure(false))(router)))
}
