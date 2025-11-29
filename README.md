# Chak - Local AI Chat System

A comprehensive local AI chat application built from scratch to explore modern AI integration patterns, RAG (Retrieval-Augmented Generation), and production-grade backend architecture. The system runs entirely locally using Ollama for LLM services, with no cloud dependencies.

## Features

- **Multi-Profile System**: Switch between different contexts (coding, paperwork, general) with separate configurations and memory
- **Semantic Memory**: Long-term memory using vector embeddings for context-aware conversations
- **Document RAG**: Automatic document indexing with intelligent Markdown-aware chunking and semantic search
- **Web Search Integration**: Optional Brave Search API integration for current information
- **Hot Reload**: Switch profiles without server restart
- **Auto-Indexing**: File watching with hash-based change detection for efficient document updates

## Architecture

### Backend (Go)
- **Clean Architecture**: Interface-manager pattern with clear separation of concerns
- **Package Structure**:
  - `config`: Profile and configuration management
  - `handler`: HTTP request handlers (chat, profile switching)
  - `memory`: Vector-based semantic memory with filtering
  - `indexer`: Document scanning and indexing with watcher
  - `embedding`: Ollama embedding integration
  - `ollama`: LLM generation interface
  - `search`: Web search providers (Brave, DuckDuckGo)
  - `prompt`: Context-aware prompt building
  - `document`: Smart text chunking respecting document structure
  - `middleware`: CORS and logging

### Frontend (Vanilla JS)
- Simple web interface with markdown rendering
- Real-time chat with streaming support
- Profile selector and model switcher
- Toggle controls for web search and RAG

## Prerequisites

- **Go 1.25.4+**
- **Ollama** running locally
- **Embedding Model**: `all-minilm:33m` (or configure your own)
- **LLM Models**: Any Ollama-compatible models
- **Brave Search API Key** (optional, for web search)

## Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd chak-server
```

2. **Install dependencies**
```bash
cd server
go mod download
```

3. **Set up Ollama**
```bash
# Install Ollama (see ollama.ai)
# Pull required models
ollama pull all-minilm:33m
ollama pull llama2  # or your preferred model
```

4. **Configure environment**
```bash
# Create .env file in server directory
OLLAMA_HOST=localhost
BRAVE_API_KEY=your_api_key_here  # Optional
```

5. **Configure profiles**

Edit `server/config.json` to customize profiles:
```json
{
  "active_profile": "coding",
  "profiles": {
    "coding": {
      "name": "Coding Assistant",
      "description": "Programming help and code examples",
      "directories": ["./documents/coding"],
      "memory_file": "memory_coding.json",
      "index_file": "index_coding.json",
      "extensions": [".txt", ".md"],
      "max_file_size": 5242880
    }
  }
}
```

## Usage

1. **Start the server**
```bash
cd server
go run main.go
```

2. **Open the web interface** I use python in this case
```
cd web
python -m http.server 8080
```

3. **Chat with the AI**
- Select a model from the dropdown
- Toggle web search or RAG as needed
- Switch profiles on the fly
- Start chatting!

## Project Structure

```
.
├── server/
│   ├── main.go                 # Application entry point
│   ├── config.json             # Profile configurations
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   ├── handler/           # HTTP handlers
│   │   ├── memory/            # Semantic memory system
│   │   ├── indexer/           # Document indexing
│   │   ├── embedding/         # Vector embeddings
│   │   ├── ollama/            # LLM integration
│   │   ├── search/            # Web search providers
│   │   ├── prompt/            # Prompt building
│   │   ├── document/          # Text chunking
│   │   └── middleware/        # HTTP middleware
│   └── documents/             # Document directories per profile
└── web/
    ├── index.html             # Main interface
    ├── css/style.css          # Styling
    └── js/scripts.js          # Frontend logic
```

## Key Design Patterns

### Memory Architecture
- **Short-term**: Conversational history with sliding window
- **Long-term**: Vector embeddings with semantic search
- **Metadata filtering**: Prevents cross-contamination between document and conversation memories

### Document Chunking
- Respects Markdown structure (headers, paragraphs)
- Smart sentence-based splitting
- Code block preservation
- Configurable chunk sizes

### Hot Reload
- Profile switching without restart
- Proper watcher lifecycle management
- Thread-safe operations with mutexes
- Goroutine management with stop channels

## API Endpoints

- `GET /` - Health check
- `POST /chat` - Send chat message
- `GET /profiles` - List available profiles
- `GET /profile/active` - Get current active profile
- `POST /profile/switch` - Switch to different profile

## Configuration Options

### Profile Settings
- `directories`: Paths to watch for documents
- `memory_file`: JSON file for storing memories
- `index_file`: JSON file for index state
- `extensions`: Allowed file extensions
- `max_file_size`: Maximum file size in bytes

## Development Notes

### Thread Safety
- Memory operations use `sync.RWMutex`
- Config access is protected
- Proper goroutine cleanup on profile switch

### Performance
- Hash-based change detection avoids redundant indexing
- Vector similarity using cosine distance
- Efficient metadata filtering

## Limitations & Future Work

- Currently single-user (no authentication)
- In-memory vector store (consider persistent storage)
- Basic chunking algorithm (could use more sophisticated methods)
- No conversation branching or editing

## Cross-Platform Memory Sharing

For sharing memory files across Windows/Linux machines:
- Use SMB/CIFS network shares
- Mount shared directory containing memory files
- Update `config.json` paths to point to network location

## License

Havent thought of it yet

## Contributing

This is a learning project built to understand AI integration patterns. Feedback and suggestions welcome!

## Acknowledgments

- Built with [Ollama](https://ollama.ai) for local LLM inference
- Uses [Brave Search API](https://brave.com/search/api/) for web search
- Inspired by modern RAG architectures and semantic memory systems
