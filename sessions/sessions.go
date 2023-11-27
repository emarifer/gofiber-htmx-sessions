package sessions

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/emarifer/gofiber-htmx-sessions/db"
	"github.com/emarifer/gofiber-htmx-sessions/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"

	_ "github.com/mattn/go-sqlite3"
)

var store *session.Store

func createStorage() *sqlite3.Storage {
	return sqlite3.New(sqlite3.Config{
		Database:        "./fiber.db",
		Table:           "sessions",
		Reset:           false,
		GCInterval:      10 * time.Second,
		MaxOpenConns:    100,
		MaxIdleConns:    100,
		ConnMaxLifetime: 1 * time.Second,
	})
}

// Init sessions store
func InitSessionsStore() {
	store = session.New(session.Config{
		Storage:    createStorage(),
		Expiration: 1 * time.Hour,
		KeyLookup:  "cookie:myapp_session",
	})
}

func CreateUserSession(c *fiber.Ctx, uid string) error {
	// Get or create session
	s, _ := store.Get(c)
	// fmt.Println(s.Fresh())

	// If this is a new session
	if s.Fresh() {
		// Get session ID
		sid := s.ID()

		//Get user ID
		// uid := c.Params("uid")

		// Save session data
		s.Set("uid", uid)
		s.Set("sid", sid)
		s.Set("ip", c.Context().RemoteIP().String())
		s.Set("login", time.Unix(time.Now().Unix(), 0).UTC().String())
		s.Set("ua", string(c.Request().Header.UserAgent()))

		err := s.Save()
		if err != nil {
			// log.Println(err)
			return err
		}

		// Save user reference
		stmt, err := db.Db.Prepare(`UPDATE sessions SET u = ? WHERE k = ?`)
		if err != nil {
			// log.Println(err)
			return err
		}

		_, err = stmt.Exec(uid, sid)
		if err != nil {
			// log.Println(err)
			return err
		}
	}

	return nil
}

func GetUserSessionData(c *fiber.Ctx) (*models.Account, error) {
	// Get current session
	s, _ := store.Get(c)
	// fmt.Println(s.Keys())

	// If there is a valid session
	if len(s.Keys()) > 0 {
		sid := s.ID()
		// From the session that is started we obtain the user id
		uid := s.Get("uid").(string)
		// Then with its uid we get the user data
		user := new(models.User)
		user.ID = uid
		recoveredUser, err := user.GetUserById()
		if err != nil {
			return nil, err
		}

		// Get profile info
		U := &models.Account{
			Email:    recoveredUser.Email,
			Username: recoveredUser.Username,
			Session:  sid,
		}

		// Get sessions list
		rows, err := db.Db.Query(`SELECT v, e FROM sessions WHERE u = ?`, uid)
		if err != nil {
			log.Println(err)
		}

		defer rows.Close()

		// Loop through sessions
		for rows.Next() {
			var (
				data       = []byte{}
				exp  int64 = 0
			)
			if err := rows.Scan(&data, &exp); err != nil {
				log.Println(err)
				return nil, err
			}

			// If session isn't expired
			if exp > time.Now().Unix() {
				// Decode session data
				gd := gob.NewDecoder(bytes.NewBuffer(data))
				dm := make(map[string]interface{})
				if err := gd.Decode(&dm); err != nil {
					log.Println(err)
					return nil, err
				}

				// Append session
				U.Sessions = append(
					U.Sessions,
					models.UserSession{
						SID:    dm["sid"].(string),
						IP:     dm["ip"].(string),
						Login:  dm["login"].(string),
						Expiry: time.Unix(exp, 0).UTC().String(),
						UA:     dm["ua"].(string),
					},
				)
			}
		}

		return U, nil
	}

	return nil, nil
}

func RemoveUserSession(c *fiber.Ctx) (bool, error) {
	//Get session ID
	sid := c.Query("sid")
	// fmt.Println("SID: ", sid)

	// Get current session
	s, _ := store.Get(c)
	// fmt.Println(s.Fresh())

	// Check session ID
	if len(sid) > 0 {
		// Get requested session
		data, err := store.Storage.Get(sid)
		if err != nil {
			return false, err
		}

		// Decode requested session data
		gd := gob.NewDecoder(bytes.NewBuffer(data))
		dm := make(map[string]interface{})
		if err := gd.Decode(&dm); err != nil {
			return false, err
		}

		// If it belongs to current user destroy requested session
		if s.Get("uid").(string) == dm["uid"] {
			store.Storage.Delete(sid)
		}

		return false, nil
	} else {
		// Destroy current session
		s.Destroy()
	}

	return true, nil
}
