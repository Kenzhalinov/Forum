package handler

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"test/model"
)

func (h *Manager) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	switch r.Method {
	case http.MethodPost:
		var user model.User
		r.ParseForm()
		user.Email = r.Form.Get("email")
		user.Login = r.Form.Get("login")
		user.Password = r.Form.Get("password")
		if err := h.service.User.Create(user); err != nil {
			t, err := template.ParseFiles("web/template/sign_up.html")
			if err != nil {
				log.Println("error Sign_Up Method Post 1", err)
				return
			}
			t.Execute(w, nil)
			if err != nil {
				log.Println("error Sign_Up Method Post 2", err)
				return
			}
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
		return

	case http.MethodGet:
		t, err := template.ParseFiles("web/template/sign_up.html")
		if err != nil {
			log.Println("error Sign_Up Method Get 1", err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Println("error Sign_Up Method Get 2", err)
			return
		}
	default:
		fmt.Println("gadon na usere")
		return
	}
}

func (h *Manager) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	switch r.Method {
	case http.MethodGet:
		t, err := template.ParseFiles("web/template/sign_in.html")
		if err != nil {
			log.Println("errorSignIn1", err)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			fmt.Fprintln(w, http.StatusMethodNotAllowed)
			return
		}
		return
	case http.MethodPost:
		var login, password string
		r.ParseForm()
		login = r.Form.Get("login")
		password = r.Form.Get("password")

		user, err := h.service.User.Authenticate(login, password)
		if errors.Is(err, model.ErrIncorectData) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Println("Bad reques SignIn", err)
			return
		} else if errors.Is(err, model.ErrUserNotExist) {
			log.Println("Not Found SignIn", err)
			http.Redirect(w, r, "/signup", http.StatusPermanentRedirect)
			return
		} else if errors.Is(err, model.ErrUserIncorrectPasword) {
			log.Println("Incorect Password SignIn", err)
			return
		} else if err != nil {
			fmt.Fprintf(w, "Incorect password or login")
			log.Println("Iternal Server Error SignIn", err)
			return
		}
		cookie, err := h.service.User.CreateSession(user)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Println("Iternal Server Error to Session SignIn", err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:  sessionCookie,
			Value: cookie,
		})
		http.Redirect(w, r, "/", http.StatusFound)
		return
	default:
		fmt.Fprintln(w, http.StatusMethodNotAllowed)
		return
	}
}

func (h *Manager) LogOut(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	uid, ok := r.Context().Value(UserContextKey).(int)
	if !ok {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if err := h.service.User.DeleteSession(uid); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
