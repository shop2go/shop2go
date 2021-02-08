package handler

import (
	"bufio"
	"encoding/base64"
	//"encoding/json"
	"fmt"
	"os"

	fb "github.com/huandu/facebook"

	/* "golang.org/x/oauth2"
	oauth2fb "golang.org/x/oauth2/facebook" */
	"io/ioutil"
	"net/http"
	"strconv"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "GET":

		str := `
	<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />

    <title>img upload</title>
  </head>
  <body>
    <form
      enctype="multipart/form-data"
      action="https://shop2go.cloud/api/img"
      method="POST"
    >
      <input type="file" name="filer" />
      <input type="submit" value="upload" />
    </form>
  </body>
</html>
`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

	case "POST":

		file, _, err := r.FormFile("filer")

		if err != nil {
			fmt.Fprint(w, err)
		}

		/* 		fmt.Fprintf(w, "<br>Uploaded File: %+v", handler.Filename)
		   		fmt.Fprintf(w, "<br>File Size: %+v", handler.Size)
		   		fmt.Fprintf(w, "<br>MIME Header: %+v", handler.Header) */

		reader := bufio.NewReader(file)

		file.Close()

		content, err := ioutil.ReadAll(reader)

		if err != nil {
			fmt.Fprint(w, err)
		} else if content == nil {
			fmt.Fprint(w, "ERROR")
		}

		encoded := base64.StdEncoding.EncodeToString(content)

		str := `
	<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />

    <title>img</title>

  </head>
  <body>
  <br>
  <img id=image src="data:image/png;base64,` + encoded + `" />
  </body>
</html>
`

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(str)))
		w.Write([]byte(str))

		// Create a global App var to hold app id and secret.
		var globalApp = fb.New("251435286506299", os.Getenv("APP_SECRET"))

		// Facebook asks for a valid redirect URI when parsing the signed request.
		// It's a newly enforced policy starting as of late 2013.
		globalApp.RedirectUri = "https://code2go.dev/"

		// If there is another way to get decoded access token,
		// this will return a session created directly from the token.
		session := globalApp.Session(os.Getenv("FB_TOKEN"))

		// This validates the access token by ensuring that the current user ID is properly returned. err is nil if the token is valid.
		err = session.Validate()
		if err != nil {
			fmt.Fprint(w, err)
		}

		/* 		_, err = session.Post("/123190199635156/photos", fb.Params{
			"caption": "img upload",
			"url":     "https://shop2go.cloud/api/img",
		}) */

		_, err = session.Post("/123190199635156/photos", fb.Params{
			"data": encoded,
		})

		if err != nil {
			fmt.Fprint(w, err)
		}

		/* resp, err := http.Post("https://graph.facebook.com/{123190199635156}/photos?url={362047903.index.png}&access_token={EAADkrdbvqzsBAIB2glZCJJIeoZBZAHGhe3f0ZBdWRDEiFG2VGXUPV0tBX3L450FsQ9gHeGAQSsga9MUrB6U7EoElq4Pvm5yg5CJl2tpe4PRknO7UjXZAgkOwsCMnRjJgAlwZCSGnkTcUGJMbwJrZArgrhLdEbL3bxZCWCGVxHn6RxQZDZD}", "", nil)
		if err != nil {
			fmt.Println(err)
		}

		bytesBody, _ := ioutil.ReadAll(resp.Body)

		defer resp.Body.Close()

		log.Println(string(bytesBody))

		if err = json.Unmarshal(bytesBody, &data); err != nil {
			log.Println(err)
		}

		tempFile.Close()

		fmt.Fprint(w, data.Post) */

	}

}
