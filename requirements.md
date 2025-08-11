# Go Web Scraper Server - Requirements

## Project Overview
A Go-based generic content scraper and LLM parser service. Originally designed for job description extraction, the architecture supports scraping any web content and parsing it into structured data using client-provided schemas and LLM prompts. The service integrates with Next.js applications (particularly a Job Application Tracker) hosted on Vercel through a polling-based API.

## Functional Requirements

### Core Scraping Functionality
1. **Generic Parse Endpoint**
   - Accept POST requests with any web URL for content extraction
   - Support client-provided JSON schemas for structured data output
   - Return job IDs for async processing and polling
   - Handle authentication/authorization for API access

2. **Dual Scraping Strategy**
   - **Primary Method**: Simple HTTP fetch + manual DOM traversal
     - Fast, lightweight approach for basic HTML pages
     - Parse HTML using Go's `golang.org/x/net/html` package
     - Extract visible text while filtering out scripts, styles, and other non-content elements
   
   - **Fallback Method**: Headless Chrome browser scraping
     - Handle JavaScript-heavy sites that require rendering
     - Use when primary method fails or returns insufficient content
     - Implement using Chrome DevTools Protocol or similar

3. **LLM-Based Content Processing**
   - Extract meaningful content from any webpage type
   - Filter out navigation, headers, footers, and irrelevant page elements
   - Use OpenRouter API to parse content into structured data
   - Support client-provided JSON schemas and parsing prompts
   - Validate LLM output against provided schemas

## Non-Functional Requirements

### Performance
- Response time < 5 seconds for simple scraping
- Response time < 15 seconds for headless browser fallback
- Support concurrent requests
- Implement request timeout handling

### Reliability
- Graceful fallback from DOM parsing to headless browser
- Proper error handling and meaningful error messages
- Retry logic for transient failures
- Rate limiting to prevent abuse

### Security
- Input validation for URLs
- Protection against malicious URLs
- CORS configuration for integration with Next.js frontend
- UUID-based job access control (cryptographically secure UUIDs)
- 24-hour TTL for automatic data cleanup
- Future: API key authentication, client-side encryption, or JWT-based auth

### Scalability
- Stateless design for easy horizontal scaling
- Resource cleanup for headless browser instances
- Connection pooling for HTTP requests

## Integration Requirements

### Schema-Driven Client Integration
- Fetch JSON schemas and prompts from client-provided endpoints
- Support Zod-generated schemas from Next.js applications
- Polling-based result retrieval (no webhooks required)
- Generic enough to support multiple client applications and content types

### Next.js Job Application Tracker Integration
- RESTful API interface compatible with frontend JavaScript
- JSON request/response format
- CORS headers configured for Vercel domain
- Error responses in consistent format for frontend error handling
- Work within Vercel free tier constraints (no webhook support)

### Deployment
- Containerized deployment (Docker)
- Environment variable configuration
- Health check endpoint for load balancer/orchestration
- Logging for monitoring and debugging

## API Specification

### Parse Job Submission
```json
POST /parse
{
  "url": "https://example.com/content-page",
  "schema_endpoint": "https://client-app.vercel.app/api/schema/content-type",
  "client_id": "job-tracker-app",
  "metadata": {
    "user_id": "user_123",
    "content_type": "job_posting"
  },
  "options": {
    "timeout": 30,
    "force_headless": false
  }
}
```

### Immediate Response Format
```json
{
  "success": true,
  "job_id": "abc123-def456-ghi789",
  "status": "processing",
  "estimated_completion": "2025-08-11T10:35:00Z"
}
```

### Polling Endpoint
```json
GET /parse/{job_id}

// Processing Response
{
  "job_id": "abc123-def456-ghi789",
  "status": "processing",
  "progress": {
    "step": "scraping", // "scraping" | "parsing" | "completed"
    "message": "Extracting content from webpage..."
  }
}

// Completed Response
{
  "job_id": "abc123-def456-ghi789",
  "status": "completed",
  "result": {
    "structured_data": { /* parsed according to client schema */ },
    "raw_content": "extracted webpage text",
    "method": "dom",
    "processed_at": "2025-08-11T10:33:00Z"
  },
  "metadata": {
    "user_id": "user_123",
    "content_type": "job_posting"
  }
}
```

### Schema Endpoint Format (Client-Provided)
```json
GET https://client-app.vercel.app/api/schema/job

{
  "schema": {
    "type": "object",
    "properties": {
      "title": {"type": "string"},
      "company": {"type": "string"},
      "salary": {"type": "object"},
      "skills": {"type": "array"}
    }
  },
  "prompt": "Parse this job description into the following JSON structure..."
}
```

### Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "SCRAPE_FAILED",
    "message": "Failed to extract content from URL",
    "details": "Additional error context"
  }
}
```

## Technical Constraints

### Dependencies
- Go 1.21+
- Chrome/Chromium for headless browsing
- Standard library HTTP server or lightweight framework (Gin, Echo)
- OpenRouter API access for LLM processing
- Database for temporary job storage (PostgreSQL, SQLite)

### Resource Limits
- Memory limit for headless browser instances
- Concurrent request limits
- Request timeout limits

## Success Criteria
1. Successfully extracts and parses content from major sites (job boards, e-commerce, news, etc.)
2. < 10% failure rate on supported website types
3. Generic architecture supports multiple client applications and content types
4. Seamless integration with Next.js applications using Zod schemas
5. Handles Vercel timeout constraints through async processing
6. Server can handle expected load with proper job queuing
7. Deployment is reliable and maintainable
8. Schema evolution works smoothly without breaking changes

## Additional Requirements

### Job Management
- UUID-based job identification (UUID4 for cryptographic security)
- Persistent job storage with status tracking
- TTL-based cleanup of completed jobs (24-hour retention)
- Job progress reporting and error details
- Support for job cancellation and retry
- Security through obscurity: no authentication required for job access

### Multi-Client Support
- Client identification and isolation
- Rate limiting per client
- Generic job storage supporting different content types
- Configurable processing timeouts per client

### LLM Integration
- OpenRouter API integration with error handling
- Retry logic for LLM API failures
- Schema validation of LLM responses
- Support for different LLM models per client

### Monitoring & Observability
- Job processing metrics and timing
- Error rate tracking by client and content type
- Resource usage monitoring (memory, CPU)
- Health check endpoints for deployment platforms