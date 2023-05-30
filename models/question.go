package models

type Question struct {
	ID           int
	SessionID    int
	UserID       int
	Text         string
	Votes        int
	Answered     bool
	CreatedAt    string
	UpdatedAt    string
	UserHasVoted bool
}
type Vote struct {
	ID         int
	QuestionID int
	UserID     int
}

// ** CREATE **
// PostQuestion posts a question to the database for a given class session and user.
// It returns the ID of the new question and any error encountered.
func (s Question) PostQuestion(userID int, sessionID int, text string) (int, error) {
	db := NewDB()

	sqlStatement := `
	INSERT INTO
		question
	VALUES 
		(NULL,
		$1,
		$2,
		$3,
		0, 
		false, 
		datetime('now'), 
		datetime('now')
		)
	RETURNING id
	`
	var questionID int
	err := db.QueryRow(sqlStatement, sessionID, userID, text).Scan(&questionID)
	if err != nil {
		return 0, err
	}

	return questionID, nil
}

// ** READ **
// GetAllQuestions takes a class session id, a sortBy parameter, and a timezone offset.
// It returns all questions for a given class session as a slice of Question
// structs and any error encountered.
func (s *Question) GetAllQuestions(classSessionID int, sortBy string, timezoneOffset int) ([]Question, error) {
	db := NewDB()

	// Get all questions for this session
	sqlStatement := `
        SELECT
            id, session_id, user_id, text, votes, answered,
			strftime('%H:%M', datetime(created_at, (? || ' minutes'))) as formatted_created_at,
			datetime(updated_at, (? || ' minutes'))
        FROM
            question
        WHERE
            session_id = ?
        ORDER BY
            ` + sortBy + ` DESC
    `
	rows, err := db.Query(sqlStatement, timezoneOffset, timezoneOffset, classSessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		q := Question{}
		err := rows.Scan(&q.ID, &q.SessionID, &q.UserID, &q.Text, &q.Votes, &q.Answered, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

// HasVoted checks if a user has voted for a question.
// It returns true if the user has voted and false if not.
func (s *Question) HasVoted(questionID int, userID int) (bool, error) {
	db := NewDB()

	var voteCount int
	err := db.QueryRow(`
	SELECT
		COUNT(*)
	FROM
		vote
	WHERE
		question_id = $1 AND
		user_id = $2`,
		questionID, userID).Scan(&voteCount)
	if err != nil {
		return false, err
	}

	if voteCount > 0 {
		return true, nil
	}

	return false, nil
}

// HasVotedAll checks if a user has voted for all questions in a session.
// It returns a slice of structs the contain the question ID and a boolean.
func (s *Question) HasVotedAll(sessionID int, userID int) ([]Question, error) {
	db := NewDB()

	// Get all questions for this session
	sqlStatement := `
			SELECT
				id
			FROM
				question
			WHERE
				session_id = $1
			`
	rows, err := db.Query(sqlStatement, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		q := Question{}
		err := rows.Scan(&q.ID)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Check if the user has voted for each question
	for i, q := range questions {
		hasVoted, err := s.HasVoted(q.ID, userID)
		if err != nil {
			return nil, err
		}
		questions[i].UserHasVoted = hasVoted
	}

	return questions, nil
}

// GetVoteCount returns the number of votes for a question.
// It returns any error encountered.
func (s *Question) GetVoteCount(questionID int) (int, error) {
	db := NewDB()
	var voteCount int
	err := db.QueryRow(`
	SELECT
		votes
	FROM
		question
	WHERE
		id = $1`,
		questionID).Scan(&voteCount)
	if err != nil {
		return 0, err
	}

	return voteCount, nil
}

// GetByID returns a question based on the question id.
// It returns any error encountered and a question struct.
func (s *Question) GetByID(questionID, timezone int) (Question, error) {
	db := NewDB()
	defer db.Close()

	var q Question

	stmt, err := db.Prepare(`
		SELECT
			id,
			session_id,
			user_id,
			text,
			votes,
			answered,
			strftime('%H:%M', datetime(created_at, (? || ' minutes'))) as formatted_created_at,
			datetime(updated_at, (? || ' minutes'))
		FROM
			question
		WHERE
			id = ?
	`)
	if err != nil {
		return Question{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(timezone, timezone, questionID).Scan(&q.ID, &q.SessionID, &q.UserID, &q.Text, &q.Votes, &q.Answered, &q.CreatedAt, &q.UpdatedAt)
	if err != nil {
		return Question{}, err
	}

	return q, nil
}

// GetUnansweredQuestions returns all unanswered questions for a given class session by votes or time.
// It returns any error encountered and a slice of Question structs.
func (s *Question) GetUnansweredQuestions(classSessionID int, timezoneOffset int, sortBy string) ([]Question, error) {

	db := NewDB()

	// Get all questions for this session
	sqlStatement := `
		SELECT
			id, 
			session_id, 
			user_id, 
			text, votes, 
			answered,
			strftime('%H:%M', datetime(created_at, (? || ' minutes'))) as formatted_created_at,
			datetime(updated_at, (? || ' minutes'))
		FROM
			question
		WHERE
			session_id = ? AND
			answered = false
	`

	if sortBy == "votes" {
		sqlStatement += "ORDER BY votes DESC"
	} else if sortBy == "time" {
		sqlStatement += "ORDER BY created_at DESC"
	} else {
		// Default sorting, or handle the invalid sortBy parameter
		sqlStatement += "ORDER BY votes DESC"
	}

	rows, err := db.Query(sqlStatement, timezoneOffset, timezoneOffset, classSessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		q := Question{}
		err := rows.Scan(&q.ID, &q.SessionID, &q.UserID, &q.Text, &q.Votes, &q.Answered, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

// ** UPDATE **
// MarkQuestion marks a question as answered.
// It returns any error encountered.
func (s *Question) MarkQuestion(questionID int) error {
	db := NewDB()
	sqlStatement := `
	UPDATE
		question
	SET
		answered = true
	WHERE
		id = $1`
	_, err := db.Exec(sqlStatement, questionID)
	if err != nil {
		return err
	}

	return nil
}

// VoteQuestion adds a vote to a question based on the question ID.
// It returns any error encountered.
func (s *Question) VoteQuestion(questionID int, userID int) error {
	db := NewDB()

	// Insert vote into the database
	sqlStatement := `
	INSERT INTO
		vote
	VALUES
		(NULL,
		$1,
		$2
		)`
	_, err := db.Exec(sqlStatement, questionID, userID)
	if err != nil {
		return err
	}

	// Get the current number of votes for the question
	var currentVotes int
	err = db.QueryRow(`
	SELECT
		votes
	FROM
		question
	WHERE
		id = $1`,
		questionID).Scan(&currentVotes)
	if err != nil {
		return err
	}

	// Increment the number of votes for the question
	sqlStatement = `
	UPDATE
		question
	SET
		votes = $1
	WHERE
		id = $2`
	_, err = db.Exec(sqlStatement, currentVotes+1, questionID)
	if err != nil {
		return err
	}

	return nil
}

// ** HELPER **
// HasVotedAppend takes a slice of question structs and a slice of voted question structs.
// It returns a slice of question structs with a new voted field set to true if the question has been voted on by that user and false if not.
func (s *Question) HasVotedAppend(questions []Question, votedQuestions []Question) []Question {

	for i, q := range questions {
		for _, v := range votedQuestions {
			if q.ID == v.ID && v.UserHasVoted == true {
				questions[i].UserHasVoted = true
			}
		}
	}
	return questions
}
