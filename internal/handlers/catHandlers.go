package handlers

import (
	"cats-social/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type response struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func HandleAddNewCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		catBody := models.NewCat()
		if err := c.ShouldBindJSON(&catBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(catBody.Name) < 1 || len(catBody.Name) > 30 {
			err := errors.New("name length should be between 1 and 30 characters")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if slices.Contains(models.CatRace, catBody.Race) != true {
			err := errors.New("accepted race is only Persian, Maine Coon, Siamese, Ragdoll, Bengal, Sphynx, British Shorthair, Abyssinian, Scottish Fold, Birman")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if slices.Contains(models.CatSex, catBody.Sex) != true {
			err := errors.New("accepted sex is only male and female")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if catBody.AgeInMonth < 1 || catBody.AgeInMonth > 120082 {
			err := errors.New("your cat's age is minimum 1 month and maximum 120082 month")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(catBody.Description) < 1 || len(catBody.Description) > 200 {
			err := errors.New("description length should be between 1 and 200 characters")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(catBody.ImageUrls) < 1 {
			err := errors.New("image urls at least have 1 image")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, imageUrl := range catBody.ImageUrls {
			if len(imageUrl) < 1 {
				err := errors.New("image urls cannot have empty item")
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			_, err := url.ParseRequestURI(imageUrl)
			if err != nil {
				err := errors.New("image url should have valid url")
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}

		// TODO: delete after auth api finish
		parsed, _ := uuid.Parse("e91ce26e-9a53-4c4f-b5b5-0cad1a61d82b")
		catBody.OwnedBy = parsed

		query := `INSERT INTO cats (id, created_at, name, race, sex, age_in_month, description, image_urls, owned_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`
		_, err := db.Exec(query, catBody.ID, catBody.CreatedAt, catBody.Name, catBody.Race, catBody.Sex, catBody.AgeInMonth, catBody.Description, catBody.ImageUrls, catBody.OwnedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		res := &response{
			Id:        catBody.ID,
			CreatedAt: catBody.CreatedAt,
		}

		c.JSON(201, gin.H{"message": "success", "data": res})
	}
}

func HandleGetAllCats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
		SELECT id, name, race, sex,
			age_in_month, image_urls, description,
			created_at, has_matched
		FROM cats
		`

		queryParams := c.Request.URL.Query()
		var args []any
		var whereClause []string
		for key, value := range queryParams {
			undefinedParam := slices.Contains(models.CatQueryParams, key) != true
			limitOffset := key == "limit" || key == "offset"
			emptyValue := len(value[0]) < 1
			if undefinedParam || limitOffset || emptyValue {
				continue
			}

			if key == "id" {
				_, err := uuid.Parse(value[0])
				if err != nil {
					continue
				}
			}

			if key == "hasMatched" {
				key = "has_matched"
			}

			if key == "ageInMonth" {
				key = "age_in_month"

				re := regexp.MustCompile(`([>=<])(\d+)`)
				matches := re.FindStringSubmatch(value[0])
				if len(matches) != 3 {
					continue
				}

				opr := matches[1]
				val := matches[2]

				whereClause = append(whereClause, fmt.Sprintf("%s %s $%d", key, opr, len(args)+1))
				args = append(args, val)

				continue
			}

			if key == "owned" {
				if value[0] != "true" && value[0] != "false" {
					continue
				}

				key = "owned_by"
				// TODO: change value of userId with user id after auth api finish
				userId := "e91ce26e-9a53-4c4f-b5b5-0cad1a61d82b"

				if value[0] == "false" {
					whereClause = append(whereClause, fmt.Sprintf("%s != $%d", key, len(args)+1))
					args = append(args, userId)
					continue
				}

				// TODO: change value of value[0] with user id after auth api finish
				value[0] = "e91ce26e-9a53-4c4f-b5b5-0cad1a61d82b"
			}

			if key == "search" {
				key = "name"
			}

			whereClause = append(whereClause, fmt.Sprintf("%s = $%d", key, len(args)+1))
			args = append(args, value[0])
		}

		if len(whereClause) > 0 {
			query += " WHERE " + strings.Join(whereClause, " AND ")
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		cats := []*models.Cat{}
		m := pgtype.NewMap()

		for rows.Next() {
			cat := &models.Cat{}

			err = rows.Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, m.SQLScanner(&cat.ImageUrls), &cat.Description, &cat.CreatedAt, &cat.HasMatched)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			cats = append(cats, cat)
		}

		c.JSON(http.StatusOK, gin.H{"message": "success", "data": &cats})
	}
}
