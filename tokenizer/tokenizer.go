package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "time"
  "strings"
  jwt "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
  "github.com/google/uuid"
)

var mySigningKey = []byte("balalaika")

// input options
var (
    port_opt = flag.String("port", ":8180", "Port to listen for requests")
)

func GetJWT(session, client, category string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)

    claims["authorized"] = true
    claims["aud"] = "kos&kot"
    claims["iss"] = "koskot"
    claims["exp"] = time.Now().Add(time.Hour * 1000).Unix()
    claims["type"] = category
    claims["client"] = client

    if session == "" {
        uuidWithHyphen := uuid.New()
        uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
        claims["session"] = uuid
    } else {
        claims["session"] = session
    }

    tokenString, err := token.SignedString(mySigningKey)

    if err != nil {
        fmt.Errorf("Something Went Wrong: %s", err.Error())
        return "", err
    }

    return tokenString, nil
}

func index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "You can query following endpoints:\n")
    fmt.Fprintf(w, " /ready - always returns ok\n")
    fmt.Fprintf(w, " /alive - always returns ok\n")
    fmt.Fprintf(w, " /unauthorized - dummy endpoint to test 401\n")
    fmt.Fprintf(w, " /forbidden - dummy endpoint to test 403\n")
    fmt.Fprintf(w, " /debug - master token for querying the whole job queue\n")
    fmt.Fprintf(w, " /client - token equipped with unique/random session_id\n")
    fmt.Fprintf(w, " /service/{service_id} - service identity\n")
}

func debug(w http.ResponseWriter, r *http.Request) {
    validToken, err := GetJWT("debug", "client", "")
    if err != nil {
        fmt.Fprintf(w, "Failed to generate token: %s", err.Error())
    } else {
        fmt.Fprintf(w, string(validToken))
    }
}

func client(w http.ResponseWriter, r *http.Request) {
    validToken, err := GetJWT("", "client", "")
    if err != nil {
        fmt.Fprintf(w, "Failed to generate token: %s", err.Error())
    } else {
        fmt.Fprintf(w, string(validToken))
    }
}

func service(w http.ResponseWriter, r *http.Request) {

    // once again, we will need to parse the path parameters
    vars := mux.Vars(r)
    // we will need to extract the `id` of the article we
    // wish to delete
    service_id := vars["service_id"]

    validToken, err := GetJWT("irrelevant", service_id, "service")
    if err != nil {
        fmt.Fprintf(w, "Failed to generate token: %s", err.Error())
    } else {
        fmt.Fprintf(w, string(validToken))
    }
}

func unauthorized(w http.ResponseWriter, r *http.Request) {
    response := "Unconditionally unauthorized"
    http.Error(w, response, http.StatusUnauthorized)
    return
}


func forbidden(w http.ResponseWriter, r *http.Request) {
    response := "unconditionally forbidden"
    http.Error(w, response, http.StatusForbidden)
    return
}

func handleRequests() {
    // creates a new instance of a mux router
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", index)
    myRouter.HandleFunc("/ready", func (w http.ResponseWriter, r *http.Request){ fmt.Fprintf(w, "ok") })
    myRouter.HandleFunc("/alive", func (w http.ResponseWriter, r *http.Request){ fmt.Fprintf(w, "ok") })
    myRouter.HandleFunc("/debug", debug)
    myRouter.HandleFunc("/unauthorized", unauthorized)
    myRouter.HandleFunc("/forbidden", forbidden)
    myRouter.HandleFunc("/client", client)
    myRouter.HandleFunc("/service/{service_id}", service)
    // finally, instead of passing in nil, we want
    // to pass in our newly created router as the second
    // argument
    log.Fatal(http.ListenAndServe(*port_opt, myRouter))
}

func main() {
    flag.Parse()
    handleRequests()
}
