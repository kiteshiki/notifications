package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthPage godoc
// @Summary      Authentication page
// @Description  HTML page for entering API key
// @Tags         auth
// @Produce      html
// @Router       /auth [get]
func AuthPage(c *gin.Context) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Key Authentication</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .auth-container {
            background: white;
            padding: 40px;
            border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2);
            max-width: 450px;
            width: 100%;
        }
        h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 28px;
        }
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            color: #333;
            margin-bottom: 8px;
            font-weight: 500;
            font-size: 14px;
        }
        input[type="text"], input[type="password"] {
            width: 100%;
            padding: 12px;
            border: 2px solid #e0e0e0;
            border-radius: 6px;
            font-size: 14px;
            transition: border-color 0.3s;
        }
        input[type="text"]:focus, input[type="password"]:focus {
            outline: none;
            border-color: #667eea;
        }
        .button-group {
            display: flex;
            gap: 10px;
            margin-top: 30px;
        }
        button {
            flex: 1;
            padding: 12px;
            border: none;
            border-radius: 6px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s;
        }
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        .btn-secondary {
            background: #f5f5f5;
            color: #666;
        }
        .btn-secondary:hover {
            background: #e0e0e0;
        }
        .message {
            padding: 12px;
            border-radius: 6px;
            margin-bottom: 20px;
            font-size: 14px;
            display: none;
        }
        .message.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        .message.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        .info {
            background: #e7f3ff;
            color: #004085;
            padding: 15px;
            border-radius: 6px;
            margin-top: 20px;
            font-size: 13px;
            line-height: 1.6;
        }
        .info strong {
            display: block;
            margin-bottom: 5px;
        }
    </style>
</head>
<body>
    <div class="auth-container">
        <h1>🔐 API Key Authentication</h1>
        <p class="subtitle">Enter your API key to access the dashboard</p>
        
        <div id="message" class="message"></div>
        
        <form id="authForm" onsubmit="event.preventDefault(); submitAuth();">
            <div class="form-group">
                <label for="apiKey">API Key</label>
                <input 
                    type="password" 
                    id="apiKey" 
                    name="apiKey" 
                    placeholder="Enter your API key"
                    required
                    autocomplete="off"
                >
            </div>
            <div class="button-group">
                <button type="submit" class="btn-primary">Authenticate</button>
                <button type="button" class="btn-secondary" onclick="clearAuth()">Clear</button>
            </div>
        </form>
        
        <div class="info">
            <strong>ℹ️ Information:</strong>
            Your API key will be stored in a secure cookie and used for all subsequent requests.
            You can clear it at any time using the Clear button.
        </div>
    </div>

    <script>
        // Check if already authenticated
        window.onload = function() {
            const cookies = document.cookie.split(';');
            const apiKeyCookie = cookies.find(c => c.trim().startsWith('api_key='));
            if (apiKeyCookie) {
                showMessage('You are already authenticated. Redirecting to dashboard...', 'success');
                setTimeout(() => {
                    window.location.href = '/dashboard';
                }, 1500);
            }
        };

        function submitAuth() {
            const apiKey = document.getElementById('apiKey').value.trim();
            if (!apiKey) {
                showMessage('Please enter an API key', 'error');
                return;
            }

            fetch('/auth/set', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ api_key: apiKey })
            })
            .then(response => response.json())
            .then(data => {
                if (data.success) {
                    showMessage('Authentication successful! Redirecting...', 'success');
                    setTimeout(() => {
                        window.location.href = '/dashboard';
                    }, 1000);
                } else {
                    showMessage(data.error || 'Authentication failed', 'error');
                }
            })
            .catch(error => {
                showMessage('An error occurred. Please try again.', 'error');
                console.error('Error:', error);
            });
        }

        function clearAuth() {
            fetch('/auth/clear', {
                method: 'POST'
            })
            .then(() => {
                document.getElementById('apiKey').value = '';
                showMessage('Authentication cleared', 'success');
            })
            .catch(error => {
                console.error('Error:', error);
            });
        }

        function showMessage(text, type) {
            const messageDiv = document.getElementById('message');
            messageDiv.textContent = text;
            messageDiv.className = 'message ' + type;
            messageDiv.style.display = 'block';
            
            if (type === 'success') {
                setTimeout(() => {
                    messageDiv.style.display = 'none';
                }, 3000);
            }
        }
    </script>
</body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// SetAuthCookie godoc
// @Summary      Set API key cookie
// @Description  Stores the provided API key in a secure cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      map[string]string  true  "API Key"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Router       /auth/set [post]
func SetAuthCookie(c *gin.Context) {
	var req struct {
		APIKey string `json:"api_key" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "API key is required",
		})
		return
	}

	// Set secure HTTP-only cookie (valid for 7 days)
	c.SetCookie(
		"api_key",
		req.APIKey,
		7*24*60*60, // 7 days in seconds
		"/",
		"",
		false, // Set to true in production with HTTPS
		true,  // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "API key stored successfully",
	})
}

// ClearAuthCookie godoc
// @Summary      Clear API key cookie
// @Description  Removes the API key cookie
// @Tags         auth
// @Produce      json
// @Success      200     {object}  map[string]interface{}
// @Router       /auth/clear [post]
func ClearAuthCookie(c *gin.Context) {
	// Clear the cookie by setting it to expire in the past
	c.SetCookie(
		"api_key",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "API key cleared successfully",
	})
}

