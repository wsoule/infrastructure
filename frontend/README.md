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

## Railway Deployment

### Configuration

1. **Build Command**: `npm install && npm run build`
2. **Start Command**: `npm start`

### Environment Variables

Set the following environment variable in Railway:

- `VITE_API_GATEWAY_URL` - Your API Gateway URL (e.g., `https://your-api-gateway.railway.app`)

**Important**: Railway will automatically set the `PORT` variable. The start script will use it.

### Deployment Steps

1. Push your code to a Git repository
2. Create a new project in Railway
3. Connect your repository
4. Set the root directory to `frontend` (if not deploying from monorepo root)
5. Add the environment variable `VITE_API_GATEWAY_URL`
6. Deploy!

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
