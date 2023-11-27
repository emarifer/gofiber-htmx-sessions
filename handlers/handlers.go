package handlers

import (
	"log"
	"strings"

	"github.com/emarifer/gofiber-htmx-sessions/models"
	"github.com/emarifer/gofiber-htmx-sessions/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/sujit-baniya/flash"
	"golang.org/x/crypto/bcrypt"
)

/********** Handlers for rendering views **********/

// Redirect root URL to Login Page
func RedirectToLogin(c *fiber.Ctx) error {
	return c.Redirect("/login", fiber.StatusMovedPermanently)
}

// Render Login Page with success/error messages
func HandleViewLogin(c *fiber.Ctx) error {

	return c.Render("login", fiber.Map{
		"Page":           "Login",
		"Error":          flash.Get(c)["error"],
		"ErrorMessage":   flash.Get(c)["errorMessage"],
		"Success":        flash.Get(c)["success"],
		"SuccessMessage": flash.Get(c)["successMessage"],
	})
}

// Render Register Page with success/error messages
func HandleViewRegister(c *fiber.Ctx) error {

	return c.Render("register", fiber.Map{
		"Page":           "Register",
		"Error":          flash.Get(c)["error"],
		"ErrorMessage":   flash.Get(c)["errorMessage"],
		"Success":        flash.Get(c)["success"],
		"SuccessMessage": flash.Get(c)["successMessage"],
	})
}

// Render Profile Page (protected page)
func HandleViewProfile(c *fiber.Ctx) error {
	userProfileData, err := sessions.GetUserSessionData(c)
	if err != nil {
		log.Println(err)
		fm := fiber.Map{
			"error":        true,
			"errorMessage": "Something went wrong: cannot recover the session!!",
		}

		return flash.WithError(c, fm).Redirect("/login")
	}

	if userProfileData != nil {
		return c.Render("profile", fiber.Map{
			"Page":           "Profile",
			"email":          userProfileData.Email, // ↓ USER DATA ↓
			"username":       userProfileData.Username,
			"CurrentSession": userProfileData.Session,
			"Sessions":       userProfileData.Sessions,
			"Error":          flash.Get(c)["error"], // ↓ FLASH MESSAGES ↓
			"ErrorMessage":   flash.Get(c)["errorMessage"],
			"Success":        flash.Get(c)["success"],
			"SuccessMessage": flash.Get(c)["successMessage"],
		})
	}

	fm := fiber.Map{
		"error":        true,
		"errorMessage": "You are no longer authenticated!!",
	}

	return flash.WithError(c, fm).Redirect("/login")
}

/********** Singin/Singup & Sessions Handlers **********/

// Handle SignIn & start Session
func HandleSigninUser(c *fiber.Ctx) error {
	email := strings.Trim(c.FormValue("email"), " ")
	password := strings.Trim(c.FormValue("password"), " ")

	newUser := new(models.User)
	newUser.Email = email

	fm := fiber.Map{
		"error":        true,
		"errorMessage": "Incorrect email or password!!",
	}

	// We recover the user from their email
	storedUser, err := newUser.GetUserByEmail()
	if err != nil {
		log.Println(err)

		return flash.WithError(c, fm).Redirect("/login")
	}

	// Compare the stored hash password with the hash version
	// of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(password)); err != nil {
		// If the two passwords do not match, set
		// a message and reload the page.
		log.Println(err)

		return flash.WithError(c, fm).Redirect("/login")
	}

	if err = sessions.CreateUserSession(c, storedUser.ID); err != nil {
		log.Println(err)
		fm["errorMessage"] = "Something has gone wrong: unable to log in"

		return flash.WithError(c, fm).Redirect("/login")
	}

	fm = fiber.Map{
		"success":        true,
		"successMessage": "You have logged in successfully!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/profile")
}

// Handle SignUp & start Session
func HandleRegisterUser(c *fiber.Ctx) error {
	email := strings.Trim(c.FormValue("email"), " ")
	password := strings.Trim(c.FormValue("password"), " ")
	username := strings.Trim(c.FormValue("username"), " ")

	newUser := new(models.User)
	newUser.Email = email
	newUser.Password = password
	newUser.Username = username

	fm := fiber.Map{
		"error":        true,
		"errorMessage": "",
	}

	u, err := newUser.CreateUser()
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			fm["errorMessage"] = "The email is already in use"
		} else {
			fm["errorMessage"] = "Something went wrong"
		}

		return flash.WithError(c, fm).Redirect("/register")
	}

	if err = sessions.CreateUserSession(c, u.ID); err != nil {
		log.Println(err)
		fm["errorMessage"] = "Something has gone wrong: unable to sign up"

		return flash.WithError(c, fm).Redirect("/register")
	}

	fm = fiber.Map{
		"success":        true,
		"successMessage": "It has been successfully registered!!",
	}

	return flash.WithSuccess(c, fm).Redirect("/profile")
}

// Handle Session Logout
func HandleSessionLogout(c *fiber.Ctx) error {
	logoutResults, err := sessions.RemoveUserSession(c)
	if err != nil {
		log.Println(err)
	}

	if logoutResults {
		// Forced redirection from the client. Redirect to Login Page
		c.Set("HX-Location", "/login")

		fm := fiber.Map{
			"success":        true,
			"successMessage": "You have successfully logged out!!",
		}

		flash.WithSuccess(c, fm)
	} else {
		// Forced redirection from the client. Redirect to the same page
		c.Set("HX-Location", "/profile")
	}

	return nil
}
