#!/bin/bash

# ChatBotCore Docker Management Script

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Function to build the Docker image
build() {
    print_status "Building ChatBotCore Docker image..."
    docker-compose build --no-cache
    print_success "Docker image built successfully!"
}

# Function to start services
start() {
    print_status "Starting ChatBotCore services..."
    docker-compose up -d
    
    print_status "Waiting for services to be healthy..."
    sleep 10
    
    # Check if services are running
    if docker-compose ps | grep -q "chatbot-core.*Up"; then
        print_success "ChatBotCore is running!"
        print_status "Service URLs:"
        echo "  - ChatBotCore API: http://localhost:8086"
        echo "  - Health Check: http://localhost:8086/health"
        echo "  - PostgreSQL: localhost:55433"
    else
        print_error "Failed to start services. Check logs with: ./docker.sh logs"
        exit 1
    fi
}

# Function to stop services
stop() {
    print_status "Stopping ChatBotCore services..."
    docker-compose down
    print_success "Services stopped!"
}

# Function to restart services
restart() {
    print_status "Restarting ChatBotCore services..."
    docker-compose down
    docker-compose up -d
    print_success "Services restarted!"
}

# Function to show logs
logs() {
    if [ -n "$2" ]; then
        docker-compose logs -f "$2"
    else
        docker-compose logs -f
    fi
}

# Function to show status
status() {
    print_status "Service Status:"
    docker-compose ps
    
    print_status "Health Checks:"
    curl -s http://localhost:8086/health | jq . 2>/dev/null || curl -s http://localhost:8086/health
}

# Function to clean up
clean() {
    print_warning "This will remove all containers, images, and volumes!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Cleaning up..."
        docker-compose down -v --rmi all
        docker system prune -f
        print_success "Cleanup completed!"
    else
        print_status "Cleanup cancelled."
    fi
}

# Function to run database migration
migrate() {
    print_status "Running database migration..."
    docker-compose run --rm chatbot-migrate
    print_success "Migration completed!"
}

# Function to access database
db() {
    print_status "Connecting to PostgreSQL database..."
    docker-compose exec postgres psql -U chatbot_user -d chatbot_mvp
}

# Function to show help
help() {
    echo "ChatBotCore Docker Management Script"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  build       Build Docker images"
    echo "  start       Start all services"
    echo "  stop        Stop all services"
    echo "  restart     Restart all services"
    echo "  logs        Show logs (use 'logs <service>' for specific service)"
    echo "  status      Show service status and health"
    echo "  migrate     Run database migrations"
    echo "  db          Connect to PostgreSQL database"
    echo "  clean       Remove all containers, images, and volumes"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start                 # Start all services"
    echo "  $0 logs chatbot-core     # Show logs for ChatBotCore service"
    echo "  $0 status                # Check service status"
}

# Main script logic
main() {
    check_docker
    
    case "${1:-help}" in
        build)
            build
            ;;
        start)
            start
            ;;
        stop)
            stop
            ;;
        restart)
            restart
            ;;
        logs)
            logs "$@"
            ;;
        status)
            status
            ;;
        migrate)
            migrate
            ;;
        db)
            db
            ;;
        clean)
            clean
            ;;
        help|--help|-h)
            help
            ;;
        *)
            print_error "Unknown command: $1"
            help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
