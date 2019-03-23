# gob
###### Note: The name is not final and may change.

gob is go library which routes API calls. Each endpoint is a specific handler typed by GET, POST, PUT, DELETE. gob supports middleware functions to inject data into request contexts. It's an easy method call to generate docs for your entire API.


### gob is based around dynamic handlers

```golang

type UserContext struct {
  DAL *postgres.db // access to your DB provider.
  SomeVal int
}

// user create request model
type UserRequest struct {
  Name  string  `json:"name"`
  Age   int     `json:"age"`
}
// user database model
type UserDatabaseModel struct {
  Name          string      `json:"name" sql:"name"`
  Age           int         `json:"age" sql: "age"`
  IsVerified    bool        `json:"isVerified" sql:"is_verified"`
  HashSecret    string      `sql:"hash_secret"`
  CreatedAt     time.Time   `json:"createdAt" sql:"created_at"`
  UpdatedAt     time.Time   `json:"updatedAt" sql:"updated_at`
}
// user create response model (view model)
type UserResponse struct {
  *UserDatabaseModel
}

// define API handlers as functions that take *input models* and return an *output models*
// with optional response codes & errors
func (c *UserContext) CreateUser(reqModel *UserRequest) (*UserResponse, int, error) {
  // do input sanity checks with `reqModel`
  if reqModel.Age < 18 {
    return nil, 400, errors.New("Cannot register user under the age of 18!")
  }
  
  // more sanity checks here if you want...
  dbModel := helpers.CreateDBModelFromRequest(reqModel)
  if err := c.DAL.CreateUser(&dbModel); err != nil {
    return nil, 400, err
  }
  
  respModel := &UserResponseModel{dbModel}
  
  return respModel, 200, nil
}

func (c *UserContext) RawHandler(w http.ResponseWriter, req *http.Request) {
  // read `req` obj and do stuff
  // ...
  
  // write to the response writer
  io.WriteString(w, "You GET'ed the raw handler!")
}

// middlewares always receive the raw request and response writer objects
func (c *UserContext) MiddlewareSetContext(w http.ResponseWriter, req *http.Request) {
  c.SomeVal = 7
}

// create a new router with your own context
userRouter := gob.NewRouter(UserContext{})

// it supports middlewares too!
userRouter.Middleware((*UserContext).MiddlewareSetContext)

userRouter.Route("POST", "/user", (*UserContext).CreateUser)
userRouter.Route("GET", "/rawendpoint", (*UserContext).RawHandler)
```

#### Dynamic handlers are used by gob for endpoints that take in request models and return response models. The goal is to enable users to easily generate documentation for these endpoints.

>Ex. by running something like `gob.GenerateDocs()`


### TODOs:

- The actual router code underneath is more like proof of concept than actual usable code and needs to be rewritten.

- It has not been tested properly, so test cases have to be implemented.
