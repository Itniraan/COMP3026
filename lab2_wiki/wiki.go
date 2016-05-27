package main

// Import needed packages (included the log package so that I could track down bugs)
import (
    "html/template"
    "io/ioutil"
    "net/http"
    "regexp"
    //"errors"
    "os"
    "log"
)

// Global variables for template and valid path format
var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))
// Valid path must be either /edit/, /view/ or /save/, followed by any alphanumeric characters
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// Page represents information that each page will contain
type Page struct {
    Title string
    Body []byte
}

/**
Save function - Saves pages to data directory as text files
*/
func (p *Page) save() error {
    // Create data directory if it doesn't already exist, give full access privileges
    os.Mkdir("data", 0777)
    fileName := "data/" + p.Title + ".txt"
    // Write text file
    return ioutil.WriteFile(fileName, p.Body, 0600)
}

/**
loadPage function - Loads page from text file in data directory
*/
func loadPage(title string) (*Page, error) {
    // Grab the requested file
    fileName := "data/" + title + ".txt"
    // Read file
    body, err := ioutil.ReadFile(fileName)
    log.Println(fileName)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

/**
viewHandler function - Handles viewing the wiki pages, if page doesn't exist, 
redirects to edit page
*/
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    log.Println("Viewing " + title)
    if err != nil {
        // If file doesn't exist, create and open it in edit mode
        http.Redirect(w, r, "/edit/"+title, http.StatusFound)
        return
    }
    // Output file to view template
    renderTemplate(w, "view", p)
}

/**
editHandler - Handles editing wiki pages.  If someone attempts to view a page 
that doesn't exist, they will be redirected here to create the page.
*/
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    log.Println("Editing " + title)
    if err != nil {
        p = &Page{Title: title}
    }
    // Output file to edit template
    renderTemplate(w, "edit", p)
}

/**
saveHandler function - Handles saving the new wiki page
*/
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    // Take title and body, save to text file
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Redirect to view the file that was just created
    http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

/**
// Don't believe this function is necessary anymore

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
    m := validPath.FindStringSubmatch(r.URL.Path)
    log.Println(m)
    if m == nil {
        http.NotFound(w, r)
        return "", errors.New("Invalid Page Title")
    }
    return m[2], nil // The title is the second subexpression
}*/

/**
makeHandler - Creates the handler using one of the three handler functions (view, edit, save)
Returned function is called a closure, because it encloses values defined outside of it. 
ie. fn will be either view, edit, or save handler.
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Check to make sure path matches validPath variable
        m := validPath.FindStringSubmatch(r.URL.Path)
        log.Println(m)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        fn(w, r, m[2])
    }
}

/**
renderTemplate function - Executes one of the valid templates (from the templates global variable)
*/
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    // Load requested template
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

/**
Main function - handles each of the three handlers (view, edit, save), and sets the webserver to listen on port 8000
*/
func main() {
    http.HandleFunc("/view/", makeHandler(viewHandler))
    http.HandleFunc("/edit/", makeHandler(editHandler))
    http.HandleFunc("/save/", makeHandler(saveHandler))
    // Listen on port 8000
    http.ListenAndServe(":8000", nil)
}

