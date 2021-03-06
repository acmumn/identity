package main

import (
	"fmt"

	"github.com/acmumn/identity/server/db"
	"github.com/acmumn/identity/server/token"
	"github.com/gin-gonic/gin"
)

// GetLogin is the handler for /login/:uuid with the method GET.
func GetLogin(db *db.DB, toks *token.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Param("uuid")
		panic(fmt.Sprintf("TODO %#v", uuid))
	}
}
