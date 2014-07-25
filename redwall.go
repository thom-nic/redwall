package redwall

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"log"
	"net/http"
	"path"
)

type Item struct {
	Title  string
	URL    string
	Score  int
	Over18 bool `json:over_18`
}

type Response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

func GetSubreddit(reddit string) ([]Item, error) {

	url := fmt.Sprintf("http://reddit.com/r/%s.json", reddit)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	resp := new(Response)
	err = json.NewDecoder(r.Body).Decode(resp)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)

	items := make([]Item, len(resp.Data.Children))
	for i, child := range resp.Data.Children {
		items[i] = child.Data
	}
	return items, nil
}

func (i Item) Html() string {

	com := ""
	switch i.Score {
	case 0:
		// nothing
	case 1:
		com = " Score: 1"
	default:
		com = fmt.Sprintf(" (Score: %d)", i.Score)
	}
	return fmt.Sprintf("<p>%s<b>%s</b><br/> <a href=\"%s\">%s</a></p>", i.Title, com, i.URL, i.URL)
}

func Download(items []Item, destDir string, err error) (string, error) {
	if err != nil {
		return "", err
	}

	for i, item := range items {
		if item.Over18 {
			continue
		}
		tokens := strings.Split(item.URL, "/")
		fileName := tokens[len(tokens)-1]
		// TODO cleanup characters that occur after ? or # (or use a URL parsing lib)
		fmt.Println("Downloading", item.Title, "to", fileName)

		// TODO: check file existence first with io.IsExist
		// base file name on permalink
		if strings.HasSuffix(item.URL, ".jpg") || strings.HasSuffix(item.URL, ".png") {
			outPath = path.join(destDir, fileName)
			output, err := os.Create(outPath)
			if err != nil {   q
				log.Fatal("Can't open output file", outPath)
				continue
			}
			defer output.Close()

			response, err := http.Get(url)
			if err != nil {
				fmt.Println("Error while downloading", url, "-", err)
				return
			}
			defer response.Body.Close()

			n, err := io.Copy(outPath, response.Body)
			if err != nil {
				fmt.Println("Error while downloading", url, "-", err)
				return
			}
		}
	}

}

func Email(items []Item, err error) (string, error) {

	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	// Need to build strings from items
	for _, item := range items {
		buffer.WriteString(item.Html())
	}

	return buffer.String(), nil
}

func Send(html string, err error) (interface{}, error) {

	if err != nil {
		return html, err
	}

	sg := sendgrid.NewSendGridClient("sendgrid_user", "sendgrid_key")
	message := sendgrid.NewMail()

	message.AddTo("myemail@me.com")
	message.AddToName("Robin Johnson")
	message.SetSubject("Your Daily Golang Breakfast News!")
	message.SetFrom("rbin@sendgrid.com")

	message.SetHTML(html)
	return sg.Send(message), nil
}

func main() {

	//	rep, err := Send(Email(GetSubreddit("wallpapers")))
	rep, err := Download(GetSubreddit("wallpapers"))

	if err != nil {
		log.Fatal(err)
	}

	if rep == nil {
		fmt.Println("Email sent!")
		fmt.Println("Closing...")
	} else {
		fmt.Println(rep)
	}

}
