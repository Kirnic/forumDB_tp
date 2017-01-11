package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"
	"gopkg.in/gin-gonic/gin.v1"
)

var log = logging.MustGetLogger("main")
var format = logging.MustStringFormatter(`%{color} %{shortfunc} ▶ %{level:.5s} %{id:03x}%{color:reset} %{message}`)

func main() {
	logging.SetFormatter(format)
	config := loadConfig()
	dbmap := initDB(&config)
	defer dbmap.Map.Db.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	common := router.Group("/db/api/")
	{
		common.POST("clear/", dbmap.commonClear)
		common.GET("status/", dbmap.commonStatus)
	}
	forum := router.Group("/db/api/forum/")
	{
		forum.POST("create/", dbmap.forumCreate)
		forum.GET("details/", dbmap.forumDetails)
		forum.GET("listPosts/", dbmap.forumListPosts)
		forum.GET("listThreads/", dbmap.forumListThreads)
		forum.GET("listUsers/", dbmap.forumListUsers)
	}
	thread := router.Group("/db/api/thread/")
	{
		thread.POST("create/", dbmap.threadCreate)
		thread.GET("details/", dbmap.threadDetails)
		thread.POST("close/", dbmap.threadClose)
		thread.GET("list/", dbmap.threadList)
		thread.GET("listPosts/", dbmap.threadListPosts)
		thread.POST("open/", dbmap.threadOpen)
		thread.POST("remove/", dbmap.threadRemove)
		thread.POST("restore/", dbmap.threadRestore)
		thread.POST("subscribe/", dbmap.threadSubscribe)
		thread.POST("unsubscribe/", dbmap.threadUnsubscribe)
		thread.POST("update/", dbmap.threadUpdate)
		thread.POST("vote/", dbmap.threadVote)
	}
	post := router.Group("/db/api/post/")
	{
		post.POST("create/", dbmap.postCreate)
		post.GET("details/", dbmap.postDetails)
		post.GET("list/", dbmap.postList)
		post.POST("remove/", dbmap.postRemove)
		post.POST("restore/", dbmap.postRestore)
		post.POST("update/", dbmap.postUpdate)
		post.POST("vote/", dbmap.postVote)
	}
	user := router.Group("/db/api/user/")
	{
		user.POST("create/", dbmap.userCreate)
		user.GET("details/", dbmap.userDetails)
		user.POST("follow/", dbmap.userFollow)
		user.GET("listFollowers/", dbmap.userFollowersList)
		user.GET("listFollowing/", dbmap.userFollowingList)
		user.GET("listPosts/", dbmap.userListPosts)
		user.POST("unfollow/", dbmap.userUnfollow)
		user.POST("updateProfile/", dbmap.userUpdate)
	}

	err := router.Run(":" + config.PORT)
	errCheck(err)
}

func errCheck(err error) {
	if err != nil {
		log.Critical(err)
		panic(err)
	}
}

func loadConfig() Config {
	file, err := os.Open("config.json")
	errCheck(err)

	decoder := json.NewDecoder(file)
	conf := Config{}
	err = decoder.Decode(&conf)
	errCheck(err)

	return conf
}

func initDB(config *Config) *DB {
	connection := config.USER + ":" + config.PASS + "@unix(" + config.PATH + ")/" + config.DB + "?loc=Local"
	db, err := sql.Open(config.DIAL, connection)
	// db, err := sql.Open("mysql", "root:0072@/forumDB?charset=utf8")
	errCheck(err)
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{Encoding: "utf8", Engine: "InnoDB"}}
	return &DB{Map: dbmap}
}

// Config struct
type Config struct {
	DB   string
	DIAL string
	HOST string
	PORT string
	PATH string
	USER string
	PASS string
}

// DB wrapper
type DB struct {
	Map *gorp.DbMap
}

// Forum entity
type Forum struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	ShortName string `json:"short_name" db:"short_name"`
	User      string `json:"user" db:"user"`
}

