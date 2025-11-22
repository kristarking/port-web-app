package main

import (
	"fmt"
	"net/http"
)

// Project represents a portfolio project
type Project struct {
	Title       string
	Description string
	Technology  string
}

// Global projects slice to store portfolio projects
var projects = []Project{
	{
		Title:       "Container Orchestration Platform",
		Description: "Kubernetes-based container orchestration solution",
		Technology:  "Docker, Kubernetes, Helm",
	},
	{
		Title:       "CI/CD Pipeline",
		Description: "Automated deployment pipeline",
		Technology:  "Jenkins, GitLab CI, ArgoCD",
	},
}

// Handler for the home page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>DevOps Portfolio</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    line-height: 1.6;
                    margin: 0;
                    padding: 20px;
                    background-color: #f4f4f4;
                }
                header {
                    background-color: #333;
                    color: white;
                    text-align: center;
                    padding: 1rem;
                    margin-bottom: 20px;
                }
                nav {
                    margin-bottom: 20px;
                }
                nav a {
                    margin-right: 15px;
                    text-decoration: none;
                    color: #333;
                }
                .container {
                    max-width: 800px;
                    margin: auto;
                }
            </style>
        </head>
        <body>
            <header>
                <h1>My DevOps Portfolio</h1>
            </header>
            <div class="container">
                <nav>
                    <a href="/">Home</a>
                    <a href="/projects">Projects</a>
                    <a href="/contact">Contact</a>
                </nav>
                <h2>Welcome to my DevOps Portfolio</h2>
                <p>I am a DevOps engineer specializing in cloud infrastructure, automation, and continuous integration/deployment.</p>
                <p>Browse through my projects to see examples of my work in action.</p>
            </div>
        </body>
    </html>`

	fmt.Fprintf(w, html)
}

// Handler for the projects page
func projectsHandler(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>Projects - DevOps Portfolio</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    line-height: 1.6;
                    margin: 0;
                    padding: 20px;
                    background-color: #f4f4f4;
                }
                header {
                    background-color: #333;
                    color: white;
                    text-align: center;
                    padding: 1rem;
                    margin-bottom: 20px;
                }
                nav {
                    margin-bottom: 20px;
                }
                nav a {
                    margin-right: 15px;
                    text-decoration: none;
                    color: #333;
                }
                .container {
                    max-width: 800px;
                    margin: auto;
                }
                .project {
                    background: white;
                    padding: 20px;
                    margin-bottom: 20px;
                    border-radius: 5px;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
                }
            </style>
        </head>
        <body>
            <header>
                <h1>My Projects</h1>
            </header>
            <div class="container">
                <nav>
                    <a href="/">Home</a>
                    <a href="/projects">Projects</a>
                    <a href="/contact">Contact</a>
                </nav>`

	// Add projects dynamically
	for _, project := range projects {
		projectHTML := fmt.Sprintf(`
            <div class="project">
                <h3>%s</h3>
                <p>%s</p>
                <p><strong>Technologies:</strong> %s</p>
            </div>`,
			project.Title, project.Description, project.Technology)
		html += projectHTML
	}

	html += `
            </div>
        </body>
    </html>`

	fmt.Fprintf(w, html)
}

// Handler for the contact page
func contactHandler(w http.ResponseWriter, r *http.Request) {
	html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>Contact - DevOps Portfolio</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    line-height: 1.6;
                    margin: 0;
                    padding: 20px;
                    background-color: #f4f4f4;
                }
                header {
                    background-color: #333;
                    color: white;
                    text-align: center;
                    padding: 1rem;
                    margin-bottom: 20px;
                }
                nav {
                    margin-bottom: 20px;
                }
                nav a {
                    margin-right: 15px;
                    text-decoration: none;
                    color: #333;
                }
                .container {
                    max-width: 800px;
                    margin: auto;
                }
                .contact-info {
                    background: white;
                    padding: 20px;
                    border-radius: 5px;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
                }
            </style>
        </head>
        <body>
            <header>
                <h1>Contact Me</h1>
            </header>
            <div class="container">
                <nav>
                    <a href="/">Home</a>
                    <a href="/projects">Projects</a>
                    <a href="/contact">Contact</a>
                </nav>
                <div class="contact-info">
                    <h2>Get in Touch</h2>
                    <p>Email: your.email@example.com</p>
                    <p>LinkedIn: linkedin.com/in/yourprofile</p>
                    <p>GitHub: github.com/yourusername</p>
                </div>
            </div>
        </body>
    </html>`

	fmt.Fprintf(w, html)
}

func main() {
	// Register handlers
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/contact", contactHandler)

	// Start the server
	fmt.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
