package main

import (
	"context"
	//"encoding/base64"
	"fmt"
	"net/http"
	"os"
	// 	"sort"
	"strconv"
	"strings"

	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	//"github.com/plutov/paypal"
)

type Access struct {
	//Reference *f.RefV `fauna:"ref"`
	Timestamp int    `fauna:"ts"`
	Secret    string `fauna:"secret"`
	Role      string `fauna:"role"`
}

type UserEntry struct {
	ID       graphql.ID     `graphql:"_id"`
	Username graphql.String `graphql:"username"`
}

func Handler(w http.ResponseWriter, r *http.Request) {

	var costumer, content string

	//replace with token
	c, err := r.Cookie("shop2go.cloud")

	if err != nil {

		c.Value = ""

	}

	id := r.Host

	id = strings.TrimSuffix(id, "shop2go.cloud")

	id = strings.TrimSuffix(id, ".")

	fc := f.NewFaunaClient(os.Getenv("FAUNA_ACCESS"))

	x, err := fc.Query(f.CreateKey(f.Obj{"database": f.Database("access"), "role": "server"}))

	if err != nil {

		fmt.Fprintf(w, "connection error: %v\n", err)

	}

	var access *Access

	x.Get(&access)

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: access.Secret},
	)

	httpClient := oauth2.NewClient(context.Background(), src)

	call := graphql.NewClient("https://graphql.fauna.com/graphql", httpClient)

	switch id {

	case "":

		if c.Value != "" {

			var q struct {
				FindUserByID struct {
					UserEntry
				} `graphql:"findCartByID(id: $ID)"`
			}

			v := map[string]interface{}{
				"ID": graphql.ID(c.Value),
			}

			if err = call.Query(context.Background(), &q, v); err != nil {
				fmt.Fprintf(w, "error with user: %v\n", err)
			}

			costumer = fmt.Sprintf("%s", q.FindUserByID.ID)

		}

	default:

		var q struct {
			FindUserByID struct {
				UserEntry
			} `graphql:"findCartByID(id: $ID)"`
		}

		v := map[string]interface{}{
			"ID": graphql.ID(id),
		}

		if err = call.Query(context.Background(), &q, v); err != nil {
			fmt.Fprintf(w, "error with user: %v\n", err)
		}

		costumer = fmt.Sprintf("%s", q.FindUserByID.ID)

	}

	content = `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>SHOP2GO</title>
		<link href="https://fonts.googleapis.com/css?family=Roboto:300,300i,400,400i,500,500i,700,700i|Roboto+Mono:300,400,700|Roboto+Slab:300,400,700" rel="stylesheet">
		<link href="https://fonts.googleapis.com/icon?family=Material+Icons" rel="stylesheet">
   		<link href="https://assets.medienwerk.now.sh/material.min.css" rel="stylesheet">
		</head>
		<body style="background-color: #bcbcbc;">
		<br>
		<br>
		<div class="container" id="data" style="color:white;">
		</div>
		<script src="https://assets.medienwerk.now.sh/material.min.js">
		</script>
		</body>
		</html>
		`

	switch r.Method {

	case "GET":

		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		w.Write([]byte(content))
	}

}
