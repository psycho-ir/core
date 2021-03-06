package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"upper.io/db.v3/lib/sqlbuilder"
)

/**
 * @api {post} /v1/companies Insert New Company
 * @apiName AddCompany
 * @apiGroup Companies
 * @apiVersion 0.1.0
 *
 * @apiUse jwt
 *
 * @apiParam (Request body) {String} name Name Of Company.
 */
func insertNewCompany(c *gin.Context) {
	var company Companies
	loginUser := c.MustGet("user").(LoginUser)

	if err := c.ShouldBindWith(&company, binding.JSON); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	ctx := context.Background()

	err := MySql.Tx(ctx, func(tx sqlbuilder.Tx) error {
		newCompany, err := tx.InsertInto("companies").Columns("name").Values(company.Name).Exec()
		if err != nil {
			return err
		}

		newCompanyId, _ := newCompany.LastInsertId()

		_, err = tx.InsertInto("employees").Columns(
			"company_id",
			"user_id",
			"type",
			"status_by_employee",
			"status_by_company",
		).Values(
			newCompanyId,
			loginUser.Id,
			"manager",
			"active",
			"active",
		).Exec()

		if err != nil {
			return err
		}

		// If the function returns no error the transaction is committed.
		return nil
	})

	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "successful"})
	return
}

/**
 * @api {get} /v1/companies Get All Companies
 * @apiName GetAllCompanies
 * @apiGroup Companies
 * @apiVersion 0.1.0
 *
 * @apiUse jwt
 */
func getAllCompanies(c *gin.Context) {
	var companies []Companies
	loginUser := c.MustGet("user").(LoginUser)

	err := MySql.Select(
		"companies.id as id",
		"companies.name as name",
		"companies.status as status",
		"companies.created_at as created_at",
		"companies.deleted_at as deleted_at",
	).From("companies").
		Where("employees.user_id LIKE ?", loginUser.Id).
		Join("employees").On("employees.company_id = companies.id").
		All(&companies)
	if err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": companies})
	return
}

type Companies struct {
	Id        uint   `db:"id" json:"id"`
	Name      string `db:"name" json:"name" binding:"required,gte=2,lte=64"`
	Status    string `db:"status" json:"status"`
	CreatedAt string `db:"created_at" json:"created_at"`
	DeletedAt string `db:"deleted_at" json:"deleted_at"`
}
