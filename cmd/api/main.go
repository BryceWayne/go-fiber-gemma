package main

import (
	"log"
	"net"

	"github.com/BryceWayne/go-fiber-gemma/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/hashicorp/mdns"
)

func main() {
	// 1. Setup mDNS (gemma.local)

	// Fetch local IPs to actually broadcast gemma.local A records
	var ips []net.IP
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP)
					}
				}
			}
		}
	}

	info := []string{"Gemma 4 Local API"}
	// Create service using gemma.local. as host so clients resolve it directly via mdns
	service, _ := mdns.NewMDNSService("gemma", "_http._tcp", "local.", "gemma.local.", 3000, ips, info)
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()

	// Initialize standard Go html template engine
	engine := html.New("./web/templates", ".html")

	// 2. Initialize Fiber with a 50MB body limit
	app := fiber.New(fiber.Config{
		BodyLimit: 50 * 1024 * 1024,
		Views:     engine, // set template engine
	})

	// Middleware
	app.Use(logger.New()) // Log every request to your terminal
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Serve static files (CSS, JS, Images)
	app.Static("/static", "./web/static")

	// 3. Web UI Route
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Gemma Chat UI",
		})
	})

	// HTMX Chat Route
	app.Post("/chat", handlers.HandleHTMXChat)
	app.Get("/chat", func(c *fiber.Ctx) error {
		return c.Redirect("/")
	})

	// API Routes
	app.Get("/api/status", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "online",
			"host":   "gemma.local",
			"gpu":    "RTX 4070",
		})
	})

	app.Post("/api/ask", handlers.HandleAsk)
	app.Get("/api/ask", func(c *fiber.Ctx) error {
		return c.Status(405).JSON(fiber.Map{
			"error": "Method Not Allowed. This endpoint requires a POST request with a JSON payload: {'prompt': '...'} ",
		})
	})

	log.Println("🚀 Gemma 4 API & Web UI serving at http://gemma.local:3000")
	log.Fatal(app.Listen(":3000"))
}
