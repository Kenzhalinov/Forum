package handler

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"test/model"
	"test/repository/sqlite"
)

func (h *Manager) CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/create" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	uid, ok := r.Context().Value(UserContextKey).(int)
	if !ok {
		// govno
		fmt.Println("implement me4")
		return
	}

	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		var post model.PostCreateDTO
		if cat := r.FormValue("category1"); cat != "" {
			post.Category += cat
		}
		if cat := r.FormValue("category2"); cat != "" {
			post.Category += cat
		}
		if cat := r.FormValue("category3"); cat != "" {
			post.Category += cat
		}
		if cat := r.FormValue("category4"); cat != "" {
			post.Category += cat
		}

		post.UserID = uid
		post.Title = r.Form.Get("title")
		post.Content = r.Form.Get("content")

		err := h.service.Post.Create(post)
		if errors.Is(err, model.ErrIncorectData) {
			http.Error(w, http.StatusText(http.StatusBadRequest), 400)
			log.Println("impleent me1")
			return
		} else if err != nil {
			// internal
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			log.Println("imple,ent me3")
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
		return
	case http.MethodGet:
		t, err := template.ParseFiles("web/template/create_post.html")
		if err != nil {
			log.Printf("create post: %s\n", err)
			return
		}
		if err := t.Execute(w, nil); err != nil {
			log.Printf("create post: %s\n", err)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Manager) GetPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		postID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("get post: %s\n", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		post, err := h.service.Post.Get(postID)
		if errors.Is(err, sqlite.ErrPostIsNotFound) {
			log.Printf("get post: %s\n", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		} else if err != nil {
			log.Printf("get post: %s\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		comments, err := h.service.Comm.GetByPost(postID)
		if err != nil {
			log.Printf("get post: %s\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("web/template/post.html")
		if err != nil {
			log.Printf("get post: %s\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		data := struct {
			Post     model.Post
			Comments []model.Comment
		}{
			Post:     post,
			Comments: comments,
		}

		if err = tmpl.Execute(w, data); err != nil {
			log.Printf("get post: %s\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Manager) CreateComments(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/comment" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		uid, ok := r.Context().Value(UserContextKey).(int)
		if !ok {
			log.Println("create comment no session")
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")
		comment := &model.CommentCreateDTO{
			PostID:  postID,
			UserID:  uid,
			Content: content,
		}

		if err = h.service.Comm.Create(*comment); errors.Is(err, model.ErrIncorectData) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			log.Printf("create comment: %s\n", err)
			return
		}

		returnPath := r.FormValue("return_path")
		http.Redirect(w, r, returnPath, http.StatusFound)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Manager) CreateVote(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/vote" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		uid, ok := r.Context().Value(UserContextKey).(int)
		if !ok {
			log.Println("create vote no session")
			return
		}

		postID, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
		v := r.FormValue("vote")
		vote := model.Vote{
			PostID: postID,
			UserID: uid,
			Vote:   v == "l",
		}

		if err = h.service.Vote.Vote(vote); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			log.Printf("create vote: %s\n", err)
			return
		}

		returnPath := r.FormValue("return_path")
		http.Redirect(w, r, returnPath, http.StatusFound)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Manager) CreateVoteComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/vote" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodPost:
		uid, ok := r.Context().Value(UserContextKey).(int)
		if !ok {
			log.Println("create vote comment no session")
			return
		}

		postID, err := strconv.Atoi(r.FormValue("comm_id"))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		v := r.FormValue("vote")
		vote := model.Vote{
			PostID: postID,
			UserID: uid,
			Vote:   v == "l",
		}

		if err = h.service.Voco.Vote(vote); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
			log.Printf("create vote: %s\n", err)
			return
		}

		returnPath := r.FormValue("return_path")
		http.Redirect(w, r, returnPath, http.StatusFound)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Manager) ListPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	posts, err := h.service.Post.List(r.FormValue("category"))
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}

	t, err := template.ParseFiles("web/template/home_page.html")
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}

	err = t.Execute(w, posts)
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}
}

func (h *Manager) ListPostCreatedByUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/users" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	uid, ok := r.Context().Value(UserContextKey).(int)
	if !ok {
		log.Println("list post no session")
		return
	}

	posts, err := h.service.Post.ListCreatedByUser(uid)
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}

	t, err := template.ParseFiles("web/template/home_page.html")
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}

	err = t.Execute(w, posts)
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}
}

func (h *Manager) ListPostLikedByUser(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/liked" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	uid, ok := r.Context().Value(UserContextKey).(int)
	if !ok {
		log.Println("list post no session")
		return
	}

	posts, err := h.service.Post.ListLikedByUser(uid)
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}

	t, err := template.ParseFiles("web/template/home_page.html")
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}

	err = t.Execute(w, posts)
	if err != nil {
		fmt.Println("list post: ", err)
		return
	}
}
