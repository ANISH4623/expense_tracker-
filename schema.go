package main

import (
	"awesomeProject1/database"
	"awesomeProject1/models"
	_ "context"
	"errors"

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
			Type: graphql.DateTime,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
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
			Type: graphql.DateTime,
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
		},
		"deletedAt": &graphql.Field{
			Type: graphql.DateTime,
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
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := uint(params.Context.Value("user_id").(int64))
				var expenses []models.Expense
				database.Database.Db.Where("user_id = ?", userId).Find(&expenses)
				return expenses, nil
			},
		},
		"incomes": &graphql.Field{
			Type: graphql.NewList(incomeType),
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := uint(params.Context.Value("user_id").(int64))
				var incomes []models.Income
				database.Database.Db.Where("user_id = ?", userId).Find(&incomes)
				return incomes, nil
			},
		},
	},
})

// Define the root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"createExpense": &graphql.Field{
			Type: expenseType,
			Args: graphql.FieldConfigArgument{
				"amount": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
				"category": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := uint(params.Context.Value("user_id").(int64))

				expense := models.Expense{
					UserID:   userId,
					Amount:   params.Args["amount"].(float64),
					Category: params.Args["category"].(string),
				}
				database.Database.Db.Create(&expense)
				return expense, nil
			},
		},
		"updateExpense": &graphql.Field{
			Type: expenseType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"amount": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"category": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := uint(params.Context.Value("user_id").(int64))
				var expense models.Expense
				database.Database.Db.First(&expense, params.Args["id"].(uint))
				if expense.UserID != userId {
					return nil, errors.New("unauthorized")
				}
				if amount, ok := params.Args["amount"].(float64); ok {
					expense.Amount = amount
				}
				if category, ok := params.Args["category"].(string); ok {
					expense.Category = category
				}
				database.Database.Db.Save(&expense)
				return expense, nil
			},
		},
		"createIncome": &graphql.Field{
			Type: incomeType,
			Args: graphql.FieldConfigArgument{
				"amount": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
				"category": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := uint(params.Context.Value("user_id").(int64))
				income := models.Income{
					UserID:   userId,
					Amount:   params.Args["amount"].(float64),
					Category: params.Args["category"].(string),
				}
				database.Database.Db.Create(&income)
				return income, nil
			},
		},
		"updateIncome": &graphql.Field{
			Type: incomeType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"amount": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"category": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				userId := uint(params.Context.Value("user_id").(int64))
				var income models.Income
				database.Database.Db.First(&income, params.Args["id"].(uint))
				if income.UserID != userId {
					return nil, errors.New("unauthorized")
				}
				if amount, ok := params.Args["amount"].(float64); ok {
					income.Amount = amount
				}
				if category, ok := params.Args["category"].(string); ok {
					income.Category = category
				}
				database.Database.Db.Save(&income)
				return income, nil
			},
		},
	},
})

// Define the schema
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
