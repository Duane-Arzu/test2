// file name: internal/data/comments.go
// package data

// import (
// 	"context"
// 	"database/sql"
// 	"errors"
// 	"time"

// 	"github.com/Duane-Arzu/comments/internal/validator"
// 	_ "github.com/Duane-Arzu/comments/internal/validator"
// )

<<<<<<< HEAD
// // each name begins with uppercase to make them exportable/ public
// type Comment struct {
// 	ID        int64     `json:"id"`      //unique value per comment
// 	Content   string    `json:"content"` //comment data
// 	Author    string    `json:"author"`  //person who wrote comment
// 	CreatedAt time.Time `json:"-"`       //database timestamp
// 	Version   int32     `json:"version"` //icremented on each update
// }

// // commentModel that expects a connection pool
// type CommentModel struct {
// 	DB *sql.DB
// }

// // Insert Row to comments table
// // expects a pointer to the actual comment content
// func (c CommentModel) Insert(comment *Comment) error {
// 	//the sql query to be executed against the database table
// 	query := `
// 	INSERT INTO comments (content, author)
// 	VALUES ($1, $2)
// 	RETURNING id, created_at, version`

// 	//the actual values to be passed into $1 and $2
// 	args := []any{comment.Content, comment.Author}

// 	// Create a context with a 3-second timeout. No database
// 	// operation should take more than 3 seconds or we will quit it
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	// execute the query against the comments database table. We ask for the the
// 	// id, created_at, and version to be sent back to us which we will use
// 	// to update the Comment struct later on
// 	return c.DB.QueryRowContext(ctx, query, args...).Scan(
// 		&comment.ID,
// 		&comment.CreatedAt,
// 		&comment.Version)

// }

// func ValidateComment(v *validator.Validator, comment *Comment) {
// 	//check if the content field is empty
// 	v.Check(comment.Content != "", "content", "must be provided")
// 	//check if the Author field is empty
// 	v.Check(comment.Author != "", "author", "must be provided")
// 	//check if the content field is empty
// 	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
// 	//check is author field is empty
// 	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")
// }

// // get a comment from DB based on ID
// func (c CommentModel) Get(id int64) (*Comment, error) {
// 	//check if the id is valid
// 	if id < 1 {
// 		return nil, ErrRecordNotFound
// 	}
// 	//the sql query to be excecuted against the database table
// 	query := `
// 	SELECT id, created_at, content, author, version
// 	FROM comments
// 	WHERE id = $1
// 	`

// 	//declare a variable of type Comment to hold the returned values
// 	var comment Comment

// 	//set 3-second context/timer
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	err := c.DB.QueryRowContext(ctx, query, id).Scan(
// 		&comment.ID,
// 		&comment.Content,
// 		&comment.Content,
// 		&comment.Author,
// 		&comment.Version,
// 	)
// 	//check for errors
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrRecordNotFound
// 		default:
// 			return nil, err
// 		}
// 	}
// 	return &comment, nil
// }

// func (c CommentModel) GetAll(content string, author string, filters Filters) (*[]Comment, error) {
// 	query := `
// 	SELECT id, created_at, content, author, version
// 	FROM comments
// 	WHERE (to_tsvector('simple',content) @@
// 		plainto_tsquery('simple', $1) OR $1 = '')
// 	AND (to_tsvector('simple',author) @@
// 		plainto_tsquery('simple',$2) OR $2 = '')
// 	ORDER BY id
// 	LIMIT $3 OFFSET $4
// 	`
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	// rows, err := c.DB.QueryContext(ctx, query, content, author)

// 	comments, err := c.DB.QueryContext(ctx, query, content, author, filters.limit(), filters.offset())
// 	//check for errors
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, sql.ErrNoRows):
// 			return nil, ErrRecordNotFound
// 		default:
// 			return nil, err
// 		}
// 	}

// 	var cmts []Comment
// 	for comments.Next() {
// 		var com Comment
// 		if err := comments.Scan(&com.ID, &com.CreatedAt, &com.Content, &com.Author, &com.Version); err != nil {
// 			return nil, err
// 		}
// 		cmts = append(cmts, com)
// 	}
// 	return &cmts, nil
// }

// // update  a specific record from the comments table
// func (c CommentModel) Update(comment *Comment) error {
// 	//the sql query to be excecuted against the DB table
// 	//Every time make an update, version number is incremented

// 	query := `
// 	UPDATE comments
// 	SET content=$1, author=$2, version=version+1
// 	WHERE id = $3
// 	RETURNING version
// 	`

// 	args := []any{comment.Content, comment.Author, comment.ID}
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
// 	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
=======
// each name begins with uppercase so that they are exportable/public
type Comment struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

type CommentModel struct {
	DB *sql.DB
}

