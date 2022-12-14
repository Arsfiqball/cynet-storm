package integration

import (
	"app/pkg/book"
	"app/pkg/restpl"
	"app/pkg/restql"
	"app/test/testhelper"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BookTestSuite struct {
	suite.Suite
	testhelper.AppSuite
}

func (s *BookTestSuite) SetupTest() {
	s.Start(s.T())
	s.ExecSQL(s.T(), "DELETE FROM books")
}

func (s *BookTestSuite) TearDownTest() {
	s.Stop(s.T())
}

func (s *BookTestSuite) TestCreateOne() {
	type scenarioTemplate struct {
		name        string
		payload     book.PayloadBook
		checkEntity book.EntityBook
	}

	sampleTime, err := time.Parse(time.RFC3339, "2023-01-01T23:11:57+07:00")
	if err != nil {
		s.Fail(err.Error())
	}

	scenarios := []scenarioTemplate{
		{
			name: "Insert a book",
			payload: book.PayloadBook{
				Mutation: restpl.Mutation{
					Fields:     []string{"title", "author", "series", "volume", "fileUrl", "coverUrl", "publishDate"},
					NullFields: []string{},
				},

				Title:       "Clean Architecture",
				Author:      "Uncle Bob",
				Series:      "Programming",
				Volume:      3,
				FileUrl:     "https://book.com/clean.pdf",
				CoverUrl:    "https://book.com/cover_clean.jpg",
				PublishDate: sampleTime,
			},
			checkEntity: book.EntityBook{
				Title:       "Clean Architecture",
				Author:      "Uncle Bob",
				Series:      "Programming",
				Volume:      3,
				FileUrl:     sql.NullString{String: "https://book.com/clean.pdf", Valid: true},
				CoverUrl:    sql.NullString{String: "https://book.com/cover_clean.jpg", Valid: true},
				PublishDate: sql.NullTime{Time: sampleTime, Valid: true},
			},
		},
		{
			name: "Insert a book only required field",
			payload: book.PayloadBook{
				Mutation: restpl.Mutation{
					Fields:     []string{"title", "author", "series", "volume"},
					NullFields: []string{},
				},

				Title:  "Clean Architecture",
				Author: "Uncle Bob",
				Series: "Programming",
				Volume: 3,
			},
			checkEntity: book.EntityBook{
				Title:       "Clean Architecture",
				Author:      "Uncle Bob",
				Series:      "Programming",
				Volume:      3,
				FileUrl:     sql.NullString{},
				CoverUrl:    sql.NullString{},
				PublishDate: sql.NullTime{},
			},
		},
	}

	for _, scenario := range scenarios {
		s.Run(scenario.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := s.AppSuite.
				GetApp().
				BookSet.
				BookRepository.
				CreateOne(ctx, scenario.payload)

			if err != nil {
				s.Fail(err.Error())
			}

			// DEBUG by adding verbose flag (-v)
			// testhelper.JsonPrint(result)

			assert.Equal(s.T(), scenario.checkEntity.Title, result.Title)
			assert.Equal(s.T(), scenario.checkEntity.Author, result.Author)
			assert.Equal(s.T(), scenario.checkEntity.Series, result.Series)
			assert.Equal(s.T(), scenario.checkEntity.Volume, result.Volume)
			assert.Equal(s.T(), scenario.checkEntity.FileUrl, result.FileUrl)
			assert.Equal(s.T(), scenario.checkEntity.CoverUrl, result.CoverUrl)
			// assert.True(s.T(), scenario.checkEntity.PublishDate.Time.Equal(result.PublishDate.Time), "PublishDate is not equal")
		})
	}
}

func SimpleMultiInt(key string, op string, vals int) restql.MultiInt {
	cond := map[string][]interface{}{}
	cond[op] = []interface{}{vals}

	return restql.MultiInt{
		Key:       key,
		Condition: cond,
	}
}

func SimpleMultiString(key string, op string, vals string) restql.MultiString {
	cond := map[string][]interface{}{}
	cond[op] = []interface{}{vals}

	return restql.MultiString{
		Key:       key,
		Condition: cond,
	}
}

func (s *BookTestSuite) TestGetOne() {
	s.ExecSQLTestFile(s.T(), "book/book_samples.sql")

	type scenarioTemplate struct {
		name        string
		query       book.QueryBook
		checkEntity book.EntityBook
	}

	scenarios := []scenarioTemplate{
		{
			name: "Get a book by id",
			query: book.QueryBook{
				ID: SimpleMultiInt("ID", "eq", 3),
			},
			checkEntity: book.EntityBook{},
		},
		{
			name: "Get a book by author name contain sub-string (case sensitive)",
			query: book.QueryBook{
				Author: SimpleMultiString("Author", "contains", "Shield"),
			},
			checkEntity: book.EntityBook{},
		},
		{
			name: "Get a book from a series (case insensitive) with specific volume",
			query: book.QueryBook{
				Series: SimpleMultiString("Series", "contain", "kane"),
				Volume: SimpleMultiInt("Volume", "eq", 3),
			},
			checkEntity: book.EntityBook{},
		},
	}

	for _, scenario := range scenarios {
		s.Run(scenario.name, func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			result, err := s.AppSuite.
				GetApp().
				BookSet.
				BookRepository.
				GetOne(ctx, scenario.query)

			if err != nil {
				s.Fail(err.Error())
			}

			// DEBUG by adding verbose flag (-v)
			testhelper.JsonPrint(result)

			// assert.Equal(s.T(), scenario.checkEntity.Title, result.Title)
			// assert.Equal(s.T(), scenario.checkEntity.Author, result.Author)
			// assert.Equal(s.T(), scenario.checkEntity.Series, result.Series)
			// assert.Equal(s.T(), scenario.checkEntity.Volume, result.Volume)
			// assert.Equal(s.T(), scenario.checkEntity.FileUrl, result.FileUrl)
			// assert.Equal(s.T(), scenario.checkEntity.CoverUrl, result.CoverUrl)
			// assert.True(s.T(), scenario.checkEntity.PublishDate.Time.Equal(result.PublishDate.Time), "PublishDate is not equal")
		})
	}
}

func TestBookTestSuite(t *testing.T) {
	suite.Run(t, new(BookTestSuite))
}
