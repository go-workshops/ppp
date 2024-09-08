package main

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"go.uber.org/zap"
)

type req struct {
	Author string
	Book   string
}

func tx1(logger *zap.Logger, session db.Session, wg *sync.WaitGroup, limiter chan struct{}, r req) {
	defer wg.Done()
	defer func() { <-limiter }()
	_ = session.Tx(func(sess db.Session) error {
		authorID := uuid.New().String()
		a := author{ID: authorID, Name: r.Author}
		if _, err := sess.SQL().InsertInto("authors1").Values(a).Exec(); err != nil {
			return err
		}
		logger.Info("inserted author", zap.String("author_id", authorID))
		bookID := uuid.New().String()
		b := book{ID: bookID, Title: r.Book, AuthorID: authorID}
		if _, err := sess.SQL().InsertInto("books1").Values(b).Exec(); err != nil {
			return err
		}
		bookCopyID := uuid.New().String()
		logger.Info("inserted book", zap.String("book_id", bookID))
		c := bookCopy{ID: bookCopyID, Status: "AVAILABLE", BookID: bookID}
		if _, err := sess.SQL().InsertInto("book_copies1").Values(c).Exec(); err != nil {
			return err
		}
		logger.Info("inserted book", zap.String("book_copy_id", bookCopyID))
		return nil
	})
}

func tx2(logger *zap.Logger, session db.Session, wg *sync.WaitGroup, limiter chan struct{}, r req) {
	defer wg.Done()
	defer func() { <-limiter }()
	_ = session.Tx(func(sess db.Session) error {
		authorID := uuid.New().String()
		a := author{ID: authorID, Name: r.Author}
		logger.Info("inserting author", zap.Any("author", a))
		if _, err := sess.SQL().InsertInto("authors2").Values(a).Exec(); err != nil {
			return err
		}
		logger.Info("inserted author", zap.Any("author", a))
		bookID := uuid.New().String()
		b := book{ID: bookID, Title: r.Book, AuthorID: authorID}
		logger.Info("inserting book", zap.Any("book", b))
		if _, err := sess.SQL().InsertInto("books2").Values(b).Exec(); err != nil {
			return err
		}
		logger.Info("inserted book", zap.Any("book", b))
		bookCopyID := uuid.New().String()
		c := bookCopy{ID: bookCopyID, Status: "AVAILABLE", BookID: bookID}
		logger.Info("inserting book copy", zap.Any("book_copy", c))
		if _, err := sess.SQL().InsertInto("book_copies2").Values(c).Exec(); err != nil {
			return err
		}
		logger.Info("inserted book copy", zap.Any("book_copy", c))
		return nil
	})
}

type author struct {
	ID   string `db:"id,omitempty"`
	Name string `db:"name"`
}

type book struct {
	ID       string `db:"id,omitempty"`
	AuthorID string `db:"author_id"`
	Title    string `db:"title"`
}

type bookCopy struct {
	ID     string `db:"id,omitempty"`
	BookID string `db:"book_id"`
	Status string `db:"status"`
}

func postgres() db.Session {
	connURL, err := postgresql.ParseURL("postgres://user:password@localhost:5432/db")
	if err != nil {
		log.Fatalf("could not parse db connection url: %v", err)
	}
	sess, err := postgresql.Open(connURL)
	if err != nil {
		log.Fatalf("could not open db connection: %v", err)
	}
	return sess
}
