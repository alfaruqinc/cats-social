package handler

import (
	"cats-social/internal/domain"
	"cats-social/internal/repository"
	"database/sql"
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

func HandleAddNewCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		catBody := domain.NewCat()
		if err := c.ShouldBindJSON(&catBody); err != nil {
			c.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		err := validateRequestBody(*catBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		userReq, _ := c.Get("userData")
		user := userReq.(*domain.User)

		catBody.OwnedById = user.Id

		err = repository.NewCatRepository().CreateCat(db, catBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "something went wrong")
			panic(err)
		}

		res := &domain.CreateCatResponse{
			ID:        catBody.ID,
			CreatedAt: catBody.CreatedAt,
		}

		c.JSON(http.StatusCreated, gin.H{"message": "success", "data": res})
	}
}

func HandleGetAllCats(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
		SELECT id, name, race, sex,
			age_in_month, image_urls, description,
			created_at, has_matched
		FROM cats
		WHERE deleted = false
		`

		userReq, _ := c.Get("userData")
		user := userReq.(*domain.User)

		queryParams := c.Request.URL.Query()
		whereClause, limitOffsetClause, args := validateGetAllCatsQueryParams(queryParams, user.Id.String())

		if len(whereClause) > 0 {
			query += "AND " + strings.Join(whereClause, " AND ")
		}
		query += "\n" + "ORDER BY created_at DESC"
		query += "\n" + strings.Join(limitOffsetClause, " ")

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "something went wrong")
			panic(err)
		}
		defer rows.Close()

		cats := []*domain.Cat{}
		m := pgtype.NewMap()

		for rows.Next() {
			cat := &domain.Cat{}

			err = rows.Scan(&cat.ID, &cat.Name, &cat.Race, &cat.Sex, &cat.AgeInMonth, m.SQLScanner(&cat.ImageUrls), &cat.Description, &cat.CreatedAt, &cat.HasMatched)
			if err != nil {
				c.JSON(http.StatusInternalServerError, "something went wrong")
				panic(err)
			}

			cats = append(cats, cat)
		}

		c.JSON(http.StatusOK, gin.H{"message": "success", "data": &cats})
	}
}

func HandleUpdateCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		catId := c.Param("catId")
		parsedCatId, err := uuid.Parse(catId)
		if err != nil {
			c.JSON(http.StatusNotFound, domain.NewNotFoundError("cat is not found"))
			return
		}

		userReq, _ := c.Get("userData")
		user := userReq.(*domain.User)

		catBody := domain.NewCat()
		if err := c.ShouldBindJSON(catBody); err != nil {
			c.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		err = validateRequestBody(*catBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		catBody.ID = parsedCatId
		catRepo := repository.NewCatRepository()

		err = catRepo.CheckCatIdExists(db, catBody.ID, user.Id)
		if err != nil {
			c.JSON(http.StatusNotFound, domain.NewNotFoundError(err.Error()))
			return
		}

		err = catRepo.CheckEditableSex(db, catBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, domain.NewBadRequest(err.Error()))
			return
		}

		err = catRepo.UpdateCat(db, catBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "something went wrong")
			panic(err)
		}

		updatedAt := time.Now().Format(time.RFC3339)
		parsedUpdatedAt, _ := time.Parse(time.RFC3339, updatedAt)

		updatedCat := domain.UpdateCatResponse{
			ID:        parsedCatId,
			UpdatedAt: parsedUpdatedAt,
		}

		c.JSON(http.StatusOK, domain.NewStatusOk("success", updatedCat))
	}
}

func HandleDeleteCat(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		catId := c.Param("catId")
		parsedCatId, err := uuid.Parse(catId)
		if err != nil {
			c.JSON(http.StatusNotFound, domain.NewNotFoundError("cat is not found"))
			return
		}

		userReq, _ := c.Get("userData")
		user := userReq.(*domain.User)

		catRepo := repository.NewCatRepository()

		err = catRepo.CheckCatIdExists(db, parsedCatId, user.Id)
		if err != nil {
			c.JSON(http.StatusNotFound, domain.NewNotFoundError(err.Error()))
			return
		}

		err = catRepo.DeleteCat(db, parsedCatId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "something went wrong")
			panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"message": "success delete cat"})
	}
}

func validateGetAllCatsQueryParams(queryParams url.Values, userId string) ([]string, []string, []any) {
	var limitOffsetClause []string
	var whereClause []string
	var args []any

	for key, value := range queryParams {
		undefinedParam := slices.Contains(domain.CatQueryParams, key) != true
		limitOffset := key == "limit" || key == "offset"
		emptyValue := len(value[0]) < 1

		if limitOffset {
			limitOffsetClause = append(limitOffsetClause, fmt.Sprintf("%s $%d", key, len(args)+1))

			if key == "limit" && emptyValue {
				value[0] = "5"
			}
			if key == "offset" && emptyValue {
				value[0] = "0"
			}

			args = append(args, value[0])
			continue
		}

		qParamsToSkip := undefinedParam || limitOffset || emptyValue
		if qParamsToSkip {
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

			// regex to extract operator (>,=,<) and number
			extractOperatorAndNumber := regexp.MustCompile(`([>=<])(\d+)`)
			matches := extractOperatorAndNumber.FindStringSubmatch(value[0])
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

			key = "owned_by_id"

			if value[0] == "false" {
				whereClause = append(whereClause, fmt.Sprintf("%s != $%d", key, len(args)+1))
				args = append(args, userId)
				continue
			}

			value[0] = userId
		}

		if key == "search" {
			key = "name"
		}

		whereClause = append(whereClause, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value[0])
	}

	return whereClause, limitOffsetClause, args
}
