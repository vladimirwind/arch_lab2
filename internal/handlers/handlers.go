package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
type UserMask struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
type Report struct {
	Name   string `json:"name"`
	UserId int    `json:"user_id"`
}

func connectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "admin:mai@tcp(mariadb:3306)/maidb")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateUser godoc
//
//	@Summary		Creates a new user
//	@Description	create a new user
//	@Accept			json
//	@Produce		json
//
//	@Tags			mai lab API
//
//	@Param			user_data	body		User	true	"User Data"
//	@Success		200			{object}	string
//	@Failure		400			{object}	gin.H
//	@Failure		404			{object}	gin.H
//	@Failure		500			{object}	gin.H
//	@Router			/user/create [post]
func CreateUser(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	var usr User
	err = json.Unmarshal(jsonData, &usr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	if len(usr.Login) == 0 || len(usr.Password) == 0 || len(usr.Name) == 0 || len(usr.Surname) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error ": "not all fields are provided"})
		return
	}
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()
	_, err = db.Exec("insert into users (user_name, user_surname, user_login, user_password) "+
		"VALUES (?,?,?,?);", usr.Name, usr.Surname, usr.Login, usr.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	c.JSON(200, "user has been created")
}

// FindUserByLogin godoc
//
//	@Summary		Find User By Login
//	@Description	Find User By Login
//	@Accept			json
//	@Produce		json
//
//	@Tags			mai lab API
//
//	@Param			user_log	path		string	true	"User Login"
//	@Success		200			{object}	[]UserMask
//	@Failure		400			{object}	gin.H
//	@Failure		404			{object}	gin.H
//	@Failure		500			{object}	gin.H
//	@Router			/user/findLogin/{user_log} [get]
func FindUserByLogin(c *gin.Context) {
	usr_log := c.Params.ByName("user_log")
	if len(usr_log) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error ": "login not provided"})
		return
	}
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT user_name, user_surname from users WHERE user_login = ?", usr_log)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	defer rows.Close()
	var res []UserMask
	for rows.Next() {
		var tmp UserMask
		err := rows.Scan(&tmp.Name, &tmp.Surname)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
			return
		}
		res = append(res, tmp)
	}
	if len(res) == 0 {
		c.JSON(http.StatusNotFound, "user not found")
		return
	}
	c.JSON(200, gin.H{"found user(s) ": res})
}

// FindUserByMask godoc
//
//	@Summary		Find User By Mask
//	@Description	Find User By Mask
//	@Accept			json
//	@Produce		json
//
//	@Tags			mai lab API
//
//	@Param			user_log	body		UserMask	true	"User Data with mask"
//	@Success		200			{object}	[]UserMask
//	@Failure		400			{object}	gin.H
//	@Failure		404			{object}	gin.H
//	@Failure		500			{object}	gin.H
//	@Router			/user/findMask [post]
func FindUserByMask(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()
	var usr UserMask
	err = json.Unmarshal(jsonData, &usr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	var res []UserMask
	rows, err := db.Query("SELECT user_name, user_surname FROM users where user_name LIKE ? AND user_surname LIKE ?",
		strings.Replace(usr.Name, "*", "%", -1), strings.Replace(usr.Surname, "*", "%", -1))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	for rows.Next() {
		var tmp UserMask
		err := rows.Scan(&tmp.Name, &tmp.Surname)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
			return
		}
		res = append(res, tmp)
	}
	if len(res) == 0 {
		c.JSON(404, "no users found with this mask")
		return
	}
	c.JSON(200, gin.H{"found user(s) ": res})
}

