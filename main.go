package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Admin struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

type Team struct {
	gorm.Model
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	Founded int    `json:"founded"`
	Address string `json:"address"`
	City    string `json:"city"`
	Players []Player
}

type Player struct {
	gorm.Model
	Name     string `json:"name"`
	Height   int    `json:"height"`
	Weight   int    `json:"weight"`
	Position string `json:"position"`
	Number   int    `json:"number"`
	TeamID   uint   `json:"team_id"`
}

type Match struct {
	gorm.Model
	Date       string  `json:"date"`
	Time       string  `json:"time"`
	HomeTeamID uint    `json:"home_team_id"`
	AwayTeamID uint    `json:"away_team_id"`
	HomeTeam   Team    `gorm:"foreignKey:HomeTeamID"`
	AwayTeam   Team    `gorm:"foreignKey:AwayTeamID"`
	Results    []Result
}

type Result struct {
	gorm.Model
	MatchID    uint   `json:"match_id"`
	HomeScore  int    `json:"home_score"`
	AwayScore  int    `json:"away_score"`
	ScorerName string `json:"scorer_name"`
	GoalMinute int    `json:"goal_minute"`
}

var db *gorm.DB
var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func initDB() {
	dsn := "host=localhost user=postgres password=dawad123321 dbname=football port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	db.AutoMigrate(&Admin{}, &Team{}, &Player{}, &Match{}, &Result{})
	var count int64
	db.Model(&Admin{}).Count(&count)
	if count == 0 {
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		db.Create(&Admin{Username: "admin", Password: string(hash)})
		log.Println("Default admin created: username=admin, password=admin123")
	}
}

func login(c *gin.Context) {
	var input Admin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	var admin Admin
	if err := db.Where("username = ?", input.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}
	exp := time.Now().Add(2 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: input.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	})
	tokenStr, _ := token.SignedString(jwtKey)
	c.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tok := c.GetHeader("Authorization")
		if len(tok) < 8 || tok[:7] != "Bearer " {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}
		tok = tok[7:]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tok, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func createTeam(c *gin.Context) {
	var t Team
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&t)
	c.JSON(http.StatusOK, t)
}

func updateTeam(c *gin.Context) {
	id := c.Param("id")
	var t Team
	if db.First(&t, id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	var input Team
	c.ShouldBindJSON(&input)
	db.Model(&t).Updates(input)
	c.JSON(http.StatusOK, t)
}

func deleteTeam(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Team{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Team soft deleted"})
}

func getTeams(c *gin.Context) {
	var teams []Team
	db.Preload("Players").Find(&teams)
	c.JSON(http.StatusOK, teams)
}

func createPlayer(c *gin.Context) {
	var p Player
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var exist Player
	if err := db.Where("team_id = ? AND number = ?", p.TeamID, p.Number).First(&exist).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Player number already used in team"})
		return
	}
	db.Create(&p)
	c.JSON(http.StatusOK, p)
}

func createMatch(c *gin.Context) {
	var m Match
	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if m.HomeTeamID == m.AwayTeamID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Teams cannot be the same"})
		return
	}
	db.Create(&m)
	c.JSON(http.StatusOK, m)
}

func reportMatch(c *gin.Context) {
	var r Result
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&r)
	c.JSON(http.StatusOK, r)
}

func getReport(c *gin.Context) {
	var matches []Match
	db.Preload("HomeTeam").Preload("AwayTeam").Preload("Results").Find(&matches)
	type Report struct {
		MatchID     uint   `json:"match_id"`
		HomeTeam    string `json:"home_team"`
		AwayTeam    string `json:"away_team"`
		HomeScore   int    `json:"home_score"`
		AwayScore   int    `json:"away_score"`
		Status      string `json:"status"`
		TopScorer   string `json:"top_scorer"`
		HomeWins    int    `json:"home_total_wins"`
		AwayWins    int    `json:"away_total_wins"`
	}
	var reports []Report
	homeWins := map[uint]int{}
	awayWins := map[uint]int{}
	for _, m := range matches {
		hs, as := 0, 0
		scorers := map[string]int{}
		for _, r := range m.Results {
			hs = r.HomeScore
			as = r.AwayScore
			scorers[r.ScorerName]++
		}
		status := "Draw"
		if hs > as {
			status = "Home Win"
			homeWins[m.HomeTeamID]++
		} else if as > hs {
			status = "Away Win"
			awayWins[m.AwayTeamID]++
		}
		top := ""
		max := 0
		for n, c := range scorers {
			if c > max {
				max = c
				top = n
			}
		}
		reports = append(reports, Report{
			MatchID:   m.ID,
			HomeTeam:  m.HomeTeam.Name,
			AwayTeam:  m.AwayTeam.Name,
			HomeScore: hs,
			AwayScore: as,
			Status:    status,
			TopScorer: top,
			HomeWins:  homeWins[m.HomeTeamID],
			AwayWins:  awayWins[m.AwayTeamID],
		})
	}
	c.JSON(http.StatusOK, reports)
}

func main() {
	initDB()
	r := gin.Default()
	r.POST("/auth/login", login)
	r.GET("/teams", getTeams)
	auth := r.Group("/")
	auth.Use(authMiddleware())
	{
		auth.POST("/teams", createTeam)
		auth.PUT("/teams/:id", updateTeam)
		auth.DELETE("/teams/:id", deleteTeam)
		auth.POST("/players", createPlayer)
		auth.POST("/matches", createMatch)
		auth.POST("/matches/result", reportMatch)
		auth.GET("/reports", getReport)
	}
	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
