# go-repo-structure
Go Repo Structure Based on Uncle Bob' Clean Architecture

### Basic
As we know the constraint before designing the Clean Architecture are :
- Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
- Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
- Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
- Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
- Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

So, based on this constraint, every layer must independent and testable.

Uncle Bob’s Architecture has 4 layers :
- Entities
- Usecase
- Controller
- Framework & Driver

In my projects, I’m using 4 too :
- Models
- Repository
- Manager
- Delivery

### Models
Same as Entities, will used in all layer. This layer, will store any Object’s Struct and its method. Example : Article, Student, Book.

### Repository
Repository will store any Database handler. Querying, or Creating/ Inserting into any database will stored here. This layer will act for CRUD to database only. No business process happen here. Only plain function to Database.

This layer also have responsibility to choose what DB will used in Application. Could be Mysql, MongoDB, MariaDB, Postgresql whatever, will decided here.

If using ORM, this layer will control the input, and give it directly to ORM services.

If calling microservices, will be handled here. Create HTTP Request to other services, and sanitize the data. This layer, must fully act as a repository. Handle all data input - output no specific logic happen.

This Repository layer will depend on Connected DB , or other microservices if exists.

### Manager
This layer will act as the business process handler. Any process will be handled here. This layer will decide, which repository layer will use. And have responsibility to provide data to serve into delivery. Process the data doing calculation or anything will be done here.

Manager layer will accept sanitized input from Delivery layer and then process the input could be storing into DB , or Fetching from DB ,etc.

This layer will depend on Repository(s) Layer

### Delivery
This layer will act as the presenter. Decide how the data will presented. Could be as REST API, or HTML File, or gRPC whatever the delivery type.

This layer also will accept the input from user. Sanitize the input and sent it to Manager layer.

For my sample project, I’m using REST API as the delivery method.

Client will call the resource endpoint over network, and the Delivery layer will get the input or request, and sent it to Manager Layer.

This layer will depend on Manager Layer.

### Communications Between Layers
Except Models, each layer will communicate through intefaces. For example, Manager layer need the Repository layer, so how they communicate? Repository will provide an interface to their contract and communication.

Same with Manager, Delivery layer will use Manager contract interface. 

### Final Output and The Merging

After finishing all layers, You should merge into one system in main.go in root/app project.

Here you will define, and create every needs to environment, and merge all layers into one.

### Final Structure

```
 main.go

 domain
 ├── article.go
 ├── author.go
 ├── codes.go 
 └── response.go
 
 article
 ├── delivery
 │   └── http
 │       └── article_handler.go
 │   └── gRPC
 ├── repository //Encapsulated Implementation of Repository Interface
 │   └── mysql
 │       └── mysql_article.go
 └── manager //Encapsulated Implementation of Manager Interface
     └── article_manager.go
 config
 └── config.go
 middleware
 ├── apm.go
 └── cors.go 