func ValidateComment(v *validator.Validator, comment *Comment) {
	// check if the Content field is empty
	v.Check(comment.Content != "", "content", "must be provided")
	// check if the Author field is empty
	v.Check(comment.Author != "", "author", "must be provided")
	// check if the Content field is empty
	v.Check(len(comment.Content) <= 100, "content", "must not be more than 100 bytes long")
	// check if the Author field is empty
	v.Check(len(comment.Author) <= 25, "author", "must not be more than 25 bytes long")
}

func (c CommentModel) Insert(comment *Comment) error {
	// the SQL query to be executed against the database table
	query := `
		 INSERT INTO comments (content, author)
		 VALUES ($1, $2)
		 RETURNING id, created_at, version
		 `
	// the actual values to replace $1, and $2
	args := []any{comment.Content, comment.Author}

	// Create a context with a 3-second timeout. No database
	// operation should take more than 3 seconds or we will quit it
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// execute the query against the comments database table. We ask for the the
	// id, created_at, and version to be sent back to us which we will use
	// to update the Comment struct later on
	return c.DB.QueryRowContext(ctx, query, args...).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.Version)
}

// Get a specific Comment from the comments table
func (c CommentModel) Get(id int64) (*Comment, error) {
	// check if the id is valid
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// the SQL query to be executed against the database table
	query := `
		 SELECT id, created_at, content, author, version
		 FROM comments
		 WHERE id = $1
	   `
	// declare a variable of type Comment to store the returned comment
	var comment Comment

	// Set a 3-second context/timer
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.Content,
		&comment.Author,
		&comment.Version,
	)
	// Cont'd on the next slide
	// check for which type of error
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &comment, nil
}

func (c CommentModel) Update(comment *Comment) error {
	// The SQL query to be executed against the database table
	// Every time we make an update, we increment the version number
	query := `
			UPDATE comments
			SET content = $1, author = $2, version = version + 1
			WHERE id = $3
			RETURNING version 
			`

	args := []any{comment.Content, comment.Author, comment.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return c.DB.QueryRowContext(ctx, query, args...).Scan(&comment.Version)
>>>>>>> 0cc1270f48216b9318fcc1ef24b827397488e322

// }

<<<<<<< HEAD
// // delete a specific comment form the comments table
// func (c CommentModel) Delete(id int64) error {
// 	//check if the id is valid
// 	if id < 1 {
// 		return ErrRecordNotFound
// 	}

// 	//sql querry to be excecuted against the database table
// 	query := `
// 	DELETE FROM comments
// 	WHERE id = $1
// 	`

// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()
=======
func (c CommentModel) Delete(id int64) error {

	// check if the id is valid
	if id < 1 {
		return ErrRecordNotFound
	}
	// the SQL query to be executed against the database table
	query := `
        DELETE FROM comments
        WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
>>>>>>> 0cc1270f48216b9318fcc1ef24b827397488e322

// 	// ExecContext does not return any rows unlike QueryRowContext.
// 	// It only returns  information about the the query execution
// 	// such as how many rows were affected
// 	result, err := c.DB.ExecContext(ctx, query, id)
// 	if err != nil {
// 		return err
// 	}

<<<<<<< HEAD
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	//maybe wrong id for record was given so we sort of try checking
// 	if rowsAffected == 0 {
// 		return ErrRecordNotFound
// 	}
// 	return nil
// }
=======
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Probably a wrong id was provided or the client is trying to
	// delete an already deleted comment
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil

}

func (c CommentModel) GetAll(content string, author string, filters Filters) ([]*Comment, Metadata, error) {

	// the SQL query to be executed against the database table
	query := fmt.Sprintf(`
	SELECT COUNT(*) OVER(), id, created_at, content, author, version
	FROM comments
	WHERE (to_tsvector('simple', content) @@
		  plainto_tsquery('simple', $1) OR $1 = '') 
	AND (to_tsvector('simple', author) @@ 
		 plainto_tsquery('simple', $2) OR $2 = '') 
	ORDER BY %s %s, id ASC 
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := c.DB.QueryContext(ctx, query, content, author, filters.limit(), filters.offset())

	if err != nil {
		return nil, Metadata{}, err
	}

	// clean up the memory that was used
	defer rows.Close()
	totalRecords := 0
	// we will store the address of each comment in our slice
	comments := []*Comment{}

	// process each row that is in rows

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&totalRecords,
			&comment.ID,
			&comment.CreatedAt,
			&comment.Content,
			&comment.Author,
			&comment.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// add the row to our slice
		comments = append(comments, &comment)
	} // end of for loop

	// after we exit the loop we need to check if it generated any errors
	err = rows.Err()
	if err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)

	return comments, metadata, nil

}
>>>>>>> 0cc1270f48216b9318fcc1ef24b827397488e322
