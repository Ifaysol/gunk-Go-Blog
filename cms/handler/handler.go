package handler

import (
	//"go/token"
	"net/http"

	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"


	tpb "grpc-blog/gunk/v1/categories"
	ppb "grpc-blog/gunk/v1/post"
)

//const sessionName = "cms-session"

type Handler struct {
	templates *template.Template
	decoder   *schema.Decoder
	sess      *sessions.CookieStore
	cc        tpb.CategoryServiceClient
	pc        ppb.PostServiceClient
}

func New(decoder *schema.Decoder, sess *sessions.CookieStore, cc tpb.CategoryServiceClient, pc ppb.PostServiceClient) *mux.Router {
	h := &Handler{
		decoder: decoder,
		sess:    sess,
		cc: cc,
		pc: pc,
	}

	h.parseTemplate()

	r := mux.NewRouter()
	// s := r.NewRoute().Subrouter()

	//r.HandleFunc("/", h.Home)

	r.HandleFunc("/categories/create", h.createCategory)
	r.HandleFunc("/categories/store", h.storeCategory)
	r.HandleFunc("/categories", h.listCategory)
	r.HandleFunc("/categories/{id}/delete", h.deleteCategory)
	r.HandleFunc("/categories/{id}/edit", h.editCategory)
	r.HandleFunc("/categories/{id}/update", h.updateCategory)

	r.HandleFunc("/posts/create", h.createPost)
	r.HandleFunc("/posts/store", h.storePost)
	r.HandleFunc("/posts", h.listPost)
	r.HandleFunc("/posts/{id}/delete", h.deletePost)
	r.HandleFunc("/posts/{id}/edit", h.editPost)
	r.HandleFunc("/posts/{id}/view", h.viewPost)
	r.HandleFunc("/posts/{id}/update", h.updatePost)
	r.PathPrefix("/asset/").Handler(http.StripPrefix("/asset/", http.FileServer(http.Dir("./"))))

	//s.Use(h.authMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := h.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return r
}

func (h *Handler) parseTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"cms/assets/templates/create-category.html",
		"cms/assets/templates/list-category.html",
		"cms/assets/templates/edit-category.html",
		"cms/assets/templates/create-post.html",
		"cms/assets/templates/list-post.html",
		"cms/assets/templates/edit-post.html",
		"cms/assets/templates/view-post.html",
		"cms/assets/templates/404.html",
	))
}

// func (h *Handler) authMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(rw, r)
// 		return
// 		// session, _ := h.sess.Get(r, sessionName)
// 		// // if err != nil {
// 		// // 	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
// 		// // 	return
// 		// // }

// 		// // Check if user is authenticated
// 		// ok := session.Values["authUserId"]
// 		// if ok == nil {
// 		// 	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
// 		// 	return
// 		// }
// 		// next.ServeHTTP(rw, r)
// 	})
//}
