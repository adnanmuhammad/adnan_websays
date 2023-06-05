package main

//============ Import the Required Packages for the Project ============
import (
    "github.com/gorilla/mux"
    "database/sql"
    _"github.com/go-sql-driver/mysql"
    "os"
    "strconv"
    "log"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "fmt"
)

//============  Product Entity ==========
type Product struct {
  ProdId string `json:"prod_id"`
  ProdTitle string `json:"prod_title"`
}

//============ Article Entity ============
type Article struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
    Author string `json:"author"`
    Country string `json:"country"`
}
var articles []Article

//============ Category Entity ============
type Category struct {
  ID   int    `json:"id"`
  Name string `json:"name"`
}

type Categories struct {
  Data []Category `json:"data"`
}

var categories Categories


var db *sql.DB
var err error


func main() {
  db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_websays")
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()
  router := mux.NewRouter()

  // ================ CRUD for products using a database MYSQL =================
  router.HandleFunc("/products", getProducts).Methods("GET")
  router.HandleFunc("/products", createProduct).Methods("POST")
  router.HandleFunc("/products/{prod_id}", getProduct).Methods("GET")
  router.HandleFunc("/products/{prod_id}", updateProduct).Methods("PUT")
  router.HandleFunc("/products/{prod_id}", deleteProduct).Methods("DELETE")

  // ================ CRUD for Articles entity with memory persistence ================
  router.HandleFunc("/articles", getArticles).Methods("GET")
  router.HandleFunc("/articles", createArticle).Methods("POST")
  router.HandleFunc("/articles/{id}", getArticle).Methods("GET")
  router.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")
  router.HandleFunc("/articles/{id}", deleteArticle).Methods("DELETE")

  // ================ CRUD for Categories entity with text file persistence ================
  loadCategories()
  router.HandleFunc("/categories", getCategories).Methods("GET")
  router.HandleFunc("/categories/{id}", getCategory).Methods("GET")
  router.HandleFunc("/categories", createCategory).Methods("POST")
  router.HandleFunc("/categories/{id}", updateCategory).Methods("PUT")
  router.HandleFunc("/categories/{id}", deleteCategory).Methods("DELETE")

  log.Fatal(http.ListenAndServe(":1337", router))

}


// ==================== CRUD with Database MYSQL for Products ====================
func getProducts(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var products []Product
  result, err := db.Query("SELECT prod_id, prod_title from tbl_products")
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()
  for result.Next() {
    var product Product
    err := result.Scan(&product.ProdId, &product.ProdTitle)
    if err != nil {
      panic(err.Error())
    }
    products = append(products, product)
  }
  json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  stmt, err := db.Prepare("INSERT INTO tbl_products(prod_title) VALUES(?)")
  if err != nil {
    panic(err.Error())
  }
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err.Error())
  }
  keyVal := make(map[string]string)
  json.Unmarshal(body, &keyVal)
  prod_title := keyVal["prod_title"]
  _, err = stmt.Exec(prod_title)
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "New Product Added Successfully")
}

func getProduct(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  result, err := db.Query("SELECT prod_id, prod_title FROM tbl_products WHERE prod_id = ?", params["prod_id"])
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()
  var product Product
  for result.Next() {
    err := result.Scan(&product.ProdId, &product.ProdTitle)
    if err != nil {
      panic(err.Error())
    }
  }
  json.NewEncoder(w).Encode(product)
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  stmt, err := db.Prepare("UPDATE tbl_products SET prod_title = ? WHERE prod_id = ?")
  if err != nil {
    panic(err.Error())
  }
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err.Error())
  }
  keyVal := make(map[string]string)
  json.Unmarshal(body, &keyVal)
  prodTitleNew := keyVal["prod_title"]
  _, err = stmt.Exec(prodTitleNew, params["prod_id"])
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "Product with ProdId = %s Updated Successfully", params["prod_id"])
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  stmt, err := db.Prepare("DELETE FROM tbl_products WHERE prod_id = ?")
  if err != nil {
    panic(err.Error())
  }
  _, err = stmt.Exec(params["prod_id"])
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "Product with ProdId = %s deleted Successfully", params["prod_id"])
}

