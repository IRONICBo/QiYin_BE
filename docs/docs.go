// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/check": {
            "get": {
                "description": "Check whether the token is valid",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "CheckToken",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/comment/add": {
            "post": {
                "description": "add comment",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comment"
                ],
                "summary": "CommentAdd",
                "parameters": [
                    {
                        "description": "CommentAddParams",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requestparams.CommentAddParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/comment/delete": {
            "post": {
                "description": "Test API",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comment"
                ],
                "summary": "delete comment",
                "parameters": [
                    {
                        "description": "CommentDelParams",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requestparams.CommentDelParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/comment/list": {
            "get": {
                "description": "get comment list by videoId",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "comment"
                ],
                "summary": "CommentList",
                "parameters": [
                    {
                        "type": "string",
                        "description": "query video id",
                        "name": "videoId",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "videoId",
                        "name": "videoId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/favorite/action": {
            "post": {
                "description": "like or dislike",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "favorite"
                ],
                "summary": "FavoriteAction",
                "parameters": [
                    {
                        "description": "FavoriteParams",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requestparams.FavoriteParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/favorite/list": {
            "get": {
                "description": "get favorite video list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "favorite"
                ],
                "summary": "GetFavoriteList",
                "parameters": [
                    {
                        "type": "string",
                        "description": "query user id",
                        "name": "userId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/login": {
            "post": {
                "description": "user login",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "UserLogin",
                "parameters": [
                    {
                        "description": "UserParams",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requestparams.UserParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/ping": {
            "post": {
                "description": "Test API",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Ping"
                ],
                "summary": "Ping",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/qiniu/pfop/callback": {
            "get": {
                "description": "Get QiNiu Pfop callback result",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "QiNiu"
                ],
                "summary": "GetPfopCallback",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/qiniu/proxy": {
            "get": {
                "description": "Get QiNiu image by proxy",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "QiNiu"
                ],
                "summary": "GetImageByProxy",
                "parameters": [
                    {
                        "type": "string",
                        "description": "url",
                        "name": "url",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/qiniu/token": {
            "post": {
                "description": "Get QiNiu upload token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "QiNiu"
                ],
                "summary": "UserLogin",
                "parameters": [
                    {
                        "description": "QiNiuTokenParams",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requestparams.QiNiuTokenParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/register": {
            "post": {
                "description": "user register",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "UserRegister",
                "parameters": [
                    {
                        "description": "UserParams",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requestparams.UserParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/searchUser": {
            "get": {
                "description": "search user by name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "SearchUser",
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/setStyle": {
            "post": {
                "description": "set user style",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "SetStyle",
                "parameters": [
                    {
                        "description": "StyleParams",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requestparams.StyleParams"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/userinfo": {
            "get": {
                "description": "get userinfo by id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "UserInfo",
                "parameters": [
                    {
                        "type": "string",
                        "description": "query user id",
                        "name": "userId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/video/hots": {
            "get": {
                "description": "hot list",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "video"
                ],
                "summary": "GetHots",
                "parameters": [
                    {
                        "type": "string",
                        "description": "searchValue",
                        "name": "searchValue",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        },
        "/api/v1/video/search": {
            "get": {
                "description": "search videos by text",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "video"
                ],
                "summary": "Search",
                "parameters": [
                    {
                        "type": "string",
                        "description": "searchValue",
                        "name": "searchValue",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/response.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "msg": {
                                            "type": "string"
                                        }
                                    }
                                }
                            ]
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "requestparams.CollectionParams": {
            "type": "object",
            "properties": {
                "actionType": {
                    "description": "1 点赞，-1 取消",
                    "type": "integer"
                },
                "videoId": {
                    "description": "UserID     string ` + "`" + `json:userId` + "`" + `",
                    "type": "integer"
                }
            }
        },
        "requestparams.CommentAddParams": {
            "type": "object",
            "properties": {
                "commentText": {
                    "type": "string"
                },
                "videoId": {
                    "type": "integer"
                }
            }
        },
        "requestparams.CommentDelParams": {
            "type": "object",
            "properties": {
                "commentId": {
                    "type": "integer"
                }
            }
        },
        "requestparams.FavoriteParams": {
            "type": "object",
            "properties": {
                "actionType": {
                    "description": "1 点赞，-1 取消",
                    "type": "integer"
                },
                "videoId": {
                    "description": "UserID     string ` + "`" + `json:userId` + "`" + `",
                    "type": "integer"
                }
            }
        },
        "requestparams.QiNiuTokenParams": {
            "type": "object",
            "required": [
                "ticket"
            ],
            "properties": {
                "ticket": {
                    "type": "string"
                }
            }
        },
        "requestparams.StyleParams": {
            "type": "object",
            "properties": {
                "style": {
                    "type": "string"
                }
            }
        },
        "requestparams.UserParams": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "response.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization.",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v0.0.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "QiYin Backend",
	Description:      "QiYin Backend API Docs.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
