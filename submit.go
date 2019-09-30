/* Submit client data http handler. */

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
    "time"
)

func submitAjax(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("Received ajax data.")

	// Invoke ParseForm before reading form values
	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("%s: %s\n", key, value)
	}

	// Acknowledge receiving the data
	fmt.Printf("User name: %s\n", r.FormValue("username"))
	fmt.Printf("Message body: %s\n", r.FormValue("message"))

	// Format current time
	now := time.Now()
	//fmt.Println(now.Format("Jan 2, 2006 - 3:04 pm MST"))
	var now_time string = now.Format("Jan 2, 2006 - 3:04 pm MST")

	// Client data received
	user_name := r.FormValue("username")
	user_message := r.FormValue("message")
	
	// Insert into database
	statement, _ := database.Prepare("INSERT INTO comments (username, comment, date) VALUES (?, ?, ?)")
	statement.Exec(user_name, user_message, now_time)

	// Read database
	rows, _ := database.Query("SELECT id, username, comment, date FROM comments ORDER BY id DESC LIMIT 0, 10")

	var db_id int
	var db_name string
	var db_message string
	var db_timestamp string
	var user_comment UserComment
	var user_comments []UserComment

	for rows.Next() {
		rows.Scan(&db_id, &db_name, &db_message, &db_timestamp)
	
		user_comment.Id = db_id
		user_comment.Name = db_name
		user_comment.Message = db_message
		user_comment.Timestamp = db_timestamp
		user_comments = append(user_comments, user_comment)
	}

	defer rows.Close()

	json_str, _ := json.Marshal(user_comments)
	fmt.Printf("%s\n", json_str)
	
	// Send back data to client
	w.Write(json_str)
 }