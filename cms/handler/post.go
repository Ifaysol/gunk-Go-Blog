package handler

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"net/http"
	"strconv"

	tpb "grpc-blog/gunk/v1/categories"
	ppb "grpc-blog/gunk/v1/post"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Post struct {
	ID        int64  `db:"id"`
	CatID     int64  `db:"catid"`
	PostName  string `db:"postname"`
	PostImage string `db:"image"`
	CategoryName string 
	Errors    map[string]string
}

type ListPost struct {
	Post []Post
}

func (p *Post) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.PostName, validation.Required.Error("Must insert a post name"), validation.Length(3, 0)),
		
	)
}

type FormData struct {
	Post Post
	Category []Category
	Errors   map[string]string
}

func (h *Handler) createPost(rw http.ResponseWriter, r *http.Request) {
	post := Post{}
	vErrs := map[string]string{}
	res, err := h.cc.List(r.Context(), &tpb.ListCategoryRequest{})
	if err != nil {
		log.Fatal(err)
	}
	category:= []Category{}

    for _,v:=range res.Category{
	category = append(category, Category{
		ID: int(v.ID),
		CategoryName: v.CategoryName,
	})
}
	h.loadCreatedPostForm(rw, post, category, vErrs)

}

func (h *Handler) storePost(rw http.ResponseWriter, r *http.Request) {

	res, err := h.cc.List(r.Context(), &tpb.ListCategoryRequest{})
	
	if err != nil {
		log.Fatal(err)
	}
	category:= []Category{}

    for _,v:=range res.Category{
	category = append(category, Category{
		ID: int(v.ID),
		CategoryName: v.CategoryName,
	})
}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var post Post
	if err := h.decoder.Decode(&post, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := post.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[strings.Title(key)] = value.Error()

			}
			h.loadCreatedPostForm(rw, post,category,vErr)
			return

		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	file, _, err := r.FormFile("PostImage")
    if file == nil {
		vErrs := map[string]string{"PostImage" : "The image field is required"}
		h.loadCreatedPostForm(rw, post, category, vErrs)
			return
	}
    
	if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
    defer file.Close()
   
	img := "upload-*.png"
    tempFile, err := ioutil.TempFile("cms/assets/images", img)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
    defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
	
    tempFile.Write(fileBytes)
	a := tempFile.Name()
	

	_, err = h.pc.CreatePost(r.Context(), &ppb.CreatePostRequest{
		Post: &ppb.Post{
			CatID:     post.CatID,
			PostName:  post.PostName,
			PostImage: a,
		},
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/posts", http.StatusTemporaryRedirect)

}

func (h *Handler) listPost(rw http.ResponseWriter, r *http.Request) {

	res, err := h.pc.ListPost(r.Context(), &ppb.ListPostRequest{})
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(res)

	if err := h.templates.ExecuteTemplate(rw, "list-post.html", res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) editPost(rw http.ResponseWriter, r *http.Request) {
	cat, err := h.cc.List(r.Context(), &tpb.ListCategoryRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var category []Category
	for _, val := range cat.Category {
		category = append(category, Category{
			ID:           int(val.ID),
			CategoryName: val.CategoryName,
		})
	}
	vars := mux.Vars(r)
	Id := vars["id"]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	pl, err := h.pc.GetPost(r.Context(), &ppb.GetPostRequest{
		ID: id,
	})

	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	post := Post{
		ID: pl.Post.ID,
		CatID: pl.Post.CatID,
		PostName: pl.Post.PostName,
		PostImage: pl.Post.PostImage,

	}
	res := map[string]string{}

	h.loadEditPostForm(rw, post, category, res)
}
func (h *Handler) loadEditPostForm(rw http.ResponseWriter, post Post, category []Category, err map[string]string) {
	form := FormData{
		Post:     post,
		Category: category,
		Errors:   err,
	}

	if err := h.templates.ExecuteTemplate(rw, "edit-post.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) viewPost(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	Id := vars["id"]

	if Id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	id, err := strconv.ParseInt(Id, 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	pl, err := h.pc.GetPost(r.Context(), &ppb.GetPostRequest{
		ID: id,
	})

	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	post := Post{
		ID: pl.Post.ID,
		CatID: pl.Post.CatID,
		PostName: pl.Post.PostName,
		PostImage: pl.Post.PostImage,
		CategoryName: pl.Post.CategoryName,

	}
	res := map[string]string{}

	h.loadViewPostForm(rw, post, res)
}
func (h *Handler) loadViewPostForm(rw http.ResponseWriter, post Post, err map[string]string) {
	form := FormData{
		Post:     post,
		Errors:   err,
	}

	if err := h.templates.ExecuteTemplate(rw, "view-post.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) updatePost(rw http.ResponseWriter, r *http.Request) {
	cat, err := h.cc.List(r.Context(), &tpb.ListCategoryRequest{})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var category []Category
	for _, val := range cat.Category {
		category = append(category, Category{
			ID:           int(val.ID),
			CategoryName: val.CategoryName,
		})
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}
	Id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.pc.GetPost(r.Context(), &ppb.GetPostRequest{
		ID: Id,
	})

	if err != nil {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}

	var post Post
	

	if err := h.decoder.Decode(&post, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("PostImage")
    
	var imageName string
	
    if err == nil {
		defer file.Close()
		tempFile, err := ioutil.TempFile("cms/assets/images", "upload-*.png")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		
		tempFile.Write(fileBytes)
		
		imageName = tempFile.Name()

		if err := os.Remove(res.Post.PostImage); err != nil {
				http.Error(rw, "Unable to delete image", http.StatusInternalServerError)
				return
			}
	} else {
		imageName = res.Post.PostImage
	}

	if err := post.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range vErrors {
				vErrs[key] = value.Error()
			}
			h.loadEditPostForm(rw, post, category, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	_, err = h.pc.UpdatePost(r.Context(), &ppb.UpdatePostRequest{
		Post: &ppb.Post{
			ID:           Id,
			CatID:        post.CatID,
			PostName:     post.PostName,
			PostImage:    imageName,
		},
	})
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/posts", http.StatusTemporaryRedirect)
}

func (h Handler) deletePost(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	

	if id == "" {
		http.Error(rw, "invalid delete", http.StatusTemporaryRedirect)
		return
	}
	Id, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := h.pc.GetPost(r.Context(), &ppb.GetPostRequest{
		ID: Id,
	})

	if err != nil {
		http.Error(rw, "invalid URL", http.StatusInternalServerError)
		return
	}
	if err := os.Remove(res.Post.PostImage); err != nil {
		http.Error(rw, "Unable to delete image", http.StatusInternalServerError)
		return
	}

	_, err = h.pc.DeletePost(r.Context(), &ppb.DeletePostRequest{
		ID: Id,
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/posts", http.StatusTemporaryRedirect)
}

func (h *Handler) loadCreatedPostForm(rw http.ResponseWriter, post Post, category []Category, errs map[string]string) {
	form := FormData{
		Post: post,
		Category: category,
		Errors:   errs,
	}

	if err := h.templates.ExecuteTemplate(rw, "create-post.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
