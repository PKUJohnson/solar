# Preparation
In each project's models/define.go file, add 
```go
var (
	db *gorm.DB
	dw *std.DBExtension
)

func InitModel(config std.ConfigMysql) {
	db = toolkit.CreateDB(config)
	dw = std.NewDBWrapper(db)
	std.LogInfoLn("start init mysql model")
	std.LogInfoLn("end init mysql model")
}

func DB() *std.DBExtension {
	return dw
}

func Session() *gorm.DB {
	return db.Begin()
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

```

# Examples 
## GetOne
```go
query := models.Topic{TopicId: topicId}
result := models.Topic{}

if found, err := models.DB().GetOne(&result, query); !found {
	//not found
    if err != nil {
    	// has error
        return err
    }
}

```

## GetList
```go
query := models.InstInfo{
    Valid:1,
}
result:= make([]*models.InstInfo, 0)

if err := models.DB().GetList(&result, query); err != nil{
    // error handling
    return err
}
```

```go
tids := []int{1, 2, 3, 4}
result := []*models.TuChart{}

if err :=models.DB().GetList(&result, "valid = 1 and tid in (?)", tids); err != nil{
    //error handling
    return err
}
```

## GetOrderedList
It shares the same usage of GetList, except that there is one more order field.

```go
query := models.InstInfo{
    Valid:1,
}
result:= make([]*models.InstInfo, 0)

if err := models.DB().GetList(&result,"create_time desc", query); err != nil{
    // error handling
    return err
}
```

## GetPageRangeList
```go
result := []*models.MpFeedInfo{}

if err :=models.DB().GetPageRangeList(&result, "update_time asc", limit, offset,
        "valid = 1 and update_time > ? and update_time < ?", startTime, endTime);err != nil{
    return err
}
```

## SaveOne
Update All Fields, the object's primary_key, defined in gorm format definition, must have value.

```go
instInfo  := models.InstInfo{
    InstId:req.InstId,
    Name:req.Name,
    Contact:req.Contact,
    Address:req.Address,
    Email:req.Email,
    Phone:req.Phone,
    OwnerId: req.OwnerId,
    Valid: req.Valid,
}

if err := models.DB().SaveOne(&instInfo); err != nil{
    // error handling
    return err
}
```

## Update
Update partial Fields, if attrs is an object, it will ignore default value field; if attrs is map, it will ignore unchanged field.

```go
//todo: add later
```

## ExecuteSql

```go
sql := "SELECT ta.*, tac.content FROM tucloud.tu_article ta inner join tucloud.tu_article_content tac on ta.article_id = tac.article_id and ta.valid = 1 and ta.article_id in (?)"
result := []*models.TuArticleWithContent{}

if err:= models.DB().ExecSql(&result, sql, tuArticleIds);err != nil{
    // error handling
    return err
}

```

