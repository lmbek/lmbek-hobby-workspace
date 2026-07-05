package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type ServiceStatus struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Status      string `json:"status"` // "online", "offline"
	InternalURL string `json:"-"`
}

var (
	services = []ServiceStatus{
		{
			Name:        "Service 1",
			Description: "Placeholder Service 1",
			URL:         "http://service1.localhost",
			InternalURL: "http://placeholder1-service:8081/health",
		},
		{
			Name:        "Service 2",
			Description: "Placeholder Service 2",
			URL:         "http://service2.localhost",
			InternalURL: "http://placeholder2-service:8082/health",
		},
		{
			Name:        "Traefik Dashboard",
			Description: "Infrastructure Proxy",
			URL:         "http://localhost:8080",
			InternalURL: "http://traefik:8080/api/overview",
		},
	}
	mu sync.RWMutex
)

func main() {
	// Initialize status to offline
	for i := range services {
		services[i].Status = "offline"
	}

	go startHealthChecks()

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/error", errorHandler)

	fmt.Println("Proxy Dashboard starting on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}

func startHealthChecks() {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	for {
		for i := range services {
			resp, err := client.Get(services[i].InternalURL)
			mu.Lock()
			if err == nil {
				if resp.StatusCode >= 200 && resp.StatusCode < 300 {
					services[i].Status = "online"
				} else {
					services[i].Status = "offline"
				}
				resp.Body.Close()
			} else {
				services[i].Status = "offline"
			}
			mu.Unlock()
		}
		time.Sleep(5 * time.Second)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("service")
	if serviceName == "" {
		serviceName = "Target Service"
	}
	w.WriteHeader(http.StatusBadGateway)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>502 Bad Gateway</title>
    <style>
        body { font-family: sans-serif; text-align: center; padding: 50px; background: #f8f9fa; color: #333; }
        h1 { color: #dc3545; font-size: 5rem; margin: 0; }
        h2 { font-size: 2rem; margin-top: 0; }
        p { color: #6c757d; font-size: 1.2rem; }
        .container { max-width: 600px; margin: auto; background: white; padding: 40px; border-radius: 10px; box-shadow: 0 4px 15px rgba(0,0,0,0.1); }
        .btn { display: inline-block; margin-top: 20px; padding: 12px 24px; background: #007bff; color: white; text-decoration: none; border-radius: 5px; font-weight: bold; }
        .btn:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <h1>502</h1>
        <h2>Bad Gateway</h2>
        <p>Oops! The service seems to be down or unreachable.</p>
        <p>Make sure you have started the system with <code>make up</code> and wait a few seconds for services to boot.</p>
        <a href="http://localhost" class="btn">Back to Dashboard</a>
    </div>
</body>
</html>
`)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<!DOCTYPE html>
<html>
<head>
    <title>Workspace Dashboard</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; margin: 0; background: #f0f2f5; color: #1c1e21; }
        header { background: #fff; padding: 1rem 2rem; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 2rem; }
        h1 { margin: 0; font-size: 1.5rem; color: #007bff; }
        .container { padding: 0 2rem; max-width: 1200px; margin: auto; }
        .card-container { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 1.5rem; }
        .card { 
            background: white; 
            padding: 1.5rem; 
            border-radius: 12px; 
            box-shadow: 0 1px 3px rgba(0,0,0,0.1); 
            text-align: left;
            transition: all 0.2s ease-in-out;
            border: 1px solid #e1e4e8;
            display: flex;
            flex-direction: column;
        }
        .card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0,0,0,0.15); }
        .card h2 { margin: 0 0 0.5rem 0; font-size: 1.25rem; }
        .card p { margin: 0; color: #65676b; flex-grow: 1; }
        .card a { 
            display: inline-block; 
            margin-top: 1.5rem; 
            padding: 0.75rem 1rem; 
            background: #007bff; 
            color: white; 
            text-decoration: none; 
            border-radius: 6px; 
            text-align: center;
            font-weight: 600;
        }
        .card a:hover { background: #0056b3; }
        .card a.offline { background: #ebedf0; color: #bcc0c4; pointer-events: none; }
        
        .status-badge {
            display: inline-flex;
            align-items: center;
            font-size: 0.875rem;
            font-weight: 600;
            margin-top: 1rem;
        }
        .dot { height: 8px; width: 8px; border-radius: 50%; display: inline-block; margin-right: 6px; }
        
        .online .dot { background-color: #31a24c; box-shadow: 0 0 0 2px rgba(49, 162, 76, 0.2); }
        .online { color: #31a24c; }
        
        .offline .dot { background-color: #dc3545; }
        .offline { color: #dc3545; }

        .loader {
            border: 3px solid #f3f3f3;
            border-top: 3px solid #007bff;
            border-radius: 50%;
            width: 20px;
            height: 20px;
            animation: spin 1s linear infinite;
            display: inline-block;
            vertical-align: middle;
            margin-left: 10px;
        }
        @keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
    </style>
</head>
<body>
    <header>
        <div style="display: flex; justify-content: space-between; align-items: center; max-width: 1200px; margin: auto;">
            <h1>Hobby Workspace <span id="loading" class="loader" style="display:none;"></span></h1>
            <div id="last-update" style="font-size: 0.8rem; color: #65676b;"></div>
        </div>
    </header>
    
    <div class="container">
        <div id="cards" class="card-container">
            <!-- Cards will be injected here -->
        </div>
    </div>

    <script>
        async function updateStatus() {
            const loading = document.getElementById('loading');
            loading.style.display = 'inline-block';
            try {
                const response = await fetch('/api/status');
                const services = await response.json();
                const container = document.getElementById('cards');
                container.innerHTML = '';

                services.forEach(s => {
                    const card = document.createElement('div');
                    card.className = 'card';
                    
                    const isOnline = s.status === 'online';
                    const statusClass = isOnline ? 'online' : 'offline';
                    const statusText = isOnline ? 'Online' : 'Offline';
                    
                    card.innerHTML = '<h2>' + s.name + '</h2>' +
                        '<p>' + s.description + '</p>' +
                        '<div class="status-badge ' + statusClass + '">' +
                            '<span class="dot"></span> ' + statusText +
                        '</div>' +
                        '<a href="' + s.url + '" target="_blank" class="' + (isOnline ? '' : 'offline') + '">' +
                            (isOnline ? 'Open Service' : 'Service Down') +
                        '</a>';
                    container.appendChild(card);
                });
                document.getElementById('last-update').innerText = 'Last updated: ' + new Date().toLocaleString();
            } catch (e) {
                console.error("Failed to fetch status", e);
            } finally {
                setTimeout(() => { loading.style.display = 'none'; }, 500);
            }
        }

        setInterval(updateStatus, 5000);
        updateStatus();
    </script>
</body>
</html>
`)
}
