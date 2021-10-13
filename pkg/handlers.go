package pkg

import (
	"encoding/json"
	ctx "github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"reflect"
	"time"
)

func GetTeaHandler(w http.ResponseWriter, r *http.Request) {

	params := ctx.Get(r, "params").(httprouter.Params)
	tea, err := getTeaByName(params.ByName("name"))
	if err != nil {
		log.Printf("failed to get tea data: %v\n", err)
		return
	}

	setHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(TeaResource{Status: http.StatusOK, Data: tea})
}

func CreateTeaHandler(w http.ResponseWriter, r *http.Request) {

	body := ctx.Get(r, "body").(*TeaResource)

	_, err := getTeaByName(body.Data.Name)
	if err == mongo.ErrNoDocuments {
		if err := createTea(body.Data); err != nil {
			log.Printf("failed to create a tea: %v\n", err)
			return
		}

		setHeaders(w, http.StatusCreated)
		json.NewEncoder(w).Encode(TeaResource{Status: http.StatusCreated, Data: body.Data})
		return
	}
	if err != nil {
		log.Printf("failed to get current tea data: %v\n", err)
		return
	}

	log.Println("tea is already exists")
	setErrResp(w, errBadRequest)
}

func DeleteTeaHandler(w http.ResponseWriter, r *http.Request) {

	params := ctx.Get(r, "params").(httprouter.Params)
	if err := deleteTeaByName(params.ByName("name")); err != nil {
		log.Printf("failed to delete a tea: %v\n", err)
		return
	}

	setHeaders(w, http.StatusNoContent)
}

func UpdateTeaHandler(w http.ResponseWriter, r *http.Request) {

	params := ctx.Get(r, "params").(httprouter.Params)
	body := ctx.Get(r, "body").(*TeaResource)
	if err := updateTea(body.Data, params.ByName("name")); err != nil {
		log.Printf("failed to update a tea: %v\n", err)
		return
	}

	setHeaders(w, http.StatusCreated)
}

func GetAllTeasHandler(w http.ResponseWriter, r *http.Request) {

	teas, err := getAllTeas()
	if err != nil {
		log.Printf("failed to update a tea: %v\n", err)
		return
	}

	setHeaders(w, http.StatusOK)
	json.NewEncoder(w).Encode(TeasCollection{Data: teas})
}

func setHeaders(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(status)
}

func setErrResp(w http.ResponseWriter, err *Error) {
	setHeaders(w, err.Status)
	json.NewEncoder(w).Encode(Errors{[]*Error{err}})
}

func AcceptHandler(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("Accept") != "application/vnd.api+json" {
			setErrResp(w, errNotAcceptable)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func ContentTypeHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/vnd.api+json" {
			setErrResp(w, errUnsupportedType)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func RecoverHandler(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				setErrResp(w, errInternalServer)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func LoggingHandler(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func BodyParserHandler(v interface{}) func(http.Handler) http.Handler {

	t := reflect.TypeOf(v)
	m := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			val := reflect.New(t).Interface()
			if err := json.NewDecoder(r.Body).Decode(val); err != nil {
				setErrResp(w, errBadRequest)
				return
			}

			ctx.Set(r, "body", val)
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

	return m
}