// ==================== CRUD with Database MYSQL for Products ====================




// ===================== Articles CRUD with memory persistence =================

func getArticles(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(articles)
}

func createArticle(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var article Article
    json.NewDecoder(r.Body).Decode(&article)
    articles = append(articles, article)
    json.NewEncoder(w).Encode(article)
}

func getArticle(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id := vars["id"]
    for _, article := range articles {
        if article.ID == id {
            json.NewEncoder(w).Encode(article)
            return
        }
    }
    http.NotFound(w, r)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id := vars["id"]
    var updatedArticle Article
    json.NewDecoder(r.Body).Decode(&updatedArticle)
    for i, article := range articles {
        if article.ID == id {
            articles[i] = updatedArticle
            json.NewEncoder(w).Encode(updatedArticle)
            return
        }
    }
    http.NotFound(w, r)
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    vars := mux.Vars(r)
    id := vars["id"]
    for i, article := range articles {
        if article.ID == id {
            articles = append(articles[:i], articles[i+1:]...)
            json.NewEncoder(w).Encode(article)
            return
        }
    }
    http.NotFound(w, r)
}

// ===================== Articles CRUD with memory persistence =================




// ===================== Categories with text file persistence =================

func loadCategories() {
  file, err := os.Open("text_file_categories.txt")
  if err != nil {
    log.Println("Error loading categories:", err)
    return
  }
  defer file.Close()

  bytes, err := ioutil.ReadAll(file)
  if err != nil {
    log.Println("Error loading categories:", err)
    return
  }

  if err := json.Unmarshal(bytes, &categories); err != nil {
    log.Println("Error unmarshalling categories:", err)
  }

}

func saveCategories() {
  bytes, err := json.Marshal(categories)
  if err != nil {
    log.Println("Error marshalling categories:", err)
    return
  }

  if err := ioutil.WriteFile("text_file_categories.txt", bytes, 0644); err != nil {
    log.Println("Error saving categories:", err)
  }
}

func getNextID() int {
  if len(categories.Data) == 0 {
    return 1
  }

  lastCategory := categories.Data[len(categories.Data)-1]
  return lastCategory.ID + 1
}

func getCategories(w http.ResponseWriter, r *http.Request) {
  json.NewEncoder(w).Encode(categories)
}

func getCategory(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    http.Error(w, "Invalid category ID", http.StatusBadRequest)
    return
  }

  for _, category := range categories.Data {
    if category.ID == id {
      json.NewEncoder(w).Encode(category)
      return
    }
  }

  http.Error(w, "Category not found", http.StatusNotFound)
}

func createCategory(w http.ResponseWriter, r *http.Request) {
  var category Category
  if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  category.ID = getNextID()
  categories.Data = append(categories.Data, category)
  saveCategories()

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(category)
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    http.Error(w, "Invalid category ID", http.StatusBadRequest)
    return
  }

  var updatedCategory Category
  if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  for i, category := range categories.Data {
    if category.ID == id {
      updatedCategory.ID = id
      categories.Data[i] = updatedCategory
      saveCategories()

      w.Header().Set("Content-Type", "application/json")
      json.NewEncoder(w).Encode(updatedCategory)
      return
    }
  }

  http.Error(w, "Category not found", http.StatusNotFound)
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    http.Error(w, "Invalid category ID", http.StatusBadRequest)
    return
  }

  for i, category := range categories.Data {
    if category.ID == id {
      // Remove the category from the slice
      categories.Data = append(categories.Data[:i], categories.Data[i+1:]...)
      saveCategories()

      w.WriteHeader(http.StatusNoContent)
      return
    }
  }

  http.Error(w, "Category not found", http.StatusNotFound)
}
// ===================== Categories with text file persistence =================