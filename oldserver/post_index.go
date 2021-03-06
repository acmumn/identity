package main

import (
	"database/sql"
	"net/http"

	"github.com/acmumn/identity/server/db"
	"github.com/acmumn/mailer/go-client"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// PostIndex is the handler for / with the method POST.
func PostIndex(db *db.DB, mailer *mailer.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Redirect string `form:"redirect" json:"redirect" xml:"redirect"`
			X500     string `form:"x500" json:"x500" xml:"x500" binding:"required"`
		}

		err := c.ShouldBind(&body)
		if err != nil {
			c.HTML(http.StatusBadRequest, "get-index", gin.H{
				"Error":    err.Error(),
				"Redirect": body.Redirect,
			})
			return
		}

		id, email, err := db.GetMemberIDAndEmailFromX500(body.X500)
		if err == sql.ErrNoRows {
			c.HTML(http.StatusNotFound, "get-index", gin.H{
				"Error":    "No such member could be found.",
				"Redirect": body.Redirect,
			})
			return
		} else if err != nil {
			log.Error("When querying for member", err)
			c.HTML(http.StatusNotFound, "error", gin.H{
				"Error": err.Error(),
			})
			return
		}

		uuid, err := db.NewLoginUUID(id)
		if err != nil {
			log.Error("When creating login UUID", err)
			c.HTML(http.StatusNotFound, "error", gin.H{
				"Error": err.Error(),
			})
			return
		}

		err = mailer.Send("identity", "login", email, "Log In", map[string]interface{}{
			"uuid": uuid,
		})
		if err != nil {
			log.Error("When sending mail", err)
			c.HTML(http.StatusNotFound, "error", gin.H{
				"Error": err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "post-index", gin.H{"Email": email})
	}
}