// CreateReport godoc
//
//	@Summary		Create New Report
//	@Description	Create New Report
//	@Accept			json
//	@Produce		json
//
//	@Tags			mai lab API
//
//	@Param			user_log	body		Report	true	"Report's data"
//	@Success		200			{object}	string
//	@Failure		400			{object}	gin.H
//	@Failure		404			{object}	gin.H
//	@Failure		500			{object}	gin.H
//	@Router			/report/create [post]
func CreateReport(c *gin.Context) {
	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	var rep Report
	err = json.Unmarshal(jsonData, &rep)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()
	_, err = db.Exec("insert into reports (report_name, user_id) "+
		"VALUES (?, ?);", rep.Name, rep.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	c.JSON(200, "report has been created")
}

// CreateConference godoc
//
//	@Summary		Create New Conference
//	@Description	Create New Conference
//	@Accept			json
//	@Produce		json
//
//	@Tags			mai lab API
//
//	@Param			conference_name	path		string	true	"conference name"
//	@Success		200				{object}	string
//	@Failure		400				{object}	gin.H
//	@Failure		404				{object}	gin.H
//	@Failure		500				{object}	gin.H
//	@Router			/conference/create/{conference_name} [post]
func CreateConference(c *gin.Context) {
	conf := c.Params.ByName("conference_name")
	if len(conf) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error ": "conference name not provided"})
		return
	}
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()
	_, err = db.Exec("insert into conferences (conference_name) "+
		"VALUES (?);", conf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	c.JSON(200, "conference has been created")
}

// GetAllReports godoc
//
//	@Summary		Get All Reports
//	@Description	Get All Reports
//	@Accept			json
//	@Produce		json
//
//	@Tags			mai lab API
//
//	@Success		200	{object}	[]string
//	@Failure		400	{object}	gin.H
//	@Failure		404	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/report/getAll [get]
func GetAllReports(c *gin.Context) {
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()
	var res []string
	rows, err := db.Query("SELECT report_name FROM reports")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
			return
		}
		res = append(res, tmp)
	}
	if len(res) == 0 {
		c.JSON(404, "no users found with this mask")
		return
	}
	c.JSON(200, gin.H{"found reports(s) ": res})
}

// AddReport godoc
//
//	@Summary		Add New Report
//	@Description	Add New Report
//	@Accept			json
//	@Produce		json
//
//	@Tags			mai lab API
//
//	@Param			conference_id	path		string	true	"conference id"
//	@Param			report_id		path		string	true	"report id"
//	@Success		200				{object}	string
//	@Failure		400				{object}	gin.H
//	@Failure		404				{object}	gin.H
//	@Failure		500				{object}	gin.H
//	@Router			/conference/addReport/{conference_id}/{report_id}/ [post]
func AddReport(c *gin.Context) {
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()
	rep_id := c.Params.ByName("report_id")
	cf_id := c.Params.ByName("conference_id")
	if len(rep_id) == 0 || len(cf_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error ": "id not provided"})
		return
	}
	_, err = db.Exec("UPDATE reports "+
		"SET conference_id=? "+
		"WHERE report_id=?;", cf_id, rep_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	c.JSON(200, "report has been added to conference")
}

// GetAllReportsInConfs godoc
//
//	@Summary		Get All Reports In Conference
//	@Description	Get All Reports In Conference
//	@Accept			json
//	@Produce		json
//	@Param			conference_id	path	string	true	"conference id"
//
//	@Tags			mai lab API
//
//	@Success		200	{object}	[]string
//	@Failure		400	{object}	gin.H
//	@Failure		404	{object}	gin.H
//	@Failure		500	{object}	gin.H
//	@Router			/conference/getAllReports/{conference_id}/ [get]
func GetAllReportsInConf(c *gin.Context) {
	db, err := connectDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}
	defer db.Close()
	cf_id := c.Params.ByName("conference_id")
	if len(cf_id) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error ": "id not provided"})
		return
	}
	rows, err := db.Query("SELECT report_name FROM reports WHERE conference_id = ?", cf_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
		return
	}
	var res []string
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error ": err.Error()})
			return
		}
		res = append(res, tmp)
	}
	if len(res) == 0 {
		c.JSON(404, "no users found with this mask")
		return
	}
	c.JSON(200, gin.H{"found reports(s) in conference": res})
}
