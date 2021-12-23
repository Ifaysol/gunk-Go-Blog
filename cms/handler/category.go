package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	tpb "grpc-blog/gunk/v1/categories"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

type Category struct {
	ID                  int    `db:"id"`
	CategoryName        string `db:"categoryname"`
	CategoryDescription string `db:"categorydescription"`
	Errors               map[string]string
}

type ListCategory struct {
	Categories []Category
}

func (c *Category) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.CategoryName, validation.Required.Error("Must insert a category name"), validation.Length(3, 0)),
		validation.Field(&c.CategoryDescription, validation.Required.Error("Must insert a category description"), validation.Length(3, 0)),
	)
}

type FData struct {
	Category Category
	Errors   map[string]string
}

func (h *Handler) createCategory(rw http.ResponseWriter, r *http.Request) {
	vErrs := map[string]string{}
	category := Category{}
	h.loadCreatedCategoryForm(rw, category, vErrs)

}

func (h *Handler) storeCategory(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := category.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[strings.Title(key)] = value.Error()

			}

			h.loadCreatedCategoryForm(rw, category, vErr)
			return

		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return

	}

	_, err := h.cc.Create(r.Context(), &tpb.CreateCategoryRequest{
		Category: &tpb.Category{
			CategoryName:        category.CategoryName,
			CategoryDescription: category.CategoryDescription,
		},
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// const insertCategory = `INSERT INTO categories(categoryname) VALUES($1);`
	// res := h.db.MustExec(insertCategory, category.CategoryName)

	// if ok, err := res.RowsAffected(); err != nil || ok == 0 {
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)

}

func (h *Handler) listCategory(rw http.ResponseWriter, r *http.Request) {

	res, err := h.cc.List(r.Context(), &tpb.ListCategoryRequest{})
	if err != nil {
		log.Fatal(err)
	}
	// log.Println(res)

	if err := h.templates.ExecuteTemplate(rw, "list-category.html", res); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) editCategory(rw http.ResponseWriter, r *http.Request) {
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

	catl, err := h.cc.Get(r.Context(), &tpb.GetCategoryRequest{
		ID: id,
	})

	if err != nil {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	res := map[string]string{}

	h.loadEditCategoryForm(rw, int(id), catl.Category.CategoryName, catl.Category.CategoryDescription, res)
}
func (h *Handler) loadEditCategoryForm(rw http.ResponseWriter, id int, name string, des string, err map[string]string) {
	form := Category{
		ID: id,
		CategoryName:        name,
		CategoryDescription: des,
		Errors:              err,
	}
	// fmt.Println(form)

	if err := h.templates.ExecuteTemplate(rw, "edit-category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h Handler) updateCategory(rw http.ResponseWriter, r *http.Request) {
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

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	

	if err := category.Validate(); err != nil {
		valError, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range valError {
				vErrs[key] =value.Error()
			}
			h.loadEditCategoryForm(rw, int(id), category.CategoryName, category.CategoryDescription, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	_, err = h.cc.Update(r.Context(), &tpb.UpdateCategoryRequest{
		Category: &tpb.Category{
			ID: id,
			CategoryName: category.CategoryName,
			CategoryDescription: category.CategoryDescription,
		},
	})
	
	if err != nil{
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)
}

func (h Handler) deleteCategory(rw http.ResponseWriter, r *http.Request) {
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

	_, err = h.cc.Delete(r.Context(), &tpb.DeleteCategoryRequest{
		ID: Id,
	})

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)
}

func (h *Handler) loadCreatedCategoryForm(rw http.ResponseWriter, category Category, errs map[string]string) {
	form := FData{
		Category: category,
		Errors:   errs,
	}

	if err := h.templates.ExecuteTemplate(rw, "create-category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
