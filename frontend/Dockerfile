# Build stage
FROM node:20-alpine

WORKDIR /app

# Copy package files
COPY package.json package-lock.json* ./

# Install dependencies
RUN npm install

# Copy source files
COPY . .

# Build the app
RUN npm run build

# Expose port (Railway will use PORT env var)
EXPOSE 5173

# Start the preview server
CMD ["npm", "start"]
