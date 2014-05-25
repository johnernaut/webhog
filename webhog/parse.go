// webhog is a package that stores and downloads a
// given URL (including js, css, and images) for offline
// use and uploads it to a given AWS-S3 account.
package webhog

import (
	"bytes"
	"code.google.com/p/go.net/html"
	uuid "github.com/nu7hatch/gouuid"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
)

// Regex to match css and js extentions.
var rxExt = regexp.MustCompile(`(\.(?:css|js|gif|png|jpg))\/?$`)

// HTML nodes and their respective attrs that we need.
var matchVals = map[string]string{"link": "href", "script": "src", "img": "src"}

// Start the scraping process.
func NewScraper(url string) (*Entity, error) {
	entity := new(Entity)

	// Return existing entity if it exists, otherwise create a
	// new one.
	err := checkExistingEntity(url, entity)
	if err != nil {
		if err = createNewEntity(url, entity); err != nil {
			log.Panicln("Error creating new entity: ", err)
		}
	}

	// Create a new entity.

	return entity, err
}

// Make a GET request to the given URL and start parsing
// its HTML.
func ExtractData(entity *Entity, url string) {
	// Parsing completion channel.
	done := make(chan bool, 1)

	res, err := http.Get(url)
	if err != nil {
		log.Panicln("Error requesting URL data: ", err)
	}

	defer res.Body.Close()

	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Println("Error parsing URL body: ", err)
	}

	go ParseHTML(doc, entity, done)

	for {
		select {
		case <-done:
			go finalizeEntity(entity, doc, EntityDir)
		default:
		}
	}
}

// Parse the HTML - pull the href/src attributes for js,
// css, and images for download.
func ParseHTML(n *html.Node, entity *Entity, done chan bool) {
	var wg sync.WaitGroup

	wg.Add(1)
	go extractAttrs(n, entity, &wg)
	wg.Wait()
	done <- true
}

// Extract attributes on elements that match css/js/link tags.
func extractAttrs(n *html.Node, entity *Entity, wg *sync.WaitGroup) {
	defer wg.Done()

	// loop through file types and their extensions to check for matches
	for j, p := range matchVals {
		if n.Type == html.ElementNode && n.Data == j {
			for i := range n.Attr {
				attr := &n.Attr[i]
				if attr.Key == p {
					matchAttrs(attr, entity)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		wg.Add(1)
		go extractAttrs(c, entity, wg)
	}
}

// match found attr file types and check src/href for a
// fullpath URL
func matchAttrs(attr *html.Attribute, entity *Entity) {
	// match css and js extentions
	match := rxExt.FindStringSubmatch(attr.Val)

	// check if a full URL is given to the entity, otherwise
	// we need to create one for the resource.
	matched, err := regexp.MatchString("http.*", attr.Val)
	if err != nil {
		log.Println("Error matching regex http(s): ", err)
	}

	// valid filetype and there is a full URL we can use
	if len(match) > 0 && matched {
		// Download and persist the current resource and insert it's
		// new name in the value of the HTML tree.
		name, err := StoreResource(attr.Val, string(match[0]), EntityDir)
		if err != nil {
			log.Println("Error storing resource: ", err)
		}
		attr.Val = name
	}

	// valid filetype but there is a relative URL
	if len(match) > 0 && !matched {
		// new resource name after adding in full URL
		var updName string
		// check for trailing slash on the entities' URL
		ln := len(entity.Url)
		if string(entity.Url[ln-1]) == "/" {
			updName = entity.Url + attr.Val
		} else {
			updName = entity.Url + "/" + attr.Val
		}

		name, err := StoreResource(updName, string(match[0]), EntityDir)
		if err != nil {
			log.Println("Error storing resource: ", err)
		}
		attr.Val = name
	}
}

// See if this URL has already been saved into the
// database and isn't expired - if so, return it.
func checkExistingEntity(url string, e *Entity) error {
	err := Find(e, bson.M{"url": url}).One(&e)

	exp := e.CreatedAt.Add(ExpirationTime)

	// update the entity and re-download / parse assets
	if time.Now().After(exp) {
		err = Update(e, bson.M{"uuid": e.UUID},
			bson.M{"$set": bson.M{"aws_link": "", "status": ParsingStatus, "created_at": time.Now()}})
		if err != nil {
			log.Println("Error updating entity: ", err)
		}

		go ExtractData(e, url)
	}

	return err
}

// Create a new entity to persist into the database - start HTML extraction.
func createNewEntity(url string, entity *Entity) error {
	err := NewEntityDir()
	if err != nil {
		log.Println("Error creating entity dir: ", err)
	}

	go ExtractData(entity, url)

	id, err := uuid.NewV4()
	if err != nil {
		log.Println("Error creating UUID: ", err)
	}

	// Set new entities' initial data from the request.
	entity.Status = ParsingStatus
	entity.Url = url
	entity.UUID = id.String()
	entity.CreatedAt = time.Now()

	// Persist new entity into the database.
	err = Create(entity)
	if err != nil {
		return err
	}

	err = Find(entity, bson.M{"uuid": entity.UUID}).One(&entity)
	if err != nil {
		return err
	}

	return err
}

// Archive / upload final entity HTML
func finalizeEntity(entity *Entity, doc *html.Node, entDir string) {
	log.Println("Entity tmp directory: ", entDir)

	var finalHTML bytes.Buffer
	bl := html.Render(&finalHTML, doc)
	if bl != nil {
		log.Println(bl)
	}

	err := StoreHTML(finalHTML, entDir)
	if err != nil {
		log.Println("Error in StoreHTML: ", err)
	}

	zipName, err := ArchiveFinalFiles(entDir)
	if err != nil {
		log.Println("Error archiving files: ", err)
	}

	err = Update(entity, bson.M{"uuid": entity.UUID}, bson.M{"$set": bson.M{"status": UploadingStatus}})
	if err != nil {
		log.Println("Error updating entity: ", err)
	}

	awsLink, err := UploadEntity(zipName, entity)
	if err != nil {
		log.Println("Error uploading final files: ", err)
	}

	err = Update(entity, bson.M{"uuid": entity.UUID}, bson.M{"$set": bson.M{"aws_link": awsLink, "status": CompleteStatus}})
	if err != nil {
		log.Println("Error updating entity: ", err)
	}
}
