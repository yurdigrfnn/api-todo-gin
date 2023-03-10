# Todo Rest Api With JWT auth

## REST API docs

- this is API docs
- Inspired by Swagger API docs style & structure: https://petstore.swagger.io/#/pet

------------------------------------------------------------------------------------------

#### Creating new user & Login

<details>
 <summary><code>POST</code> <code><b>/</b></code> <code>(This method for register or login)</code></summary>

##### Parameters

> | name               |  type       | data type               | body                                                     |description                                       |
> |--------------------|-------------|-------------------------|----------------------------------------------------------|------------------------------------------------|
> | `api/register`     |  `required` | `object (JSON)`         | `{"email" : "example@email" ,"password" : "example"}`    |`create user`                                       |
> | `api/signin`       |  `required` | `object (JSON)`         | `{"email" : "example@email" ,"password" : "example"}`    |`Login user `                                      |


##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`                | `{"isError": false, "message": "succes to create user/login"}`      |
> | `401`         | `application/json`                | `{"isError": true, "message": "unautorization"}`                    |

</details>

------------------------------------------------------------------------------------------

#### Get Todo

<details>
 <summary><code>GET</code> <code><b>/</b></code> <code>(this method for get todo and get validation)</code></summary>

##### Parameters

> | name               |  type       | data type               | response                                                 |description                                   |
> |--------------------|-------------|-------------------------|----------------------------------------------------------|----------------------------------------------|
> | `api/todo`         |  `required` | `object (JSON)`         | `{"isError": false,"todos": [],"page": 2,`               |`get todo with parameter limit & page`        |           > |                    |             |                         |`"limit": 10,"totalPages": 1,"totalTodos": 3}`            |                                              |
> | `api/validate`     |  `required` | `object (JSON)`         | `{"email" : "example@email" ,"password" : "example"}`    |`get user `                                   |

</details>


<details>
  <summary><code>PUT</code> <code>(edit todo by id)</code></summary>

##### Parameters

> | name            |  type      | body                                               | description                                          |
> |-----------------|------------|----------------------------------------------------|------------------------------------------------------|
> | `api/todo/{id}` |  required  | `{"name" : "example" ,"complete" : bool }`         | Edit tode                                            |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`                | `{"isError": false,"todos": {} }`                                   |
> | `401`         | `application/json`                | `{"isError": true ,"message":"unauthoriz/todo not found"}`          |

</details>


<details>
  <summary><code>DELETE</code> <code>(DELETE TODO)</code></summary>
    
##### Parameters

> | name            |  type      | description                                        |
> |-----------------|------------|----------------------------------------------------|------------------------------------------------------|
> | `api/todo/{id}` |  required  | DELETE tode                                            |

##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`                | `{"isError": false,"message": "Todo deleted"}`                      |
> | `401`         | `application/json`                | `{"isError": true ,"message":"unauthoriz/todo not found"}`          |

</details>
