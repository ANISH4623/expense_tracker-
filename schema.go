package main

import (
	"awesomeProject1/database"
	"awesomeProject1/models"

	"github.com/graphql-go/graphql"
)

// Define the User type
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"expenses": &graphql.Field{
			Type: graphql.NewList(expenseType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user := p.Source.(models.User)
				var expenses []models.Expense
				database.Database.Db.Where("user_id = ?", user.ID).Find(&expenses)
				return expenses, nil
			},
		},
		"incomes": &graphql.Field{
			Type: graphql.NewList(incomeType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user := p.Source.(models.User)
				var incomes []models.Income
				database.Database.Db.Where("user_id = ?", user.ID).Find(&incomes)
				return incomes, nil
			},
		},
	},
})

// Define the Expense type
var expenseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Expense",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
		"deletedAt": &graphql.Field{
			Type: graphql.String,
		},
		"userId": &graphql.Field{
			Type: graphql.String,
		},
		"amount": &graphql.Field{
			Type: graphql.Float,
		},
		"category": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// Define the Income type
var incomeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Income",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.String,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.String,
		},
		"deletedAt": &graphql.Field{
			Type: graphql.String,
		},
		"userId": &graphql.Field{
			Type: graphql.String,
		},
		"amount": &graphql.Field{
			Type: graphql.Float,
		},
		"category": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// Define the root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"users": &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				var users []models.User
				database.Database.Db.Find(&users)
				return users, nil
			},
		},
		"user": &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(string)
				var user models.User
				database.Database.Db.First(&user, id)
				return user, nil
			},
		},
		"expenses": &graphql.Field{
			Type: graphql.NewList(expenseType),
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId, _ := params.Args["userId"].(string)
				var expenses []models.Expense
				database.Database.Db.Where("user_id = ?", userId).Find(&expenses)
				return expenses, nil
			},
		},
		"incomes": &graphql.Field{
			Type: graphql.NewList(incomeType),
			Args: graphql.FieldConfigArgument{
				"userId": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId, _ := params.Args["userId"].(string)
				var incomes []models.Income
				database.Database.Db.Where("user_id = ?", userId).Find(&incomes)
				return incomes, nil
			},
		},
	},
})

// Define the schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})
