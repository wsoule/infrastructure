# API Gateway Frontend

A simple Vite-based frontend to demonstrate the API Gateway reverse proxy pattern.

## Local Development

1. Install dependencies:
```bash
npm install
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Update `.env` with your API Gateway URL (defaults to `http://localhost:8080`)

4. Run the development server:
```bash
npm run dev
```

The app will be available at `http://localhost:5173`

## Railway Deployment (Docker)

### Configuration

Railway will automatically detect and use the Dockerfile.

### Environment Variables

Set the following environment variable in Railway:

- `VITE_API_GATEWAY_URL` - Your API Gateway URL (e.g., `https://your-api-gateway.railway.app`)

**Important**: This environment variable must be set at **build time** since Vite bakes env vars into the bundle.

### Deployment Steps

1. Push your code to a Git repository
2. Create a new project in Railway
3. Connect your repository
4. Set the root directory to `frontend` (if deploying from monorepo root)
5. Add the environment variable `VITE_API_GATEWAY_URL`
6. Railway will automatically build and deploy using the Dockerfile
7. Railway will set the `PORT` environment variable automatically

## Docker Build & Run Locally

Build the Docker image:
```bash
docker build -t api-gateway-frontend .
```

Run the container:
```bash
docker run -p 3000:5173 api-gateway-frontend
```

Access at `http://localhost:3000`

## Build for Production

```bash
npm run build
```

This creates an optimized build in the `dist/` folder.

## Preview Production Build

```bash
npm run preview
```

This serves the production build locally at `http://localhost:5173`

## Available Scripts

- `npm run dev` - Start development server with hot reload
- `npm run build` - Build for production
- `npm run preview` - Preview production build locally
- `npm start` - Start production server (used by Railway)
