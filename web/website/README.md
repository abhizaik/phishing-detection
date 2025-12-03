# SafeSurf Web Application

A SvelteKit web application that serves as the frontend for the SafeSurf phishing detection tool.

## Features

- **Clean, Minimal UI**: Modern interface with URL input and submit functionality
- **Comprehensive Security Checks**: Integrates with multiple backend API endpoints
- **Real-time Feedback**: Loading states and error handling
- **Organized Results**: Results grouped into logical sections (Reputation, Security, Structure, Network)
- **TypeScript Support**: Full type safety throughout the application

## API Endpoints

The application calls the following backend endpoints (assumed to be running on `http://localhost:8080`):

### Reputation Checks

- `GET /api/v1/rank` - Global domain rank
- `GET /api/v1/trusted-tld` - Trusted TLD check
- `GET /api/v1/risky-tld` - Risky TLD check
- `GET /api/v1/url-shortener` - URL shortener detection

### Security Checks

- `GET /api/v1/hsts` - HSTS header check
- `GET /api/v1/punycode` - Punycode detection
- `GET /api/v1/status-code` - HTTP status code

### Structure Checks

- `GET /api/v1/length` - URL length validation
- `GET /api/v1/depth` - URL depth validation
- `GET /api/v1/redirects` - Redirect chain analysis

### Network Checks

- `GET /api/v1/ip/check` - IP-based URL detection
- `GET /api/v1/ip/resolve` - IP resolution
- `GET /api/v1/whois` - WHOIS and domain age information

## Project Structure

```
src/
├── lib/
│   ├── components/
│   │   ├── LoadingSpinner.svelte
│   │   ├── ErrorMessage.svelte
│   │   └── ResultSection.svelte
│   ├── api.ts          # API service layer
│   ├── types.ts        # TypeScript interfaces
│   └── utils.ts        # Utility functions
├── routes/
│   ├── +layout.svelte  # Layout with CSS import
│   └── +page.svelte    # Main application page
├── app.css             # TailwindCSS imports
└── app.html            # HTML template
```

## Getting Started

### Prerequisites

- Node.js 18+
- Backend API running on `http://localhost:8080`

### Installation

1. Install dependencies:

   ```bash
   npm install
   ```

2. Start the development server:

   ```bash
   npm run dev
   ```

3. Open your browser to `http://localhost:5173` (or the port shown in the terminal)

### Building for Production

```bash
npm run build
```

## Usage

1. Enter a URL in the input field
2. Click "Check URL" to start the security analysis
3. View results organized by category:
   - **Reputation**: Domain rank, TLD trustworthiness, URL shortener detection
   - **Security**: HSTS, punycode usage, HTTP status codes
   - **Structure**: URL length, depth, redirect analysis
   - **Network**: IP-based checks, WHOIS information

## Features

### URL Validation

- Client-side URL validation before API calls
- Automatic protocol addition (https://) if missing
- Clear error messages for invalid URLs

### Loading States

- Individual loading spinners for each check
- Progress indication during API calls
- Graceful error handling for failed requests

### Result Display

- Color-coded results (green for safe, red for suspicious, yellow for warnings)
- Detailed WHOIS information including domain age
- Error messages for failed API calls

## Technologies Used

- **SvelteKit**: Full-stack web framework
- **TypeScript**: Type safety
- **TailwindCSS**: Styling
- **Vite**: Build tool and dev server

## API Configuration

The backend API URL is configured in `src/lib/api.ts`. To change the API endpoint, modify the `PUBLIC_BASE_URL` constant:

```typescript
const PUBLIC_BASE_URL = 'http://localhost:8080/api/v1';
```

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run check` - Type checking
- `npm run check:watch` - Watch mode type checking

### Adding New API Endpoints

1. Add the response type to `src/lib/types.ts`
2. Add the API method to `src/lib/api.ts`
3. Update the main page to include the new check
4. Add formatting logic to `ResultSection.svelte`

## Contributing

1. Follow the existing code structure
2. Use TypeScript for type safety
3. Maintain consistent styling with TailwindCSS
4. Test with various URL types and error conditions
