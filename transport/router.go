package transport

import "net/http"

func (s *Server) route() {
	sessReq := s.handler.SessMiddleware

	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	s.mux.HandleFunc("/", s.handler.ListPost)

	s.mux.HandleFunc("/signin", s.handler.SignIn)
	s.mux.HandleFunc("/signup", s.handler.SignUp)
	s.mux.HandleFunc("/logout", sessReq(s.handler.LogOut))

	s.mux.HandleFunc("/post/create", s.handler.SessMiddleware(s.handler.CreatePost))

	s.mux.HandleFunc("/post", s.handler.GetPost)
	s.mux.HandleFunc("/post/liked", sessReq(s.handler.ListPostLikedByUser))
	s.mux.HandleFunc("/post/users", sessReq(s.handler.ListPostCreatedByUser))
	s.mux.HandleFunc("/post/comment", sessReq(s.handler.CreateComments))
	s.mux.HandleFunc("/post/vote", sessReq(s.handler.CreateVote))
	s.mux.HandleFunc("/comment/vote", sessReq(s.handler.CreateVoteComment))
}
