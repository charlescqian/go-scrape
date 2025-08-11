# Architectural Context & Decisions

## Project Background

This Go web scraper was originally conceived as a microservice to solve a specific problem: the Next.js Job Application Tracker (deployed on Vercel) was unable to scrape JavaScript-heavy job posting sites using client-side libraries like Readability.js.

### Original Problem Statement
- Next.js app uses Readability.js to extract job descriptions
- JS-rendered sites fail with client-side scraping  
- LLM parsing via OpenRouter was working well in Next.js
- Vercel 30-second timeout kills long operations (scraping + LLM parsing)

## Key Architectural Decisions

### 1. Data Flow Architecture

**Rejected Options:**
- **Option A**: Scraper → Next.js → OpenRouter → Supabase
  - ❌ Vercel timeout on combined scraping + LLM (can take 35+ seconds)
- **Option B**: Scraper → Supabase → Next.js → OpenRouter  
  - ❌ Still has Vercel timeout risk on OpenRouter parsing

**Selected Option C**: Scraper → OpenRouter → Supabase → Next.js polls
- ✅ Eliminates ALL Vercel timeout issues
- ✅ Go server handles both scraping AND LLM parsing
- ✅ Next.js just polls for completed results
- ✅ Better retry logic and error handling in Go

### 2. Schema Management Strategy

**Challenge**: Avoid maintaining duplicate schemas in two repositories

**Solution**: Zod Schema Endpoints
- Next.js generates JSON schema from existing Zod definitions
- Exposes schema + prompt via API endpoint: `GET /api/schema/job`
- Go server fetches schema on startup/initialization
- Single source of truth maintained in Next.js with Zod

**Benefits**:
- ✅ Leverage existing Zod schemas and prompts
- ✅ No duplication of OpenRouter prompt engineering
- ✅ Schema evolution happens in one place
- ✅ Go server stays agnostic to specific data structures

### 3. Generic Architecture Evolution

**Key Insight**: The solution naturally becomes a generic content scraper + LLM parser service

**Generic API Design**:
```json
POST /parse
{
  "url": "https://any-website.com",
  "schema_endpoint": "https://client-app.com/api/schema/content-type",
  "client_id": "job-tracker-app",
  "metadata": {"user_id": "user_123"}
}
```

**Use Cases Enabled**:
- Job descriptions → structured job data
- Product pages → price, features, specs  
- Recipe sites → ingredients, instructions
- Real estate → property details
- News articles → structured article data

### 4. Data Routing Without Webhooks

**Constraint**: Vercel free plan doesn't support webhook endpoints

**Solution**: Generic Storage + Client Polling
- Go server maintains temporary results database
- Clients poll `GET /parse/{job_id}` for completion
- Clients handle their own data validation and storage
- TTL-based cleanup (24-hour job retention)

**Implementation**:
```
1. Client submits job → immediate job_id response
2. Go server processes async (scrape + LLM)  
3. Client polls every 2-3 seconds until complete
4. Client validates result and stores in own database
```

## Technical Architecture

### Go Server Responsibilities
1. **Content Scraping**: DOM parsing + headless browser fallback
2. **Schema Fetching**: Dynamic schema loading from client endpoints
3. **LLM Integration**: OpenRouter API calls with client-provided schemas
4. **Job Management**: Async processing with status tracking
5. **Result Storage**: Temporary results with TTL cleanup

### Next.js Client Responsibilities  
1. **Schema Definition**: Zod schemas exposed via API endpoints
2. **Prompt Engineering**: LLM prompts served alongside schemas
3. **Result Polling**: Monitor job status until completion
4. **Data Validation**: Validate LLM results against Zod schemas
5. **Database Storage**: Store validated results in Supabase

### Database Design

**Go Server Database (Temporary Storage)**:
```sql
parse_jobs (
  id UUID PRIMARY KEY,
  client_id VARCHAR,
  url TEXT,
  content_type VARCHAR,
  raw_content TEXT,
  structured_data JSONB,
  metadata JSONB,
  status VARCHAR, -- 'processing', 'completed', 'failed'
  error_message TEXT,
  created_at TIMESTAMP,
  expires_at TIMESTAMP -- TTL cleanup
)
```

**Client Database (Permanent Storage)**:
- Each client manages their own schema
- Go server results validated and transformed by client
- Enables client-specific business logic and relationships

## Benefits of This Architecture

### Scalability
- ✅ Stateless Go server design
- ✅ Horizontal scaling capability
- ✅ Multiple clients can use same service
- ✅ Language-agnostic client integration

### Maintainability  
- ✅ Single schema source of truth per client
- ✅ Separation of concerns (scraping vs. business logic)
- ✅ Independent deployment cycles
- ✅ Generic service reduces duplication

### Reliability
- ✅ No timeout issues on any platform
- ✅ Proper error handling and retry logic
- ✅ Graceful fallback strategies (DOM → headless browser)
- ✅ Built-in job persistence and recovery

### Flexibility
- ✅ Support for any content type and schema
- ✅ Client controls data validation and storage
- ✅ Easy to add new scraping strategies
- ✅ Platform-agnostic deployment options

## Implementation Phases

### Phase 1: Job-Specific Implementation
- Enhance current server for job description scraping
- Add OpenRouter integration with hardcoded job schema
- Implement basic error handling and status tracking

### Phase 2: Generic Architecture
- Add schema endpoint fetching capability
- Implement generic data storage and polling
- Create client polling utilities

### Phase 3: Production Hardening  
- Add authentication and rate limiting
- Implement proper logging and monitoring
- Add headless browser fallback strategy

### Phase 4: Multi-Client Support
- Deploy as shared service
- Documentation and client SDKs
- Support for multiple concurrent clients

## Security Model

### Current Approach: UUID-Based Security
- Each parse job assigned a cryptographically secure UUID (UUID4)
- 24-hour TTL for all job results
- Security through obscurity: UUIDs are practically impossible to guess
- No authentication required for initial implementation

**Benefits**:
- ✅ Simple to implement and maintain
- ✅ No user management or API key distribution needed  
- ✅ Automatic cleanup prevents data accumulation
- ✅ Statistically secure (2^122 possible UUIDs)

### Future Authentication Options

**Option 1: API Key Authentication**
- Client registration with API key generation
- Rate limiting and usage tracking per API key
- Revocable access control

**Option 2: Client-Side Encryption**
- Client provides public key with job submission
- Server encrypts results with client's public key
- Only client can decrypt results with private key
- Enables zero-knowledge storage

**Option 3: Signed JWTs**
- Client provides signed JWT with job submission
- Server validates JWT signature and extracts user context
- Results tied to JWT claims (user_id, permissions)

## Open Questions & Future Considerations

1. **Rate Limiting**: How to prevent abuse while allowing legitimate usage?
2. **Monitoring**: What metrics are important for service health?
3. **Scaling**: When to implement job queues vs. direct processing?
4. **Error Recovery**: How to handle partial failures and retries?
5. **Data Retention**: Should successful results be kept longer for debugging?

This architecture provides a solid foundation for both immediate needs and future expansion while working within the constraints of free-tier cloud platforms.