// User entity
type User struct {
	About       string `json:"about" db:"about"`
	Email       string `json:"email" db:"email"`
	ID          int64  `json:"id" db:"id"`
	IsAnonymous bool   `json:"isAnonymous" db:"isAnonymous"`
	Name        string `json:"name" db:"name"`
	Username    string `json:"username" db:"username"`
}

// Post entity
type Post struct {
	Date          string `json:"date" db:"date"`
	Dislikes      int    `json:"dislikes" db:"dislikes"`
	Forum         string `json:"forum" db:"forum"`
	ID            int    `json:"id" db:"id"`
	IsApproved    bool   `json:"isApproved" db:"isApproved"`
	IsDeleted     bool   `json:"isDeleted" db:"isDeleted"`
	IsEdited      bool   `json:"isEdited" db:"isEdited"`
	IsHighlighted bool   `json:"isHighlighted" db:"isHighlighted"`
	IsSpam        bool   `json:"isSpam" db:"isSpam"`
	Likes         int    `json:"likes" db:"likes"`
	Message       string `json:"message" db:"message"`
	Parent        *int   `json:"parent" db:"parent"`
	Points        int    `json:"points" db:"points"`
	Thread        int    `json:"thread" db:"thread"`
	User          string `json:"user" db:"user"`
	FirstPath     int    `json:"first_path" db:"first_path"`
	LastPath      string `json:"last_path" db:"last_path"`
}

// Thread entity
type Thread struct {
	Date      string `json:"date" db:"date"`
	Dislikes  int    `json:"dislikes" db:"dislikes"`
	Forum     string `json:"forum" db:"forum"`
	ID        int    `json:"id" db:"id"`
	IsClosed  bool   `json:"isClosed" db:"isClosed"`
	IsDeleted bool   `json:"isDeleted" db:"isDeleted"`
	Likes     int    `json:"likes" db:"likes"`
	Message   string `json:"message" db:"message"`
	Points    int    `json:"points" db:"points"`
	Posts     int    `json:"posts" db:"posts"`
	Slug      string `json:"slug" db:"slug"`
	Title     string `json:"title" db:"title"`
	User      string `json:"user" db:"user"`
}

// Follow entity
type Follow struct {
	Follower  string `json:"follower" db:"follower"`
	Following string `json:"followee" db:"following"`
}

// UpdateUser entity
type UpdateUser struct {
	About string `json:"about"`
	User  string `json:"user"`
	Name  string `json:"name"`
}

// COMMON METHODS

