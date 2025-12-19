package transport

import (
	"net/http"
	"strconv"
	"time"

	"fandom/notifications/internal/models"
	"fandom/notifications/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	logService *service.LogService
}

func NewDashboardHandler(logService *service.LogService) *DashboardHandler {
	return &DashboardHandler{logService: logService}
}

// GetLogs godoc
// @Summary      Get request logs
// @Description  Retrieve paginated request logs with optional filters
// @Tags         dashboard
// @Produce      json
// @Param        limit       query     int     false  "Limit (max 1000)"
// @Param        offset      query     int     false  "Offset"
// @Param        method      query     string  false  "Filter by HTTP method"
// @Param        status_code query     int     false  "Filter by status code"
// @Param        path        query     string  false  "Filter by path (partial match)"
// @Param        start_date  query     string  false  "Start date (RFC3339)"
// @Param        end_date    query     string  false  "End date (RFC3339)"
// @Success      200         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]string
// @Router       /dashboard/logs [get]
func (h *DashboardHandler) GetLogs(c *gin.Context) {
	var params models.LogQueryParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	if params.Limit <= 0 {
		params.Limit = 100
	}

	if params.Page > 0 {
		params.Offset = (params.Page - 1) * params.Limit
	}

	logs, total, err := h.logService.GetLogs(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
		"page":   params.Page,
	})
}

