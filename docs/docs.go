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
        "/chat": {
            "post": {
                "description": "Sends a message to the AI model with a Wikipedia context.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Chat"
                ],
                "summary": "Chat with AI",
                "parameters": [
                    {
                        "description": "Chat Request Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.Chat"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "$ref": "#/definitions/routes.ChatResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/infobox": {
            "post": {
                "description": "Generates a concise infobox using AI",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wiki"
                ],
                "summary": "Generate article infobox",
                "parameters": [
                    {
                        "description": "Article Request Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.ArticleO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "$ref": "#/definitions/routes.InfoboxRes"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/messages/{chat_id}": {
            "get": {
                "description": "Retrieves all messages for a specific chat ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Messages"
                ],
                "summary": "Get chat messages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Chat ID",
                        "name": "chat_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of messages",
                        "schema": {
                            "$ref": "#/definitions/routes.MessagesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/search": {
            "get": {
                "description": "Searches Wikipedia for articles matching the query",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Search"
                ],
                "summary": "Search Wikipedia",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search query",
                        "name": "q",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Language code (default: en)",
                        "name": "lang",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Search results",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/routes.SearchResult"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/summary": {
            "post": {
                "description": "Generates a concise summary of a Wikipedia article",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Summary"
                ],
                "summary": "Generate article summary",
                "parameters": [
                    {
                        "description": "Article Request Body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/routes.Article"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successful response",
                        "schema": {
                            "$ref": "#/definitions/routes.SummaryResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/wiki/{lang}/{wiki}": {
            "get": {
                "description": "Retrieves a Wikipedia article by language and article name",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wiki"
                ],
                "summary": "Get Wikipedia article",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Language code",
                        "name": "lang",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Wikipedia article name",
                        "name": "wiki",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Wikipedia article content",
                        "schema": {
                            "$ref": "#/definitions/routes.WikiResponse"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Article not found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "routes.Article": {
            "type": "object",
            "properties": {
                "lang": {
                    "type": "string",
                    "example": "en"
                },
                "section": {
                    "type": "string",
                    "example": "Causes"
                },
                "wiki": {
                    "type": "string",
                    "example": "French_Revolution"
                }
            }
        },
        "routes.ArticleO": {
            "type": "object",
            "properties": {
                "lang": {
                    "type": "string",
                    "example": "en"
                },
                "wiki": {
                    "type": "string",
                    "example": "French_Revolution"
                }
            }
        },
        "routes.Chat": {
            "type": "object",
            "properties": {
                "lang": {
                    "type": "string",
                    "example": "en"
                },
                "message": {
                    "type": "string",
                    "example": "What is the capital of France?"
                },
                "model": {
                    "type": "string",
                    "example": "conservative"
                },
                "uuid": {
                    "type": "string",
                    "example": "123e4567-e89b-12d3-a456-426614174000"
                },
                "wiki": {
                    "type": "string",
                    "example": "French_Revolution"
                }
            }
        },
        "routes.ChatResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Paris is the capital of France."
                }
            }
        },
        "routes.InfoboxRes": {
            "type": "object",
            "properties": {
                "infobox": {
                    "type": "string",
                    "example": "{name: 'French', leader: 'Bonaparte'}"
                }
            }
        },
        "routes.MessageResponse": {
            "type": "object",
            "properties": {
                "chat_id": {
                    "type": "string",
                    "example": "123e4567-e89b-12d3-a456-426614174000"
                },
                "content": {
                    "type": "string",
                    "example": "What is the capital of France?"
                },
                "id": {
                    "type": "integer",
                    "example": 2
                },
                "role": {
                    "type": "string",
                    "example": "User"
                },
                "topic": {
                    "type": "string",
                    "example": "French_Revolution"
                }
            }
        },
        "routes.MessagesResponse": {
            "type": "object",
            "properties": {
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/routes.MessageResponse"
                    }
                }
            }
        },
        "routes.SearchResult": {
            "type": "object",
            "properties": {
                "lang": {
                    "type": "string",
                    "example": "en"
                },
                "summary": {
                    "type": "string",
                    "example": "The French Revolution was a period of radical political and societal change in France..."
                },
                "title": {
                    "type": "string",
                    "example": "French Revolution"
                },
                "url": {
                    "type": "string",
                    "example": "/wiki/en/French_Revolution"
                }
            }
        },
        "routes.SummaryResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "The French Revolution was a period of radical political and societal change in France..."
                }
            }
        },
        "routes.WikiResponse": {
            "type": "object",
            "properties": {
                "full_body": {
                    "type": "string",
                    "example": "The French Revolution was a period of radical political and societal change in France..."
                },
                "sections": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "body": {
                                "type": "string",
                                "example": "The French Revolution was a period of radical political and societal change in France..."
                            },
                            "title": {
                                "type": "string",
                                "example": "Causes"
                            }
                        }
                    }
                },
                "summary": {
                    "type": "string",
                    "example": "The French Revolution was a period of radical political and societal change in France..."
                },
                "title": {
                    "type": "string",
                    "example": "French Revolution"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.5",
	Host:             "9.141.41.77:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Wiki? API",
	Description:      "Wiki? API",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