func (db *DB) commonClear(c *gin.Context) {
	tables := []string{"forum", "post", "user", "thread", "follow", "subscription", "user_forum"}
	for _, table := range tables {
		db.Map.Exec(`truncate table ` + table)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": "OK"})
}

func (db *DB) commonStatus(c *gin.Context) {
	tables := []string{"forum", "post", "user", "thread"}
	response := gin.H{}
	for _, table := range tables {
		rowsCount, _ := db.Map.SelectInt(`select count(*) from ` + table)
		response[table] = rowsCount
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": response})
}

// FORUM METHODS

func (db *DB) forumSelect(shortName string, full bool) gin.H {
	forum := Forum{}
	db.Map.SelectOne(&forum, "select * from `forum` where `short_name` = ?", shortName)
	forumInfo := gin.H{"id": forum.ID, "name": forum.Name, "short_name": forum.ShortName, "user": forum.User}
	if full {
		forumInfo["user"] = db.userSelect(forum.User, true)
	}
	return forumInfo
}

func (db *DB) forumCreate(c *gin.Context) {
	forum := Forum{}
	c.BindJSON(&forum)

	db.Map.Exec("insert into `forum` (name, short_name, user) values(?, ?, ?)",
		forum.Name, forum.ShortName, forum.User)
	forumInfo := db.forumSelect(forum.ShortName, false)

	c.JSON(http.StatusOK, gin.H{"code": 0, "response": forumInfo})

}

func (db *DB) forumDetails(c *gin.Context) {
	forum := c.Query("forum")
	forumInfo := gin.H{}
	if related := c.Query("related"); related == "user" {
		forumInfo = db.forumSelect(forum, true)
	} else {
		forumInfo = db.forumSelect(forum, false)
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": forumInfo})
}

func (db *DB) forumListPosts(c *gin.Context) {
	related, user, forum, thread := c.Request.URL.Query()["related"], false, false, false
	for _, entity := range related {
		if entity == "user" {
			user = true
		} else if entity == "forum" {
			forum = true
		} else if entity == "thread" {
			thread = true
		}
	}

	query := "select * from `post` where `forum` = " + "\"" + c.Query("forum") + "\""
	if since := c.Query("since"); since != "" {
		query += " and `date` >= " + "\"" + since + "\""
	}
	query += " order by `date` " + c.DefaultQuery("order", "desc")
	if limit := c.Query("limit"); limit != "" {
		query += " limit " + limit
	}
	posts := []Post{}
	db.Map.Select(&posts, query)
	response := make([]gin.H, len(posts))
	for i, post := range posts {
		response[i] = gin.H{"date": post.Date, "dislikes": post.Dislikes, "forum": post.Forum, "id": post.ID,
			"isApproved": post.IsApproved, "isDeleted": post.IsDeleted, "isEdited": post.IsEdited,
			"isHighlighted": post.IsHighlighted, "isSpam": post.IsSpam, "likes": post.Likes, "message": post.Message,
			"parent": post.Parent, "points": post.Points, "thread": post.Thread, "user": post.User}
		if user {
			response[i]["user"] = db.userSelect(response[i]["user"].(string), true)
		}
		if forum {
			response[i]["forum"] = db.forumSelect(response[i]["forum"].(string), false)
		}
		if thread {
			response[i]["thread"] = db.threadSelect(response[i]["thread"].(int), true)
		}
	}
	// log.Debug(response)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": response})
}

func (db *DB) forumListThreads(c *gin.Context) {
	related, user, forum := c.Request.URL.Query()["related"], false, false
	for _, entity := range related {
		if entity == "user" {
			user = true
		} else if entity == "forum" {
			forum = true
		}
	}
	query := "select * from `thread` where `forum` = " + "\"" + c.Query("forum") + "\""
	if since := c.Query("since"); since != "" {
		query += " and `date` >= " + "\"" + since + "\""
	}
	query += " order by `date` " + c.DefaultQuery("order", "desc")
	if limit := c.Query("limit"); limit != "" {
		query += " limit " + limit
	}
	threads := []Thread{}
	db.Map.Select(&threads, query)
	response := make([]gin.H, len(threads))
	for i, thread := range threads {
		response[i] = gin.H{"date": thread.Date, "dislikes": thread.Dislikes, "forum": thread.Forum, "id": thread.ID,
			"isClosed": thread.IsClosed, "isDeleted": thread.IsDeleted, "likes": thread.Likes, "message": thread.Message,
			"points": thread.Points, "posts": thread.Posts, "slug": thread.Slug, "title": thread.Title, "user": thread.User}
		if user {
			response[i]["user"] = db.userSelect(response[i]["user"].(string), true)
		}
		if forum {
			response[i]["forum"] = db.forumSelect(response[i]["forum"].(string), false)
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": response})
}

func (db *DB) forumListUsers(c *gin.Context) {
	query := "select * from `user` where `email` IN (select distinct `user` from `post` where `forum` = " + "\"" + c.Query("forum") + "\")"
	if sinceID := c.Query("since_id"); sinceID != "" {
		query += " and `user`.`id` >= " + sinceID
	}
	query += " order by `user`.`name` " + c.DefaultQuery("order", "desc")
	if limit := c.Query("limit"); limit != "" {
		query += " limit " + limit
	}
	users := []User{}
	db.Map.Select(&users, query)
	response := make([]gin.H, len(users))
	for i, user := range users {
		var followers []string
		db.Map.Select(&followers, "select `follower` from `follow` where `following` = ?", user.Email)
		var following []string
		db.Map.Select(&following, "select `following` from `follow` where `follower` = ?", user.Email)
		var subscriptions []int
		db.Map.Select(&subscriptions, "select `thread` from `subscription` where `user` = ?", user.Email)
		response[i] = gin.H{"about": user.About, "email": user.Email, "followers": followers, "following": following,
			"id": user.ID, "isAnonymous": user.IsAnonymous, "name": user.Name, "subscriptions": subscriptions,
			"username": user.Username}
		if user.IsAnonymous {
			response[i]["username"] = nil
			response[i]["about"] = nil
			response[i]["name"] = nil
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": response})
}

// THREAD METHODS

func (db *DB) threadSelect(id int, full bool) gin.H {
	thread := Thread{}
	db.Map.SelectOne(&thread, "select * from `thread` where id = ?", id)

	threadInfo := gin.H{"date": thread.Date, "forum": thread.Forum, "id": thread.ID, "isClosed": thread.IsClosed,
		"isDeleted": thread.IsDeleted, "message": thread.Message, "slug": thread.Slug, "title": thread.Title, "user": thread.User}

	if full {
		threadInfo["posts"] = thread.Posts
		threadInfo["likes"] = thread.Likes
		threadInfo["dislikes"] = thread.Dislikes
		threadInfo["points"] = thread.Points
	}
	return threadInfo
}

func (db *DB) threadCreate(c *gin.Context) {
	thread := Thread{}
	c.BindJSON(&thread)

	result, _ := db.Map.Exec("insert into `thread` (`forum`, `user`, `title`, `isClosed`, `slug`, `date`, `message`, `IsDeleted`) values (?, ?, ?, ?, ?, ?, ?, ?)",
		thread.Forum, thread.User, thread.Title, thread.IsClosed, thread.Slug, thread.Date, thread.Message, thread.IsDeleted)
	id, _ := result.LastInsertId()

	threadInfo := db.threadSelect(int(id), false)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": threadInfo})
}

func (db *DB) threadDetails(c *gin.Context) {
	thread, _ := strconv.Atoi(c.Query("thread"))
	threadInfo := db.threadSelect(thread, true)

	for _, entity := range c.Request.URL.Query()["related"] {
		if entity == "user" {
			threadInfo[entity] = db.userSelect(threadInfo[entity].(string), true)
		} else if entity == "forum" {
			threadInfo[entity] = db.forumSelect(threadInfo[entity].(string), false)
		} else {
			c.JSON(http.StatusOK, gin.H{"code": 3, "response": "Bad request"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": threadInfo})
}

func (db *DB) threadClose(c *gin.Context) {
	var thread struct {
		ID int `json:"thread"`
	}
	c.BindJSON(&thread)
	db.Map.Exec("update `thread` set `isClosed` = true where `id` = ?", thread.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": thread})
}

func (db *DB) threadList(c *gin.Context) {
	query := "select * from `thread` where "
	if related := c.Query("forum"); related != "" {
		query += "`forum` = " + "\"" + related + "\""
	} else if related = c.Query("user"); related != "" {
		query += "`user` = " + "\"" + related + "\""
	}
	if since := c.Query("since"); since != "" {
		query += " and `date` >= " + "\"" + since + "\""
	}
	if order := c.DefaultQuery("order", "desc"); order != "" {
		query += " order by `date` " + order
	}
	if limit := c.Query("limit"); limit != "" {
		query += " limit " + limit
	}

	threads := []Thread{}
	db.Map.Select(&threads, query)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": threads})
}

func (db *DB) threadListPosts(c *gin.Context) {
	posts := []Post{}
	query := "select * from `post` where `thread` = " + c.Query("thread")
	if since := c.Query("since"); since != "" {
		query += " and `date` >= " + "\"" + since + "\""
	}
	order := c.Query("order")

	sort := c.Query("sort")
	if sort != "parent_tree" {
		if sort == "" {
			query += " order by `date` " + c.DefaultQuery("order", "desc")
			if limit := c.Query("limit"); limit != "" {
				query += " limit " + limit
			}

		} else if sort == "flat" {
			query += " order by `date` " + c.DefaultQuery("order", "desc")
			if limit := c.Query("limit"); limit != "" {
				query += " limit " + limit
			}
		} else if sort == "tree" {
			if order == "desc" {
				query += "order by first_path desc, last_path asc "
				if limit := c.Query("limit"); limit != "" {
					query += " limit " + limit
				}
			}
			if order == "asc" {
				query += "order by first_path asc, last_path asc "
				if limit := c.Query("limit"); limit != "" {
					query += " limit " + limit
				}
			}
		}
		db.Map.Select(&posts, query)
		c.JSON(http.StatusOK, gin.H{"code": 0, "response": posts})
	}
	if sort == "parent_tree" {
		postsTemp := []Post{}
		resultList := []Post{}

		query += "order by first_path asc"
		query += ", last_path asc"
		limit := c.Query("limit")
		db.Map.Select(&postsTemp, query)
		currParFirstPath := -1
		limitInt, _ := strconv.Atoi(limit)
		counter := 0
		for i := 0; i < len(postsTemp); i++ {

			if currParFirstPath != postsTemp[i].FirstPath {
				currParFirstPath = postsTemp[i].FirstPath
				counter++
			}
			if counter > limitInt {
				break
			}

			resultList = append(resultList, postsTemp[i])
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "response": resultList})
	}
}

func (db *DB) threadOpen(c *gin.Context) {
	var thread struct {
		ID int `json:"thread"`
	}
	c.BindJSON(&thread)
	db.Map.Exec("update `thread` set `isClosed` = false where `id` = ?", thread.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": thread})
}

func (db *DB) threadRemove(c *gin.Context) {
	var thread struct {
		ID int `json:"thread"`
	}
	c.BindJSON(&thread)
	db.Map.Exec("update `thread` set `isDeleted` = true, `posts` = 0 where `id` = ?", thread.ID)
	db.Map.Exec("update `post` set `isDeleted` = true where thread = ?", thread.ID)

	c.JSON(http.StatusOK, gin.H{"code": 0, "response": thread})
}

func (db *DB) threadRestore(c *gin.Context) {
	var thread struct {
		ID int `json:"thread"`
	}
	c.BindJSON(&thread)
	posts, _ := db.Map.SelectInt("select count(id) from `post` where `thread` = ?", thread.ID)
	db.Map.Exec("update `thread` set `isDeleted` = false, `posts` = ? where `id` = ?", posts, thread.ID)
	db.Map.Exec("update `post` set `isDeleted` = false where `thread` = ?", thread.ID)

	c.JSON(http.StatusOK, gin.H{"code": 0, "response": thread})
}

func (db *DB) threadSubscribe(c *gin.Context) {
	var subs struct {
		ID   int    `json:"thread"`
		User string `json:"user"`
	}
	c.BindJSON(&subs)
	db.Map.Exec("insert into `subscription` (user, thread) values (?, ?)", subs.User, subs.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": subs})
}

func (db *DB) threadUnsubscribe(c *gin.Context) {
	var subs struct {
		ID   int    `json:"thread"`
		User string `json:"user"`
	}
	c.BindJSON(&subs)
	db.Map.Exec("delete from `subscription` where `user` = ? and `thread` = ?", subs.User, subs.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": subs})
}

func (db *DB) threadUpdate(c *gin.Context) {
	var updateThrd struct {
		Message string `json:"message"`
		Slug    string `json:"slug"`
		ID      int    `json:"thread"`
	}
	c.BindJSON(&updateThrd)
	db.Map.Exec("update `thread` set `message` = ?, `slug` = ? where `id` = ?",
		updateThrd.Message, updateThrd.Slug, updateThrd.ID)

	threadInfo := db.threadSelect(updateThrd.ID, true)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": threadInfo})
}

func (db *DB) threadVote(c *gin.Context) {
	var thread struct {
		Vote int `json:"vote"`
		ID   int `json:"thread"`
	}
	c.BindJSON(&thread)
	if thread.Vote > 0 {
		db.Map.Exec("update `thread` set `likes` = likes + 1, `points` = points + 1 where `id` = ?", thread.ID)
	} else if thread.Vote < 0 {
		db.Map.Exec("update `thread` set `dislikes` = dislikes + 1, `points` = points - 1 where `id` = ?", thread.ID)
	}
	threadInfo := db.threadSelect(thread.ID, true)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": threadInfo})
}

// POST METHODS

const sizeOfPath int = 3

func intCapacity(num int) int {
	size := 0
	for num > 0 {
		num = num / 10
		size++
	}
	return size
}

func makePath(number int) string {
	var mathPath string
	for i := sizeOfPath - intCapacity(number); i > 0; i-- {
		mathPath += "0"
	}
	numStr := strconv.Itoa(number)
	mathPath += numStr
	return mathPath
}

func (db *DB) postSelect(id int) gin.H {
	post := Post{}
	if err := db.Map.SelectOne(&post, "select * from `post` WHERE `id` = ?", id); err == nil {
		return gin.H{"date": post.Date, "dislikes": post.Dislikes, "forum": post.Forum, "id": post.ID,
			"isApproved": post.IsApproved, "isDeleted": post.IsDeleted, "isEdited": post.IsEdited,
			"isHighlighted": post.IsHighlighted, "isSpam": post.IsSpam, "likes": post.Likes, "message": post.Message,
			"parent": post.Parent, "points": post.Points, "thread": post.Thread, "user": post.User, "first_path": 0, "last_path": ""}
	}
	return nil
}

func (db *DB) postCreate(c *gin.Context) {
	post := Post{}
	c.BindJSON(&post)
	result, _ := db.Map.Exec("insert into `post` (date, forum, isApproved, isDeleted, isEdited, isHighlighted, isSpam, message, parent, thread, user) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		post.Date, post.Forum, post.IsApproved, post.IsDeleted, post.IsEdited, post.IsHighlighted,
		post.IsSpam, post.Message, post.Parent, post.Thread, post.User)
	id, _ := result.LastInsertId()

	if post.Parent == nil {
		db.Map.Exec("update `post` set `first_path` = ? where `id` = ?", id, id)
	} else {
		tempPost := Post{}
		db.Map.SelectOne(&tempPost, "select `first_path`, `last_path` from `post` where `id` = ?", post.Parent)
		parentFirstPath := tempPost.FirstPath
		parentLastPath := tempPost.LastPath
		if parentLastPath == "" {
			i := id
			var i64 int
			i64 = int(i)
			mathPathID := "."
			mathPathID += makePath(i64)
			db.Map.Exec("update `post` set `first_path` = ?, `last_path` = ? where `id` = ?",
				parentFirstPath, mathPathID, id)
		} else {
			parentLastPath += "."
			i := id
			var i64 int
			i64 = int(i)
			mathPathID := makePath(i64)
			parentLastPath += mathPathID
			db.Map.Exec("update `post` set `first_path` = ?, `last_path` = ? where `id` = ?",
				parentFirstPath, parentLastPath, id)
		}
	}
	db.Map.Exec("update `thread` set `posts` = `posts` + 1 where id = ?", post.Thread)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": gin.H{"date": post.Date, "forum": post.Forum,
		"id": id, "isApproved": post.IsApproved, "isDeleted": post.IsDeleted, "isEdited": post.IsEdited,
		"isHighlighted": post.IsHighlighted, "isSpam": post.IsSpam, "message": post.Message,
		"parent": post.Parent, "thread": post.Thread, "user": post.User}})
}

func (db *DB) postDetails(c *gin.Context) {
	a := (c.Query("post"))
	post, _ := strconv.Atoi(a)

	if response := db.postSelect(post); response != nil {

		for _, entity := range c.Request.URL.Query()["related"] {
			if entity == "user" {
				response[entity] = db.userSelect(response[entity].(string), true)
			} else if entity == "thread" {
				response[entity] = db.threadSelect(response[entity].(int), true)
			} else if entity == "forum" {
				response[entity] = db.forumSelect(response[entity].(string), false)
			}
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "response": response})
	} else {
		c.JSON(http.StatusOK, gin.H{"code": 1, "response": "Post not found"})
	}
}

func (db *DB) postList(c *gin.Context) {

	query := "SELECT * FROM post WHERE"
	if forum := c.Query("forum"); forum != "" {
		query += " forum = " + "\"" + forum + "\""
	} else {
		query += " thread = " + c.Query("thread")
	}
	if since := c.Query("since"); since != "" {
		query += " AND date >= " + "\"" + since + "\""
	}
	query += " ORDER BY date " + c.DefaultQuery("order", "desc")
	if limit := c.Query("limit"); limit != "" {
		query += " LIMIT " + limit
	}
	var posts []Post
	db.Map.Select(&posts, query)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": posts})
}

func (db *DB) postRemove(c *gin.Context) {
	var post struct {
		ID int `json:"post"`
	}
	c.BindJSON(&post)
	db.Map.Exec("update `post` set `isDeleted` = true where `id` = ? ", post.ID)
	thread, _ := db.Map.SelectInt("select `thread` from `post` where `id` = ?", post.ID)
	db.Map.Exec("update thread set `posts` = `posts` - 1 where `id` = ?", thread)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": post})

}

func (db *DB) postRestore(c *gin.Context) {
	var post struct {
		ID int `json:"post"`
	}
	c.BindJSON(&post)
	db.Map.Exec("update `post` set `isDeleted` = false where `id` = ? ", post.ID)
	thread, _ := db.Map.SelectInt("select `thread` from `post` where `id` = ?", post.ID)
	db.Map.Exec("update thread set `posts` = `posts` + 1 where `id` = ?", thread)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": post})
}

func (db *DB) postUpdate(c *gin.Context) {
	var post struct {
		ID      int    `json:"post"`
		Message string `json:"message"`
	}
	c.BindJSON(&post)
	db.Map.Exec("update `post` set `message` = ? where `id` = ?", post.Message, post.ID)

	postInfo := db.postSelect(post.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": postInfo})
}

func (db *DB) postVote(c *gin.Context) {
	var post struct {
		ID   int `json:"post"`
		Vote int `json:"vote"`
	}
	c.BindJSON(&post)
	if post.Vote > 0 {
		db.Map.Exec("update `post` set `likes` = `likes` + 1, `points` = `points` + 1 where `id` = ?", post.ID)
	} else {
		db.Map.Exec("update `post` set `dislikes` = `dislikes` + 1, `points` = `points` - 1 where `id` = ?", post.ID)
	}
	postInfo := db.postSelect(post.ID)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": postInfo})
}

// USER METHODS

func (db *DB) userSelect(email string, full bool) gin.H {
	user := User{}
	db.Map.SelectOne(&user, "select * from `user` where `email` = ?", email)

	userInfo := gin.H{"about": user.About, "id": user.ID, "name": user.Name,
		"username": user.Username, "email": user.Email, "isAnonymous": user.IsAnonymous}
	if user.IsAnonymous {
		userInfo["about"] = nil
		userInfo["name"] = nil
		userInfo["username"] = nil
	}
	if full {
		var follower []string
		var following []string
		var subs []int

		db.Map.Select(&follower, "select `follower` from `follow` where `following` = ?", email)
		db.Map.Select(&following, "select `following` from `follow` where `follower` = ?", email)
		db.Map.Select(&subs, "select `thread` from `subscription` where `user` = ?", email)

		userInfo["followers"] = follower
		userInfo["following"] = following
		userInfo["subscriptions"] = subs
	}
	return userInfo
}

func (db *DB) userCreate(c *gin.Context) {
	user := User{}
	c.BindJSON(&user)

	_, err := db.Map.Exec("insert into `user` (about, name, username, isAnonymous, email) values(?, ?, ?, ?, ?)",
		user.About, user.Name, user.Username, user.IsAnonymous, user.Email)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 5, "response": "user is already exists"})
	} else {
		userInfo := db.userSelect(user.Email, false)
		c.JSON(http.StatusOK, gin.H{"code": 0, "response": userInfo})
	}
}

func (db *DB) userDetails(c *gin.Context) {
	userInfo := db.userSelect(c.Query("user"), true)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": userInfo})
}

func (db *DB) userFollow(c *gin.Context) {
	fol := Follow{}
	c.BindJSON(&fol)

	db.Map.Exec("insert into `follow` (follower, following) values(?, ?)", fol.Follower, fol.Following)
	userInfo := db.userSelect(fol.Follower, true)

	c.JSON(http.StatusOK, gin.H{"code": 0, "response": userInfo})
}

func (db *DB) userFollowersList(c *gin.Context) {
	user := c.Query("user")
	limit := c.Query("limit")
	order := c.Query("order")
	since := c.Query("since_id")

	query := "select `follower` from `follow` join `user` on `follower` = `email` where `following` = ? "
	if since != "" {
		query += "and `id` >= ? "
	}
	if order != "asc" {
		order = "desc"
	}
	query += " order by `follower` " + order
	if limit != "" {
		query += " limit " + limit
	}
	var followers []string
	if since != "" {
		db.Map.Select(&followers, query, user, since)
	} else {
		db.Map.Select(&followers, query, user)
	}

	followList := make([]gin.H, len(followers))
	for i, flw := range followers {
		followList[i] = db.userSelect(flw, true)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "response": followList})
}

func (db *DB) userFollowingList(c *gin.Context) {
	user := c.Query("user")
	limit := c.Query("limit")
	order := c.Query("order")
	since := c.Query("since_id")

	query := "select `following` from `follow` join `user` on `following` = `email` where `follower` = ? "
	if since != "" {
		query += "and `id` >= ? "
	}
	if order != "asc" {
		order = "desc"
	}
	query += " order by `following` " + order
	if limit != "" {
		query += " limit " + limit
	}
	var following []string
	if since != "" {
		db.Map.Select(&following, query, user, since)
	} else {
		db.Map.Select(&following, query, user)
	}

	followList := make([]gin.H, len(following))
	for i, flw := range following {
		followList[i] = db.userSelect(flw, true)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "response": followList})
}

func (db *DB) userUnfollow(c *gin.Context) {
	unfol := Follow{}
	c.BindJSON(&unfol)

	db.Map.Exec("delete from `follow` where `follower` = ? and `following` = ?", unfol.Follower, unfol.Following)
	userInfo := db.userSelect(unfol.Follower, true)

	c.JSON(http.StatusOK, gin.H{"code": 0, "response": userInfo})
}

func (db *DB) userListPosts(c *gin.Context) {
	query := "select * from `post` where `user` = " + "\"" + c.Query("user") + "\""
	if since := c.Query("since"); since != "" {
		query += " and `date` >= " + "\"" + since + "\""
	}
	query += " order by `date` " + c.DefaultQuery("order", "desc")
	if limit := c.Query("limit"); limit != "" {
		query += " limit " + limit
	}
	var posts []Post
	db.Map.Select(&posts, query)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": posts})
}

func (db *DB) userUpdate(c *gin.Context) {
	params := UpdateUser{}
	c.BindJSON(&params)
	db.Map.Exec("update `user` set `about` = ?, `name` = ? where `email` = ?", params.About, params.Name, params.User)
	userInfo := db.userSelect(params.User, true)
	c.JSON(http.StatusOK, gin.H{"code": 0, "response": userInfo})
}
