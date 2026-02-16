# Vocabulary Learning App - Frontend

## Project Structure

```
frontend/
├── public/
├── src/
│   ├── components/
│   │   ├── WordCard/           # Word card component with flip/reveal effect
│   │   ├── QualityButtons/     # Quality rating buttons (0-5)
│   │   ├── ProgressBar/        # Progress bar component
│   │   └── UploadStatus/       # Upload progress component
│   ├── pages/
│   │   ├── DictionaryUpload/   # Dictionary upload page
│   │   ├── Learning/           # Learning page (core)
│   │   └── Dashboard/          # Dictionary list dashboard
│   ├── services/
│   │   ├── request.ts          # Axios instance configuration
│   │   ├── dictionary.ts       # Dictionary API
│   │   └── learning.ts         # Learning API
│   ├── hooks/
│   │   └── useLearning.ts      # Learning logic hook
│   ├── styles/
│   │   └── global.css          # Global styles
│   ├── App.tsx                 # Main app component
│   └── main.tsx                # Entry point
├── .env                        # Environment variables
├── package.json
└── vite.config.js
```

## Features

### 1. Dictionary Management
- File upload with drag & drop support
- Real-time upload progress tracking
- Dictionary list with word count display

### 2. Learning Interface (Core)
- Word card with reveal effect
- Quality rating buttons (0-5 scale)
- Progress tracking
- Completion celebration

### 3. Global State
- User login state (mock)
- Learning task queue caching

## API Integration

### Environment Variables
```bash
VITE_API_BASE_URL=http://localhost:8000/api/v1
```

### Core Endpoints
1. `POST /api/v1/dictionaries/upload` - Upload dictionary
2. `GET /api/v1/dictionaries/upload/status/:task_id` - Query upload progress
3. `GET /api/v1/learning/today-tasks` - Get today's learning tasks
4. `POST /api/v1/learning/submit` - Submit learning result

## Development

```bash
# Install dependencies
pnpm install

# Run development server
pnpm dev

# Build for production
pnpm build
```

## Tech Stack
- React 18
- TypeScript
- Vite
- Axios
- CSS Modules
