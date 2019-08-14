package

query := `INSERT INTO
    		students (”firstName”,”lastName”)
		VALUES
    		($1, $2)
  		RETURNING id`
stmt, err := db.Prepare(query)
if err != nil {
log.Fatal(err)
}
defer stmt.Close()
var studentID int
err = stmt.QueryRow("Lee", "Provoost").Scan(&studentId)
if err != nil {
log.Fatal(err)
}