// GetStats godoc
// @Summary      Get log statistics
// @Description  Get aggregated statistics about request logs
// @Tags         dashboard
// @Produce      json
// @Param        start_date  query     string  false  "Start date (RFC3339)"
// @Param        end_date    query     string  false  "End date (RFC3339)"
// @Success      200         {object}  models.LogStats
// @Failure      500         {object}  map[string]string
// @Router       /dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {
	var startDate, endDate time.Time
	var err error

	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
	}

	stats, err := h.logService.GetStats(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// DashboardPage godoc
// @Summary      Dashboard HTML page
// @Description  Serve the dashboard HTML page. Requires authentication via cookie or query parameter.
// @Tags         dashboard
// @Produce      html
// @Router       /dashboard [get]
func (h *DashboardHandler) DashboardPage(c *gin.Context) {
	// Get default stats for the page
	stats, _ := h.logService.GetStats(c.Request.Context(), time.Time{}, time.Time{})

	// Get recent logs
	params := models.LogQueryParams{
		Limit:  50,
		Offset: 0,
	}
	logs, total, _ := h.logService.GetLogs(c.Request.Context(), params)

	html := generateDashboardHTML(stats, logs, total)
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func generateDashboardHTML(stats *models.LogStats, logs []models.RequestLog, total int64) string {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Request Logs Dashboard</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: #f5f5f5;
            padding: 20px;
        }
        .container { max-width: 1400px; margin: 0 auto; }
        h1 { color: #333; margin-bottom: 30px; }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .stat-card h3 {
            color: #666;
            font-size: 14px;
            margin-bottom: 10px;
            text-transform: uppercase;
        }
        .stat-card .value {
            font-size: 32px;
            font-weight: bold;
            color: #333;
        }
        .filters {
            background: white;
            padding: 20px;
            border-radius: 8px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .filters form {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
        }
        .filters input, .filters select {
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
        }
        .filters button {
            padding: 10px 20px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
        .filters button:hover { background: #0056b3; }
        .logs-table {
            background: white;
            border-radius: 8px;
            overflow: hidden;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        table {
            width: 100%;
            border-collapse: collapse;
        }
        th {
            background: #f8f9fa;
            padding: 12px;
            text-align: left;
            font-weight: 600;
            color: #333;
            border-bottom: 2px solid #dee2e6;
        }
        td {
            padding: 12px;
            border-bottom: 1px solid #dee2e6;
        }
        tr:hover { background: #f8f9fa; }
        .status-code {
            padding: 4px 8px;
            border-radius: 4px;
            font-weight: 600;
            font-size: 12px;
        }
        .status-2xx { background: #d4edda; color: #155724; }
        .status-3xx { background: #fff3cd; color: #856404; }
        .status-4xx { background: #f8d7da; color: #721c24; }
        .status-5xx { background: #f5c6cb; color: #721c24; }
        .response-time {
            font-family: monospace;
        }
        .pagination {
            padding: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .refresh-btn {
            padding: 8px 16px;
            background: #28a745;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Request Logs Dashboard</h1>
        
        <div class="stats-grid">
            <div class="stat-card">
                <h3>Total Requests</h3>
                <div class="value">` + formatNumber(stats.TotalRequests) + `</div>
            </div>
            <div class="stat-card">
                <h3>Avg Response Time</h3>
                <div class="value">` + formatFloat(stats.AverageResponseTime) + `ms</div>
            </div>
        </div>

        <div class="filters">
            <form id="filterForm" onsubmit="event.preventDefault(); loadLogs();">
                <input type="text" name="method" placeholder="Method (GET, POST, etc.)" id="method">
                <input type="number" name="status_code" placeholder="Status Code" id="status_code">
                <input type="text" name="path" placeholder="Path" id="path">
                <input type="datetime-local" name="start_date" id="start_date">
                <input type="datetime-local" name="end_date" id="end_date">
                <button type="submit">Filter</button>
                <button type="button" onclick="clearFilters()">Clear</button>
                <button type="button" onclick="window.location.href='/auth'" style="background: #dc3545;">Logout</button>
            </form>
        </div>

        <div class="logs-table">
            <table>
                <thead>
                    <tr>
                        <th>Time</th>
                        <th>Method</th>
                        <th>Path</th>
                        <th>Status</th>
                        <th>Response Time</th>
                        <th>IP Address</th>
                    </tr>
                </thead>
                <tbody id="logsTableBody">
                    ` + generateLogsTableRows(logs) + `
                </tbody>
            </table>
            <div class="pagination">
                <div>Showing ` + formatNumber(int64(len(logs))) + ` of ` + formatNumber(total) + ` logs</div>
                <div>
                    <button onclick="loadPreviousPage()">Previous Page</button>
                    <button onclick="loadNextPage()">Next Page</button>
                </div>
                <button class="refresh-btn" onclick="loadLogs()">🔄 Refresh</button>
            </div>
        </div>
    </div>

    <script>
        // Check authentication on page load
        window.onload = function() {
            loadLogs();
        };

        function loadPreviousPage() {
            const params = new URLSearchParams(location.search);
            const currentPage = Number(params.get('page') || 1);
            if (currentPage > 1) {
                params.set('page', currentPage - 1);
                window.location.search = params.toString();
                loadLogs();
            }
        }

        function loadNextPage() {
            const params = new URLSearchParams(location.search);
            const currentPage = Number(params.get('page') || 1);
            params.set('page', currentPage + 1);
            window.location.search = params.toString();
            loadLogs();
        }

        function loadLogs() {
            const params = new URLSearchParams(location.search);
            const method = document.getElementById('method').value;
            const statusCode = document.getElementById('status_code').value;
            const path = document.getElementById('path').value;
            const startDate = document.getElementById('start_date').value;
            const endDate = document.getElementById('end_date').value;

            if (method) params.append('method', method);
            if (statusCode) params.append('status_code', statusCode);
            if (path) params.append('path', path);
            if (startDate) params.append('start_date', new Date(startDate).toISOString());
            if (endDate) params.append('end_date', new Date(endDate).toISOString());

            fetch('/dashboard/logs?' + params.toString(), {
                credentials: 'include'
            })
                .then(r => {
                    if (r.status === 403 || r.status === 401) {
                        window.location.href = '/auth';
                        return null;
                    }
                    if (!r.ok) {
                        throw new Error('Failed to load logs');
                    }
                    return r.json();
                })
                .then(data => {
                    if (!data) return;
                    const tbody = document.getElementById('logsTableBody');
                    if (data.logs && data.logs.length > 0) {
                        tbody.innerHTML = data.logs.map(log => {
                            const statusClass = getStatusClass(log.status_code);
                            return '<tr>' +
                                '<td>' + new Date(log.created_at).toLocaleString() + '</td>' +
                                '<td><strong>' + log.method + '</strong></td>' +
                                '<td>' + log.path + '</td>' +
                                '<td><span class="status-code status-' + statusClass + '">' + log.status_code + '</span></td>' +
                                '<td class="response-time">' + log.response_time_ms + 'ms</td>' +
                                '<td>' + log.ip_address + '</td>' +
                                '</tr>';
                        }).join('');
                    } else {
                        tbody.innerHTML = '<tr><td colspan="6" style="text-align: center; padding: 20px;">No logs found</td></tr>';
                    }
                })
                .catch(err => {
                    console.error('Error loading logs:', err);
                    const tbody = document.getElementById('logsTableBody');
                    tbody.innerHTML = '<tr><td colspan="6" style="text-align: center; padding: 20px; color: #dc3545;">Error loading logs. Please try again.</td></tr>';
                });
        }

        function clearFilters() {
            document.getElementById('filterForm').reset();
            loadLogs();
        }

        function getStatusClass(code) {
            if (code >= 200 && code < 300) return '2xx';
            if (code >= 300 && code < 400) return '3xx';
            if (code >= 400 && code < 500) return '4xx';
            return '5xx';
        }

        // Auto-refresh every 30 seconds
        setInterval(loadLogs, 30000);
    </script>
</body>
</html>`
	return html
}

func generateLogsTableRows(logs []models.RequestLog) string {
	rows := ""
	for _, log := range logs {
		statusClass := "2xx"
		if log.StatusCode >= 300 && log.StatusCode < 400 {
			statusClass = "3xx"
		} else if log.StatusCode >= 400 && log.StatusCode < 500 {
			statusClass = "4xx"
		} else if log.StatusCode >= 500 {
			statusClass = "5xx"
		}

		rows += `<tr>
			<td>` + log.CreatedAt.Format("2006-01-02 15:04:05") + `</td>
			<td><strong>` + log.Method + `</strong></td>
			<td>` + log.Path + `</td>
			<td><span class="status-code status-` + statusClass + `">` + strconv.Itoa(log.StatusCode) + `</span></td>
			<td class="response-time">` + strconv.FormatInt(log.ResponseTime, 10) + `ms</td>
			<td>` + log.IPAddress + `</td>
		</tr>`
	}
	return rows
}

func formatNumber(n int64) string {
	return strconv.FormatInt(n, 10)
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
