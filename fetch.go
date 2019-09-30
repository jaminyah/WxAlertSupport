/* Fetch comments from database handler. */

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func fetchComments(w http.ResponseWriter, r *http.Request) {
	

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