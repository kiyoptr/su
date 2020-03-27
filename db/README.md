This package provides gorm model and db connection handling.

The `DbTypeGetter` and `DbConnGetter` variables in this package must be set before any calls to `Open`.

Use `FindAll` and `FindSingle` for querying data and `BaseModel` as an anonymous field in all your models.
Also implement `Model` interface for all models.

Before doing any query or modification to db, you must call `CheckModelTables` once to make sure all tables of the
models are created. To define a model, use `DefineModel` in `init` functions.